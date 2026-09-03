package service

import (
	"neupaneanish.com.np/profile/internal/config"
	externalProfilev1 "neupaneanish.com.np/profile/internal/protobuf/external/profile/v1"
	gatewayProfilev1 "neupaneanish.com.np/profile/internal/protobuf/gateway/profile/v1"
	rootProfilev1 "neupaneanish.com.np/profile/internal/protobuf/root/profile/v1"
)

type ExternalProfileService struct {
	externalProfilev1.UnimplementedExternalProfileServiceServer

	cfg *config.Config
}

type GatewayProfileService struct {
	gatewayProfilev1.UnimplementedGatewayProfileServiceServer

	cfg *config.Config
}

type RootProfileService struct {
	rootProfilev1.UnimplementedRootProfileServiceServer

	cfg *config.Config
}

func NewExternalProfileService(cfg *config.Config) *ExternalProfileService {
	return &ExternalProfileService{
		cfg: cfg,
	}
}

func NewGatewayProfileService(cfg *config.Config) *GatewayProfileService {
	return &GatewayProfileService{
		cfg: cfg,
	}
}

func NewRootProfileService(cfg *config.Config) *RootProfileService {
	return &RootProfileService{
		cfg: cfg,
	}
}
