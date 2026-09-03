package transport

import (
	"google.golang.org/grpc"
	"neupaneanish.com.np/profile/internal/config"
	externalProfilev1 "neupaneanish.com.np/profile/internal/protobuf/external/profile/v1"
	gatewayProfilev1 "neupaneanish.com.np/profile/internal/protobuf/gateway/profile/v1"
	rootProfilev1 "neupaneanish.com.np/profile/internal/protobuf/root/profile/v1"
	"neupaneanish.com.np/profile/internal/service"
)

func register(
	cfg *config.Config,
	server *grpc.Server,
) {
	externalProfileService := service.NewExternalProfileService(cfg)
	externalProfilev1.RegisterExternalProfileServiceServer(server, externalProfileService)

	gatewayProfileService := service.NewGatewayProfileService(cfg)
	gatewayProfilev1.RegisterGatewayProfileServiceServer(server, gatewayProfileService)

	rootProfileService := service.NewRootProfileService(cfg)
	rootProfilev1.RegisterRootProfileServiceServer(server, rootProfileService)
}
