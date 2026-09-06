package service

import (
	"context"
	"log/slog"

	"neupaneanish.com.np/profile/internal/errs"
	profilev1 "neupaneanish.com.np/profile/internal/protobuf/common/profile/v1"
	gatewayProfilev1 "neupaneanish.com.np/profile/internal/protobuf/gateway/profile/v1"
	rootProfilev1 "neupaneanish.com.np/profile/internal/protobuf/root/profile/v1"
	"neupaneanish.com.np/profile/internal/repository"
)

func (s *GatewayProfileService) Platforms(
	ctx context.Context,
	_ *gatewayProfilev1.PlatformsRequest,
) (*gatewayProfilev1.PlatformsResponse, error) {
	serviceName := "GatewayPlatforms"

	res, err := platforms(ctx, s.cfg.Repository, s.cfg.Logger, serviceName)
	if err != nil {
		return nil, err
	}

	return &gatewayProfilev1.PlatformsResponse{Platforms: res}, nil
}

func (s *RootProfileService) Platforms(
	ctx context.Context,
	_ *rootProfilev1.PlatformsRequest,
) (*rootProfilev1.PlatformsResponse, error) {
	serviceName := "RootPlatforms"

	res, err := platforms(ctx, s.cfg.Repository, s.cfg.Logger, serviceName)
	if err != nil {
		return nil, err
	}

	return &rootProfilev1.PlatformsResponse{Platforms: res}, nil
}

func platforms(
	ctx context.Context,
	repo repository.Querier,
	logger *slog.Logger,
	serviceName string,
) ([]*profilev1.Platforms, error) {
	rows, err := repo.Platforms(ctx)
	if err != nil {
		logger.ErrorContext(ctx, "Failed to query platforms", "service", serviceName, "error", err)
		return nil, errs.ErrInternalServer
	}

	res := make([]*profilev1.Platforms, len(rows))
	for i, p := range rows {
		res[i] = &profilev1.Platforms{
			Id:        p.ID.String(),
			Name:      p.Name,
			Url:       p.Url,
			UrlSuffix: p.UrlSuffix,
			Logo:      p.Logo,
			Color:     p.Color,
		}
	}

	return res, nil
}
