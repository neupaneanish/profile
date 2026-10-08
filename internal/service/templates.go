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

func (s *GatewayProfileService) Templates(
	ctx context.Context,
	_ *gatewayProfilev1.TemplatesRequest,
) (*gatewayProfilev1.TemplatesResponse, error) {
	res, err := templates(ctx, s.cfg.Repository, s.cfg.Logger, "GatewayTemplate")
	if err != nil {
		return nil, err
	}

	return &gatewayProfilev1.TemplatesResponse{Templates: res}, nil
}

func (s *RootProfileService) Templates(
	ctx context.Context,
	_ *rootProfilev1.TemplatesRequest,
) (*rootProfilev1.TemplatesResponse, error) {
	res, err := templates(ctx, s.cfg.Repository, s.cfg.Logger, "RootTemplate")
	if err != nil {
		return nil, err
	}

	return &rootProfilev1.TemplatesResponse{Templates: res}, nil
}

func templates(
	ctx context.Context,
	querier repository.Querier,
	logger *slog.Logger,
	serviceName string,
) ([]*profilev1.Template, error) {
	rows, err := querier.Templates(ctx)
	if err != nil {
		logger.WarnContext(
			ctx,
			"Failed to query templates",
			"service",
			serviceName,
		)
		return nil, errs.ErrInternalServer
	}

	res := make([]*profilev1.Template, len(rows))

	for i, row := range rows {
		res[i] = &profilev1.Template{
			Id:          row.ID.String(),
			Name:        row.Name,
			Description: row.Description,
			Icon:        row.Icon,
		}
	}

	return res, nil
}

func (s *RootProfileService) TemplateIcons(
	ctx context.Context,
	_ *rootProfilev1.TemplateIconsRequest,
) (*rootProfilev1.TemplateIconsResponse, error) {
	serviceName := "RootTemplateIcons"

	rows, err := s.cfg.Repository.TemplateIcons(ctx)
	if err != nil {
		s.cfg.Logger.WarnContext(
			ctx,
			"Failed to query template icons",
			"service",
			serviceName,
		)
		return nil, errs.ErrInternalServer
	}

	res := make([]*profilev1.Icon, len(rows))
	for i, row := range rows {
		res[i] = &profilev1.Icon{
			Id:   row.ID.String(),
			Name: row.Name,
			Icon: row.Icon,
		}
	}

	return &rootProfilev1.TemplateIconsResponse{Icons: res}, nil
}
