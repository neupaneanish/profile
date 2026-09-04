//go:build integration || benchmark || e2e

package service_test

import (
	"context"
	"database/sql"
	"errors"
	"log/slog"
	"net"
	"os"
	"path/filepath"
	"runtime"
	"testing"
	"time"
	"uuid"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/test/bufconn"
	"neupaneanish.com.np/profile/internal/config"
	"neupaneanish.com.np/profile/internal/enum"
	externalProfilev1 "neupaneanish.com.np/profile/internal/protobuf/external/profile/v1"
	gatewayProfilev1 "neupaneanish.com.np/profile/internal/protobuf/gateway/profile/v1"
	rootProfilev1 "neupaneanish.com.np/profile/internal/protobuf/root/profile/v1"
	"neupaneanish.com.np/profile/internal/repository"
	"neupaneanish.com.np/profile/internal/service"
	"neupaneanish.com.np/profile/internal/telemetry"
	"neupaneanish.com.np/profile/internal/transport"
	"neupaneanish.com.np/profile/tests"

	// Register the file source driver for migrations.
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

var (
	cfg *config.Config

	externalProfileServiceClient externalProfilev1.ExternalProfileServiceClient
	gatewayProfileServiceClient  gatewayProfilev1.GatewayProfileServiceClient
	rootProfileServiceClient     rootProfilev1.RootProfileServiceClient
)

type container struct {
	dbURL            string
	dbCleanup        func()
	vkURL            string
	vkCleanup        func()
	telemetryURL     string
	telemetryCleanup func()
}

func TestMain(m *testing.M) {
	ctx := context.Background()

	baseLogger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	testContainer := setupContainer(baseLogger)

	testEnv := setupEnv(testContainer.dbURL, testContainer.vkURL)

	logger, loggerCleanup, loggerErr := telemetry.NewTelemetry(
		ctx,
		testContainer.telemetryURL,
		testEnv.ServiceName,
		testEnv.Environment,
	)

	if loggerErr != nil {
		baseLogger.Error("Failed to start telemetry", "error", loggerErr)
		os.Exit(1)
	}

	testCfg, testCfgErr := config.NewConfig(ctx, testEnv, logger)
	if testCfgErr != nil {
		baseLogger.Error("Failed to setup config", "error", testCfgErr)
		os.Exit(1)
	}

	cfg = testCfg

	client, server, testClientServerErr := testClientServer(testCfg, baseLogger)
	if testClientServerErr != nil {
		baseLogger.Error("Failed to start client / server")
		os.Exit(1)
	}

	externalProfileServiceClient = externalProfilev1.NewExternalProfileServiceClient(client)
	gatewayProfileServiceClient = gatewayProfilev1.NewGatewayProfileServiceClient(client)
	rootProfileServiceClient = rootProfilev1.NewRootProfileServiceClient(client)

	code := m.Run()
	loggerCleanupErr := loggerCleanup(ctx)
	if loggerCleanupErr != nil {
		baseLogger.Error("Failed to cleanup logger", "error", loggerCleanupErr)
		os.Exit(1)
	}

	err := client.Close()
	if err != nil {
		baseLogger.Error("Failed to close client", "error", err)
		os.Exit(1)
	}
	server.GracefulStop()
	testContainer.dbCleanup()
	testContainer.vkCleanup()
	testContainer.telemetryCleanup()

	os.Exit(code)
}

func setupContainer(logger *slog.Logger) *container {
	dbURL, dbCleanup, dbErr := tests.Postgres()
	if dbErr != nil {
		logger.Error("Failed to start postgres container", "error", dbErr)
		os.Exit(1)
	}

	migrationErr := runMigrations(dbURL)
	if migrationErr != nil {
		logger.Error("Failed to migrations", "error", migrationErr)
		os.Exit(1)
	}

	vkURL, vkCleanup, vkErr := tests.Valkey()
	if vkErr != nil {
		logger.Error("Failed to start valkey container", "error", vkErr)
		os.Exit(1)
	}

	telemetryURL, telemetryCleanup, telemetryErr := tests.OpenTelemetry()
	if telemetryErr != nil {
		logger.Error("Failed to start telemetry container", "error", telemetryErr)
		os.Exit(1)
	}

	return &container{
		dbURL:            dbURL,
		dbCleanup:        dbCleanup,
		vkURL:            vkURL,
		vkCleanup:        vkCleanup,
		telemetryURL:     telemetryURL,
		telemetryCleanup: telemetryCleanup,
	}
}

func runMigrations(url string) error {
	_, b, _, _ := runtime.Caller(0)
	root := filepath.Join(filepath.Dir(b), "..", "..")
	migrationsPath := filepath.Join(root, "database", "profile/migrations")

	db, dbErr := sql.Open("postgres", url)
	if dbErr != nil {
		return dbErr
	}

	defer func() {
		_ = db.Close()
	}()

	driver, driverErr := postgres.WithInstance(db, &postgres.Config{})
	if driverErr != nil {
		return driverErr
	}

	m, mErr := migrate.NewWithDatabaseInstance(
		"file://"+migrationsPath,
		"postgres",
		driver,
	)
	if mErr != nil {
		return mErr
	}

	defer func() {
		_, _ = m.Close()
	}()

	upErr := m.Up()
	if upErr != nil && !errors.Is(upErr, migrate.ErrNoChange) {
		return upErr
	}

	return nil
}

func setupEnv(db string, vk string) *config.Env {
	return &config.Env{
		DatabaseURL: db,
		ValkeyURL:   vk,
		Environment: "test",
		ServiceName: "Test",
	}
}

func testClientServer(cfg *config.Config, logger *slog.Logger) (*grpc.ClientConn, *grpc.Server, error) {
	listen := bufconn.Listen(1024 * 1024)

	opts, optsErr := transport.NewOptions(cfg)

	if optsErr != nil {
		return nil, nil, optsErr
	}

	server := grpc.NewServer(opts...)

	externalProfilev1.RegisterExternalProfileServiceServer(
		server,
		service.NewExternalProfileService(cfg),
	)
	gatewayProfilev1.RegisterGatewayProfileServiceServer(
		server,
		service.NewGatewayProfileService(cfg),
	)
	rootProfilev1.RegisterRootProfileServiceServer(
		server,
		service.NewRootProfileService(cfg),
	)

	go func() {
		if err := server.Serve(listen); err != nil {
			logger.Error("Failed to serve server", "error", err)
			os.Exit(1)
		}
	}()

	client, clientErr := grpc.NewClient(
		"passthrough://bufnet",
		grpc.WithContextDialer(func(context.Context, string) (net.Conn, error) {
			return listen.Dial()
		}),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)

	if clientErr != nil {
		return nil, nil, clientErr
	}

	return client, server, nil
}

func contextWithValue(t *testing.T, userID uuid.UUID, role enum.UserRole) context.Context {
	t.Helper()

	md := metadata.Pairs(
		"x-user-id", userID.String(),
		"x-role", string(role),
		"x-jti", uuid.NewV7().String(),
	)

	ctx := metadata.NewOutgoingContext(t.Context(), md)
	return ctx
}

func seedProfile(t *testing.T, userID uuid.UUID) {
	t.Helper()

	params := &repository.CreateProfileParams{
		UserID:    userID,
		Name:      "Anish Neupane",
		Title:     "Backend Developer",
		Dob:       time.Now(),
		CreatedBy: userID,
		UpdatedBy: userID,
	}
	_, createErr := cfg.Repository.CreateProfile(t.Context(), params)
	require.NoError(t, createErr)
}

func getProfile(t *testing.T, userID uuid.UUID) *repository.Profile {
	t.Helper()
	seedProfile(t, userID)
	params := &repository.ProfileParams{UserID: userID}
	profile, err := cfg.Repository.Profile(t.Context(), params)
	require.NoError(t, err)
	return profile
}
