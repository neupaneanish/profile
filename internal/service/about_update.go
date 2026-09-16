package service

import (
	"context"
	"errors"
	"log/slog"
	"time"
	"uuid"

	"github.com/jackc/pgx/v5"
	"google.golang.org/protobuf/types/known/timestamppb"
	"neupaneanish.com.np/profile/internal/config"
	"neupaneanish.com.np/profile/internal/errs"
	profilev1 "neupaneanish.com.np/profile/internal/protobuf/common/profile/v1"
	gatewayProfilev1 "neupaneanish.com.np/profile/internal/protobuf/gateway/profile/v1"
	rootProfilev1 "neupaneanish.com.np/profile/internal/protobuf/root/profile/v1"
	"neupaneanish.com.np/profile/internal/repository"
	"neupaneanish.com.np/profile/internal/utils"
)

func (s *GatewayProfileService) UpdateAbout(
	ctx context.Context,
	req *gatewayProfilev1.UpdateAboutRequest,
) (*gatewayProfilev1.UpdateAboutResponse, error) {
	serviceName := "GatewayUpdateAbout"
	userSession := utils.UserSessionContext(ctx)

	res, err := updateAbout(
		ctx,
		userSession.UserID,
		userSession.UserID,
		req.GetUpdatedAt().AsTime(),
		req.GetAbout(),
		userSession.Username,
		serviceName,
		s.cfg.Repository,
		s.cfg.Redpanda,
		s.cfg.Logger,
	)
	if err != nil {
		return nil, err
	}
	return &gatewayProfilev1.UpdateAboutResponse{About: res}, nil
}

func (s *RootProfileService) UpdateAbout(
	ctx context.Context,
	req *rootProfilev1.UpdateAboutRequest,
) (*rootProfilev1.UpdateAboutResponse, error) {
	serviceName := "RootUpdateAbout"
	userSession := utils.UserSessionContext(ctx)

	userID, userIDErr := utils.ParseUUID(ctx, req.GetUserId(), serviceName, s.cfg.Logger)
	if userIDErr != nil {
		return nil, userIDErr
	}

	res, err := updateAbout(
		ctx,
		userID,
		userSession.UserID,
		req.GetUpdatedAt().AsTime(),
		req.GetAbout(),
		serviceName,
		userSession.Username,
		s.cfg.Repository,
		s.cfg.Redpanda,
		s.cfg.Logger,
	)
	if err != nil {
		return nil, err
	}

	return &rootProfilev1.UpdateAboutResponse{About: res}, nil
}

func updateAbout(
	ctx context.Context,
	userID, updatedBy uuid.UUID,
	updatedAt time.Time,
	about, serviceName, username string,
	repo repository.Querier,
	redpanda *config.Redpanda,
	logger *slog.Logger,
) (*profilev1.About, error) {
	params := &repository.UpdateAboutParams{
		About:     about,
		UpdatedBy: updatedBy,
		UserID:    userID,
		UpdatedAt: updatedAt,
	}
	row, err := repo.UpdateAbout(ctx, params)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			logger.WarnContext(
				ctx,
				"Concurrent profile update detected",
				"service",
				serviceName,
				"userID",
				userID.String(),
			)
			return nil, errs.ErrConflict
		}
		logger.ErrorContext(ctx, "Update Profile Failed", "service", serviceName, "error", err)
		return nil, errs.ErrInternalServer
	}

	payload := utils.RedpandaRootEventNotificationPayload{
		ActorID:  updatedBy,
		Username: username,
		UserID:   userID,
		Message:  "updated about",
	}

	redpanda.Produce(ctx, utils.RedpandaRootEventNotifications, serviceName, payload)

	return &profilev1.About{
		UserId:    row.UserID.String(),
		About:     row.About,
		CreatedAt: timestamppb.New(row.CreatedAt),
		CreatedBy: row.CreatedBy.String(),
		UpdatedAt: timestamppb.New(row.UpdatedAt),
		UpdatedBy: row.UpdatedBy.String(),
	}, nil
}
