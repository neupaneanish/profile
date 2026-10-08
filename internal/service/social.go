package service

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"google.golang.org/protobuf/types/known/timestamppb"
	"neupaneanish.com.np/profile/internal/errs"
	rootProfilev1 "neupaneanish.com.np/profile/internal/protobuf/root/profile/v1"
	"neupaneanish.com.np/profile/internal/repository"
	"neupaneanish.com.np/profile/internal/utils"
)

func (s *RootProfileService) Social(
	ctx context.Context,
	req *rootProfilev1.SocialRequest,
) (*rootProfilev1.SocialResponse, error) {
	serviceName := "RootSocial"
	userSession := utils.UserSessionContext(ctx)

	id, idErr := utils.ParseUUID(ctx, req.GetId(), serviceName, s.cfg.Logger)
	if idErr != nil {
		return nil, idErr
	}

	userID, userIDErr := utils.ParseUUID(ctx, req.GetUserId(), serviceName, s.cfg.Logger)
	if userIDErr != nil {
		return nil, userIDErr
	}

	params := &repository.SocialParams{
		ID:     id,
		UserID: userID,
	}

	row, err := s.cfg.Repository.Social(ctx, params)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			s.cfg.Logger.WarnContext(ctx, "Failed to query social", "service", serviceName, "error", err)
			return nil, errs.ErrNotFound("Social")
		}
		s.cfg.Logger.ErrorContext(ctx, "Failed to fetch social", "service", serviceName, "error", err)
		return nil, errs.ErrInternalServer
	}

	createdUsername, updatedUsername, usernamesErr := utils.GetUsernames(
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

	return &rootProfilev1.SocialResponse{
		Id:                row.ID.String(),
		UserId:            row.UserID.String(),
		IconId:            row.IconID.String(),
		Icon:              row.Icon,
		Username:          row.Username,
		Name:              row.Name,
		Social:            row.Social,
		CreatedAt:         timestamppb.New(row.CreatedAt),
		CreatedBy:         row.CreatedBy.String(),
		UpdatedAt:         timestamppb.New(row.UpdatedAt),
		UpdatedBy:         row.UpdatedBy.String(),
		CreatedByUsername: createdUsername,
		UpdatedByUsername: updatedUsername,
	}, nil
}
