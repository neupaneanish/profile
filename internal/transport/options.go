package transport

import (
	"context"
	"time"

	"buf.build/go/protovalidate"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/auth"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/recovery"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc/filters"
	"neupaneanish.com.np/profile/internal/config"
	"neupaneanish.com.np/profile/internal/errs"

	protovalidatemiddleware "github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/protovalidate"
	"google.golang.org/grpc"
)

const (
	interceptorTimeout = 30 * time.Second
	maxTimeout         = 5 * time.Minute
)

func NewOptions(cfg *config.Config) ([]grpc.ServerOption, error) {
	oTelHandler := otelgrpc.NewServerHandler(
		otelgrpc.WithFilter(filters.Not(filters.HealthCheck())),
	)

	validator, validatorErr := protovalidate.New()
	if validatorErr != nil {
		cfg.Logger.Error("proto validate", "error", validatorErr)
		return nil, validatorErr
	}

	recoveryOpt := recovery.WithRecoveryHandler(func(p any) error {
		cfg.Logger.Error("panic recovered in gRPC handler", "panic", p)
		return errs.ErrInternalServer
	})

	externalEndpoints := map[string]struct{}{}

	gatewayEndpoints := map[string]struct{}{}

	rootEndpoints := map[string]struct{}{}

	authFunc := func(ctx context.Context) (context.Context, error) {
		return AuthInterceptor(ctx, externalEndpoints, gatewayEndpoints, rootEndpoints)
	}

	opts := []grpc.ServerOption{
		grpc.StatsHandler(oTelHandler),
		grpc.ChainUnaryInterceptor(
			recovery.UnaryServerInterceptor(recoveryOpt),
			UnaryTimeoutInterceptor(interceptorTimeout),
			protovalidatemiddleware.UnaryServerInterceptor(validator),
			logging.UnaryServerInterceptor(
				LoggerInterceptor(cfg.Logger),
				logging.WithLogOnEvents(logging.StartCall, logging.FinishCall),
			),
			auth.UnaryServerInterceptor(authFunc),
		),
		grpc.ChainStreamInterceptor(
			recovery.StreamServerInterceptor(recoveryOpt),
			StreamTimeoutInterceptor(maxTimeout),
			protovalidatemiddleware.StreamServerInterceptor(validator),
			logging.StreamServerInterceptor(
				LoggerInterceptor(cfg.Logger),
				logging.WithLogOnEvents(logging.StartCall, logging.FinishCall),
			),
			auth.StreamServerInterceptor(authFunc),
		),
	}

	return opts, nil
}
