package service_test

import (
	"context"
	"crypto/rand"
	"database/sql"
	"errors"
	"log/slog"
	"net"
	"os"
	"path/filepath"
	"runtime"
	"strings"
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
	"neupaneanish.com.np/profile/internal/utils"
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
	rpURL            string
	rpCleanup        func()
	telemetryURL     string
	telemetryCleanup func()
}

func TestMain(m *testing.M) {
	ctx := context.Background()

	baseLogger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	testContainer := setupContainer(baseLogger)

	testEnv := setupEnv(testContainer.dbURL, testContainer.vkURL, testContainer.rpURL)

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

	rpURL, rpCleanup, rpErr := tests.Redpanda()
	if rpErr != nil {
		logger.Error("Failed to start redpanda container", "error", telemetryErr)
		os.Exit(1)
	}

	return &container{
		dbURL:            dbURL,
		dbCleanup:        dbCleanup,
		vkURL:            vkURL,
		vkCleanup:        vkCleanup,
		rpURL:            rpURL,
		rpCleanup:        rpCleanup,
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

func setupEnv(db, vk, rp string) *config.Env {
	return &config.Env{
		DatabaseURL: db,
		ValkeyURL:   vk,
		RedpandaURL: rp,
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
		"x-username", rand.Text()[:8],
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
	createErr := cfg.Repository.CreateProfile(t.Context(), params)
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

func seedEducation(t *testing.T, userID uuid.UUID) uuid.UUID {
	t.Helper()
	createParams := &repository.CreateEducationParams{
		UserID:        userID,
		School:        "Westcliff University",
		Degree:        "Master of Science in Computer Science",
		Affiliation:   nil,
		FieldOfStudy:  nil,
		Concentration: nil,
		StartDate:     time.Date(2024, time.September, 1, 0, 0, 0, 0, time.UTC),
		EndDate:       nil,
		Address:       "United State",
		Description:   nil,
		CreatedBy:     userID,
		UpdatedBy:     userID,
	}
	id, err := cfg.Repository.CreateEducation(t.Context(), createParams)
	require.NoError(t, err)
	return id
}

func getEducation(t *testing.T, userID uuid.UUID) *repository.Education {
	t.Helper()
	id := seedEducation(t, userID)

	getParams := &repository.EducationParams{
		ID:     id,
		UserID: userID,
	}
	education, eduErr := cfg.Repository.Education(t.Context(), getParams)
	require.NoError(t, eduErr)

	return education
}

func seedExperience(t *testing.T, userID uuid.UUID) uuid.UUID {
	t.Helper()

	var description *string
	descriptionValue := "Computer Science Teacher"
	description = &descriptionValue

	var endDate *time.Time
	endDateValue := time.Date(2024, time.August, 14, 0, 0, 0, 0, time.UTC)
	endDate = &endDateValue

	createParams := &repository.CreateExperienceParams{
		UserID:       userID,
		Title:        "Teacher",
		CompanyName:  "Oxford College of Engineering and Management",
		Location:     "Nepal",
		LocationType: enum.LocationTypeHybrid,
		StartDate:    time.Date(2023, time.June, 15, 0, 0, 0, 0, time.UTC),
		EndDate:      endDate,
		Description:  description,
		CreatedBy:    userID,
		UpdatedBy:    userID,
	}
	id, err := cfg.Repository.CreateExperience(t.Context(), createParams)
	require.NoError(t, err)
	return id
}

func getExperience(t *testing.T, userID uuid.UUID) *repository.Experience {
	t.Helper()
	id := seedExperience(t, userID)

	getParams := &repository.ExperienceParams{
		ID:     id,
		UserID: userID,
	}
	experience, eduErr := cfg.Repository.Experience(t.Context(), getParams)
	require.NoError(t, eduErr)

	return experience
}

func seedAbout(t *testing.T) *repository.About {
	t.Helper()
	userID := uuid.NewV7()
	params := &repository.CreateAboutParams{
		UserID:    userID,
		About:     rand.Text() + rand.Text(),
		CreatedBy: userID,
		UpdatedBy: userID,
	}

	row, err := cfg.Repository.CreateAbout(t.Context(), params)
	require.NoError(t, err)
	return row
}

func externalContextWithValue(t *testing.T, userID uuid.UUID, hostname string) context.Context {
	t.Helper()

	md := metadata.Pairs(
		"x-hostname", hostname,
		"x-user-id", userID.String(),
	)

	cmd := cfg.Client.B().Hset().
		Key(utils.DomainUserSessionKey).
		FieldValue().
		FieldValue(hostname, userID.String()).
		Build()

	err := cfg.Client.Do(t.Context(), cmd).Error()
	require.NoError(t, err)

	ctx := metadata.NewOutgoingContext(t.Context(), md)
	return ctx
}

func seedNameserver(t *testing.T, ip, ipType string) uuid.UUID {
	t.Helper()

	params := &repository.CreateNameserverParams{
		Ip:        ip,
		IpType:    ipType,
		CreatedBy: uuid.Nil(),
		UpdatedBy: uuid.Nil(),
	}

	id, err := cfg.Repository.CreateNameserver(t.Context(), params)
	require.NoError(t, err)
	return id
}

func getNameserver(t *testing.T, ip, ipType string) *repository.Nameserver {
	t.Helper()

	id := seedNameserver(t, ip, ipType)

	ns, err := cfg.Repository.Nameservers(t.Context())
	require.NoError(t, err)

	for _, n := range ns {
		if n.ID == id {
			return n
		}
	}
	return nil
}

func seedDomain(t *testing.T, userID uuid.UUID, url, ip, ipType string) uuid.UUID {
	t.Helper()

	nsID := seedNameserver(t, ip, ipType)

	params := &repository.CreateDomainParams{
		UserID:       userID,
		NameserverID: nsID,
		Fqdn:         url,
		Txt:          rand.Text(),
		CreatedBy:    userID,
		UpdatedBy:    userID,
	}

	id, err := cfg.Repository.CreateDomain(t.Context(), params)
	require.NoError(t, err)
	return id
}

func getDomain(t *testing.T, userID uuid.UUID, url, ip, ipType string) *repository.DomainRow {
	t.Helper()

	id := seedDomain(t, userID, url, ip, ipType)

	params := &repository.DomainParams{
		ID:     id,
		UserID: userID,
	}

	domain, err := cfg.Repository.Domain(t.Context(), params)
	require.NoError(t, err)
	return domain
}

func seedIcon(t *testing.T, name string, siteSuffix *string) uuid.UUID {
	t.Helper()
	params := &repository.CreateIconParams{
		Name:       name,
		Site:       name + ".com",
		SiteSuffix: siteSuffix,
		Url:        name + ".com",
		Slug:       name,
		Color:      "#FFFFFF",
		CreatedBy:  uuid.Nil(),
		UpdatedBy:  uuid.Nil(),
	}

	id, err := cfg.Repository.CreateIcon(t.Context(), params)
	require.NoError(t, err)
	return id
}

func getIcon(t *testing.T, name string, siteSuffix *string) *repository.IconRow {
	t.Helper()
	id := seedIcon(t, name, siteSuffix)

	params := &repository.IconParams{ID: id}

	icon, iconErr := cfg.Repository.Icon(t.Context(), params)
	require.NoError(t, iconErr)
	return icon
}

func seedSocial(t *testing.T, userID uuid.UUID, username string) (uuid.UUID, uuid.UUID) {
	t.Helper()

	name := strings.ToLower(rand.Text()[:8])
	icon := getIcon(t, name, &name)

	params := &repository.CreateSocialParams{
		UserID:    userID,
		IconID:    icon.ID,
		Username:  username,
		CreatedBy: userID,
		UpdatedBy: userID,
	}
	id, err := cfg.Repository.CreateSocial(t.Context(), params)
	require.NoError(t, err)

	return icon.ID, id
}
