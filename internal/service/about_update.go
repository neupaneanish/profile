package service

import (
	"context"
	"errors"
	"log/slog"
	"time"
	"uuid"

	"github.com/jackc/pgx/v5"
	"github.com/twmb/franz-go/pkg/kgo"
	"github.com/valkey-io/valkey-go"
	"google.golang.org/protobuf/types/known/timestamppb"
	"neupaneanish.com.np/profile/internal/errs"
	gatewayProfilev1 "neupaneanish.com.np/profile/internal/protobuf/gateway/profile/v1"
	rootProfilev1 "neupaneanish.com.np/profile/internal/protobuf/root/profile/v1"
	"neupaneanish.com.np/profile/internal/redpanda"
	"neupaneanish.com.np/profile/internal/repository"
	"neupaneanish.com.np/profile/internal/utils"
)

func (s *GatewayProfileService) UpdateAbout(
	ctx context.Context,
	req *gatewayProfilev1.UpdateAboutRequest,
) (*gatewayProfilev1.UpdateAboutResponse, error) {
	serviceName := "GatewayUpdateAbout"
	userSession := utils.UserSessionContext(ctx)

	row, err := updateAbout(
		ctx,
		userSession,
		userSession.UserID,
		req.GetUpdatedAt().AsTime(),
		req.GetAbout(),
		serviceName,
		s.cfg.Repository,
		s.cfg.Client,
		s.cfg.Redpanda,
		s.cfg.Logger,
	)
	if err != nil {
		return nil, err
	}

	res := &gatewayProfilev1.About{
		UserId:    row.UserID.String(),
		About:     row.About,
		UpdatedAt: timestamppb.New(row.UpdatedAt),
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

	row, err := updateAbout(
		ctx,
		userSession,
		userID,
		req.GetUpdatedAt().AsTime(),
		req.GetAbout(),
		serviceName,
		s.cfg.Repository,
		s.cfg.Client,
		s.cfg.Redpanda,
		s.cfg.Logger,
	)
	if err != nil {
		return nil, err
	}

	createdByUsername, updatedByUsername, usernamesErr := utils.GetUsernames(
		ctx,
		row.CreatedBy,
		row.UpdatedBy,
		userSession,
		s.cfg.Client,
		s.cfg.Logger,
	)
	if usernamesErr != nil {
		return nil, usernamesErr
	}

	res := &rootProfilev1.About{
		UserId:            row.UserID.String(),
		About:             row.About,
		CreatedAt:         timestamppb.New(row.CreatedAt),
		CreatedBy:         row.CreatedBy.String(),
		UpdatedAt:         timestamppb.New(row.UpdatedAt),
		UpdatedBy:         row.UpdatedBy.String(),
		CreatedByUsername: createdByUsername,
		UpdatedByUsername: updatedByUsername,
	}

	return &rootProfilev1.UpdateAboutResponse{About: res}, nil
}

func updateAbout(
	ctx context.Context,
	session *utils.UserSession,
	userID uuid.UUID,
	updatedAt time.Time,
	about, serviceName string,
	repo repository.Querier,
	vkClient valkey.Client,
	client *kgo.Client,
	logger *slog.Logger,
) (*repository.About, error) {
	params := &repository.UpdateAboutParams{
		About:     about,
		UpdatedBy: session.UserID,
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

	redpanda.RootNotificationProduce(
		ctx,
		session,
		userID,
		utils.DatabaseTableAbout,
		utils.DatabaseMethodUpdate,
		serviceName,
		vkClient,
		client,
		logger,
	)

	return row, nil
}
