package service

import (
	"context"
	"errors"
	"log/slog"
	"uuid"

	"github.com/jackc/pgx/v5"
	"google.golang.org/protobuf/types/known/timestamppb"
	"neupaneanish.com.np/profile/internal/errs"
	profilev1 "neupaneanish.com.np/profile/internal/protobuf/common/profile/v1"
	externalProfilev1 "neupaneanish.com.np/profile/internal/protobuf/external/profile/v1"
	gatewayProfilev1 "neupaneanish.com.np/profile/internal/protobuf/gateway/profile/v1"
	rootProfilev1 "neupaneanish.com.np/profile/internal/protobuf/root/profile/v1"
	"neupaneanish.com.np/profile/internal/repository"
	"neupaneanish.com.np/profile/internal/utils"
)

func (s *GatewayProfileService) About(
	ctx context.Context,
	_ *gatewayProfilev1.AboutRequest,
) (*gatewayProfilev1.AboutResponse, error) {
	serviceName := "GatewayAbout"
	userSession := utils.UserSessionContext(ctx)

	res, err := about(ctx, userSession.UserID, s.cfg.Repository, s.cfg.Logger, serviceName)
	if err != nil {
		return nil, err
	}

	return &gatewayProfilev1.AboutResponse{About: res}, nil
}

func (s *RootProfileService) About(
	ctx context.Context,
	req *rootProfilev1.AboutRequest,
) (*rootProfilev1.AboutResponse, error) {
	serviceName := "RootAbout"
	userID, userIDErr := utils.ParseUUID(ctx, req.GetUserId(), serviceName, s.cfg.Logger)
	if userIDErr != nil {
		return nil, userIDErr
	}

	res, err := about(ctx, userID, s.cfg.Repository, s.cfg.Logger, serviceName)
	if err != nil {
		return nil, err
	}

	return &rootProfilev1.AboutResponse{About: res}, nil
}

func (s *ExternalProfileService) About(
	ctx context.Context,
	_ *externalProfilev1.AboutRequest,
) (*externalProfilev1.AboutResponse, error) {
	serviceName := "ExternalAbout"
	userID := utils.DomainUserSessionContext(ctx)

	res, err := about(ctx, userID, s.cfg.Repository, s.cfg.Logger, serviceName)
	if err != nil {
		return nil, err
	}

	return &externalProfilev1.AboutResponse{About: res.GetAbout()}, nil
}

func about(
	ctx context.Context,
	userID uuid.UUID,
	repo repository.Querier,
	logger *slog.Logger,
	serviceName string,
) (*profilev1.About, error) {
	params := &repository.AboutParams{UserID: userID}

	row, err := repo.About(ctx, params)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			logger.WarnContext(ctx, "About not found", "service", serviceName)
			return nil, errs.ErrNotFound("About")
		}
		logger.ErrorContext(ctx, "About query failed", "service", serviceName, "error", err)
		return nil, errs.ErrInternalServer
	}

	return &profilev1.About{
		UserId:    row.UserID.String(),
		About:     row.About,
		CreatedAt: timestamppb.New(row.CreatedAt),
		CreatedBy: row.CreatedBy.String(),
		UpdatedAt: timestamppb.New(row.UpdatedAt),
		UpdatedBy: row.UpdatedBy.String(),
	}, nil
}
