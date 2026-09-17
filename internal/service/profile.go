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

func (s *GatewayProfileService) Profile(
	ctx context.Context,
	_ *gatewayProfilev1.ProfileRequest,
) (*gatewayProfilev1.ProfileResponse, error) {
	serviceName := "GatewayProfile"
	userSession := utils.UserSessionContext(ctx)

	res, err := profile(ctx, userSession.UserID, serviceName, s.cfg.Repository, s.cfg.Logger)
	if err != nil {
		return nil, err
	}
	return &gatewayProfilev1.ProfileResponse{Profile: res}, nil
}

func (s *RootProfileService) Profile(
	ctx context.Context,
	req *rootProfilev1.ProfileRequest,
) (*rootProfilev1.ProfileResponse, error) {
	serviceName := "RootProfile"
	userID, userIDErr := utils.ParseUUID(ctx, req.GetUserId(), serviceName, s.cfg.Logger)
	if userIDErr != nil {
		return nil, userIDErr
	}

	res, err := profile(ctx, userID, serviceName, s.cfg.Repository, s.cfg.Logger)
	if err != nil {
		return nil, err
	}
	return &rootProfilev1.ProfileResponse{Profile: res}, nil
}

func (s *ExternalProfileService) Profile(
	ctx context.Context,
	_ *externalProfilev1.ProfileRequest,
) (*externalProfilev1.ProfileResponse, error) {
	serviceName := "ExternalProfile"
	userID := utils.DomainUserSessionContext(ctx)

	res, err := profile(ctx, userID, serviceName, s.cfg.Repository, s.cfg.Logger)
	if err != nil {
		return nil, err
	}

	return &externalProfilev1.ProfileResponse{
		Name:  res.GetName(),
		Title: res.GetTitle(),
	}, nil
}

func profile(
	ctx context.Context,
	userID uuid.UUID,
	serviceName string,
	repo repository.Querier,
	logger *slog.Logger,
) (*profilev1.Profile, error) {
	params := &repository.ProfileParams{UserID: userID}
	row, rowErr := repo.Profile(ctx, params)
	if rowErr != nil {
		if errors.Is(rowErr, pgx.ErrNoRows) {
			logger.WarnContext(ctx, "Profile not found", "service", serviceName)
			return nil, errs.ErrNotFound("Profile")
		}
		logger.ErrorContext(ctx, "Profile query failed", "service", serviceName, "error", rowErr)
		return nil, errs.ErrInternalServer
	}
	return &profilev1.Profile{
		UserId:    row.UserID.String(),
		Name:      row.Name,
		Title:     row.Title,
		Dob:       timestamppb.New(row.Dob),
		CreatedAt: timestamppb.New(row.CreatedAt),
		CreatedBy: row.CreatedBy.String(),
		UpdatedAt: timestamppb.New(row.UpdatedAt),
		UpdatedBy: row.UpdatedBy.String(),
	}, nil
}
