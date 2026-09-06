package service

import (
	"context"
	"errors"
	"log/slog"
	"time"
	"uuid"

	"github.com/jackc/pgx/v5"
	"google.golang.org/protobuf/types/known/timestamppb"
	"neupaneanish.com.np/profile/internal/errs"
	profilev1 "neupaneanish.com.np/profile/internal/protobuf/common/profile/v1"
	gatewayProfilev1 "neupaneanish.com.np/profile/internal/protobuf/gateway/profile/v1"
	rootProfilev1 "neupaneanish.com.np/profile/internal/protobuf/root/profile/v1"
	"neupaneanish.com.np/profile/internal/repository"
	"neupaneanish.com.np/profile/internal/utils"
)

func (s *GatewayProfileService) UpdateProfile(
	ctx context.Context,
	req *gatewayProfilev1.UpdateProfileRequest,
) (*gatewayProfilev1.UpdateProfileResponse, error) {
	serviceName := "GatewayUpdateProfile"
	userSession := utils.UserSessionContext(ctx)

	res, err := updateProfile(
		ctx,
		userSession.UserID,
		userSession.UserID,
		req.GetProfile(),
		req.GetUpdatedAt().AsTime(),
		s.cfg.Repository,
		s.cfg.Logger,
		serviceName,
	)
	if err != nil {
		return nil, err
	}

	return &gatewayProfilev1.UpdateProfileResponse{
		Profile: res,
	}, nil
}

func (s *RootProfileService) UpdateProfile(
	ctx context.Context,
	req *rootProfilev1.UpdateProfileRequest,
) (*rootProfilev1.UpdateProfileResponse, error) {
	serviceName := "RootUpdateProfile"
	userSession := utils.UserSessionContext(ctx)

	userID, userIDErr := parseUUID(ctx, req.GetUserId(), serviceName, s.cfg.Logger)
	if userIDErr != nil {
		return nil, userIDErr
	}

	res, err := updateProfile(
		ctx,
		userID,
		userSession.UserID,
		req.GetProfile(),
		req.GetUpdatedAt().AsTime(),
		s.cfg.Repository,
		s.cfg.Logger,
		serviceName,
	)
	if err != nil {
		return nil, err
	}

	return &rootProfilev1.UpdateProfileResponse{
		Profile: res,
	}, nil
}

func updateProfile(
	ctx context.Context,
	userID, updatedBy uuid.UUID,
	req *profilev1.CreateUpdateProfile,
	updatedAt time.Time,
	repo repository.Querier,
	logger *slog.Logger,
	serviceName string,
) (*profilev1.Profile, error) {
	params := &repository.UpdateProfileParams{
		Name:      req.GetName(),
		Title:     req.GetTitle(),
		UpdatedBy: updatedBy,
		UserID:    userID,
		UpdatedAt: updatedAt,
	}

	row, rowErr := repo.UpdateProfile(ctx, params)
	if rowErr != nil {
		if errors.Is(rowErr, pgx.ErrNoRows) {
			logger.WarnContext(ctx, "Concurrent profile update detected", "service", serviceName, "userID", userID)
			return nil, errs.ErrConflict
		}
		logger.ErrorContext(ctx, "Update Profile Failed", "service", serviceName, "error", rowErr)
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
