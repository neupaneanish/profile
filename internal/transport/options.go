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
	externalProfilev1 "neupaneanish.com.np/profile/internal/protobuf/external/profile/v1"
	gatewayProfilev1 "neupaneanish.com.np/profile/internal/protobuf/gateway/profile/v1"
	rootProfilev1 "neupaneanish.com.np/profile/internal/protobuf/root/profile/v1"

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

	authFunc := func(ctx context.Context) (context.Context, error) {
		return AuthInterceptor(ctx, cfg.Logger, externalEndpoints(), gatewayEndpoints(), rootEndpoints())
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

func externalEndpoints() map[string]struct{} {
	return map[string]struct{}{
		externalProfilev1.ExternalProfileService_Profile_FullMethodName:     {},
		externalProfilev1.ExternalProfileService_About_FullMethodName:       {},
		externalProfilev1.ExternalProfileService_Educations_FullMethodName:  {},
		externalProfilev1.ExternalProfileService_Experiences_FullMethodName: {},
		externalProfilev1.ExternalProfileService_Socials_FullMethodName:     {},
	}
}

func gatewayEndpoints() map[string]struct{} {
	return map[string]struct{}{
		// Profile
		gatewayProfilev1.GatewayProfileService_Profile_FullMethodName:       {},
		gatewayProfilev1.GatewayProfileService_CreateProfile_FullMethodName: {},
		gatewayProfilev1.GatewayProfileService_UpdateProfile_FullMethodName: {},
		gatewayProfilev1.GatewayProfileService_CheckProfile_FullMethodName:  {},
		// About
		gatewayProfilev1.GatewayProfileService_About_FullMethodName:       {},
		gatewayProfilev1.GatewayProfileService_CreateAbout_FullMethodName: {},
		gatewayProfilev1.GatewayProfileService_UpdateAbout_FullMethodName: {},
		// Domain
		gatewayProfilev1.GatewayProfileService_CreateDomain_FullMethodName: {},
		gatewayProfilev1.GatewayProfileService_VerifyDomain_FullMethodName: {},
		gatewayProfilev1.GatewayProfileService_DeleteDomain_FullMethodName: {},
		gatewayProfilev1.GatewayProfileService_Domains_FullMethodName:      {},
		// Education
		gatewayProfilev1.GatewayProfileService_CreateEducation_FullMethodName: {},
		gatewayProfilev1.GatewayProfileService_UpdateEducation_FullMethodName: {},
		gatewayProfilev1.GatewayProfileService_DeleteEducation_FullMethodName: {},
		gatewayProfilev1.GatewayProfileService_Education_FullMethodName:       {},
		gatewayProfilev1.GatewayProfileService_Educations_FullMethodName:      {},
		// Experience
		gatewayProfilev1.GatewayProfileService_CreateExperience_FullMethodName: {},
		gatewayProfilev1.GatewayProfileService_UpdateExperience_FullMethodName: {},
		gatewayProfilev1.GatewayProfileService_DeleteExperience_FullMethodName: {},
		gatewayProfilev1.GatewayProfileService_Experience_FullMethodName:       {},
		gatewayProfilev1.GatewayProfileService_Experiences_FullMethodName:      {},
		// Social
		gatewayProfilev1.GatewayProfileService_SocialIcons_FullMethodName:  {},
		gatewayProfilev1.GatewayProfileService_CreateSocial_FullMethodName: {},
		gatewayProfilev1.GatewayProfileService_UpdateSocial_FullMethodName: {},
		gatewayProfilev1.GatewayProfileService_DeleteSocial_FullMethodName: {},
		gatewayProfilev1.GatewayProfileService_Socials_FullMethodName:      {},
	}
}

func rootEndpoints() map[string]struct{} {
	return map[string]struct{}{
		// Profile
		rootProfilev1.RootProfileService_UpdateProfile_FullMethodName: {},
		rootProfilev1.RootProfileService_Profile_FullMethodName:       {},
		// About
		rootProfilev1.RootProfileService_About_FullMethodName:       {},
		rootProfilev1.RootProfileService_UpdateAbout_FullMethodName: {},
		// Domain
		rootProfilev1.RootProfileService_Domains_FullMethodName: {},
		// Nameserver
		rootProfilev1.RootProfileService_CreateNameserver_FullMethodName: {},
		rootProfilev1.RootProfileService_DeleteNameserver_FullMethodName: {},
		rootProfilev1.RootProfileService_Nameservers_FullMethodName:      {},
		// Education
		rootProfilev1.RootProfileService_UpdateEducation_FullMethodName: {},
		rootProfilev1.RootProfileService_Education_FullMethodName:       {},
		rootProfilev1.RootProfileService_Educations_FullMethodName:      {},
		// Experience
		rootProfilev1.RootProfileService_UpdateExperience_FullMethodName: {},
		rootProfilev1.RootProfileService_Experience_FullMethodName:       {},
		rootProfilev1.RootProfileService_Experiences_FullMethodName:      {},
		// Icon
		rootProfilev1.RootProfileService_Icons_FullMethodName:      {},
		rootProfilev1.RootProfileService_Icon_FullMethodName:       {},
		rootProfilev1.RootProfileService_CreateIcon_FullMethodName: {},
		rootProfilev1.RootProfileService_UpdateIcon_FullMethodName: {},
		rootProfilev1.RootProfileService_DeleteIcon_FullMethodName: {},
		// Social
		rootProfilev1.RootProfileService_UpdateSocial_FullMethodName: {},
		rootProfilev1.RootProfileService_Socials_FullMethodName:      {},
	}
}
