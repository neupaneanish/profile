package service

import (
	"context"
	"log/slog"

	"uuid"

	"neupaneanish.com.np/profile/internal/errs"
	profilev1 "neupaneanish.com.np/profile/internal/protobuf/common/profile/v1"
	externalProfilev1 "neupaneanish.com.np/profile/internal/protobuf/external/profile/v1"
	gatewayProfilev1 "neupaneanish.com.np/profile/internal/protobuf/gateway/profile/v1"
	"neupaneanish.com.np/profile/internal/repository"
	"neupaneanish.com.np/profile/internal/utils"
)

func (s *GatewayProfileService) Exists(
	ctx context.Context,
	_ *gatewayProfilev1.ExistsRequest,
) (*gatewayProfilev1.ExistsResponse, error) {
	serviceName := "GatewayExists"
	userSession := utils.UserSessionContext(ctx)

	res, err := exists(ctx, userSession.UserID, serviceName, s.cfg.Repository, s.cfg.Logger)
	if err != nil {
		return nil, err
	}

	return &gatewayProfilev1.ExistsResponse{Exists: res}, nil
}

func (s *ExternalProfileService) Exists(
	ctx context.Context,
	_ *externalProfilev1.ExistsRequest,
) (*externalProfilev1.ExistsResponse, error) {
	serviceName := "ExternalExists"
	userID := utils.DomainUserSessionContext(ctx)

	res, err := exists(ctx, userID, serviceName, s.cfg.Repository, s.cfg.Logger)
	if err != nil {
		return nil, err
	}

	return &externalProfilev1.ExistsResponse{Exists: res}, nil
}

func exists(
	ctx context.Context,
	userID uuid.UUID,
	serviceName string,
	repo repository.Querier,
	logger *slog.Logger,
) (*profilev1.Exists, error) {
	params := &repository.ExistsParams{UserID: userID}
	row, err := repo.Exists(ctx, params)
	if err != nil {
		logger.ErrorContext(ctx, "Failed to query Exists", "service", serviceName, "error", err)
		return nil, errs.ErrInternalServer
	}

	return &profilev1.Exists{
		Profile:     row.Profile,
		About:       row.About,
		Educations:  row.Educations,
		Experiences: row.Experiences,
		Socials:     row.Socials,
	}, nil
}
