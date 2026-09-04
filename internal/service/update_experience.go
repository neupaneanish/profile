package service

import (
	"context"
	"errors"
	"log/slog"
	"uuid"

	"github.com/jackc/pgx/v5"
	"neupaneanish.com.np/profile/internal/enum"
	"neupaneanish.com.np/profile/internal/errs"
	profilev1 "neupaneanish.com.np/profile/internal/protobuf/common/profile/v1"
	gatewayProfilev1 "neupaneanish.com.np/profile/internal/protobuf/gateway/profile/v1"
	rootProfilev1 "neupaneanish.com.np/profile/internal/protobuf/root/profile/v1"
	"neupaneanish.com.np/profile/internal/repository"
	"neupaneanish.com.np/profile/internal/utils"
)

func (s *GatewayProfileService) UpdateExperience(
	ctx context.Context,
	req *gatewayProfilev1.UpdateExperienceRequest,
) (*gatewayProfilev1.UpdateExperienceResponse, error) {
	serviceName := "GatewayUpdateExperience"
	userSession := utils.UserSessionContext(ctx)

	id, err := updateExperience(
		ctx,
		req.GetExperience(),
		userSession.UserID,
		userSession.UserID,
		s.cfg.Repository,
		serviceName,
		s.cfg.Logger,
	)
	if err != nil {
		return nil, err
	}

	return &gatewayProfilev1.UpdateExperienceResponse{Id: id.String()}, nil
}

func (s *RootProfileService) UpdateExperience(
	ctx context.Context,
	req *rootProfilev1.UpdateExperienceRequest,
) (*rootProfilev1.UpdateExperienceResponse, error) {
	serviceName := "RootUpdateExperience"
	userSession := utils.UserSessionContext(ctx)

	userID, userIDErr := parseUUID(ctx, req.GetUserId(), serviceName, s.cfg.Logger)
	if userIDErr != nil {
		return nil, userIDErr
	}

	id, err := updateExperience(
		ctx,
		req.GetExperience(),
		userID,
		userSession.UserID,
		s.cfg.Repository,
		serviceName,
		s.cfg.Logger,
	)
	if err != nil {
		return nil, err
	}

	return &rootProfilev1.UpdateExperienceResponse{Id: id.String()}, nil
}

func updateExperience(
	ctx context.Context,
	req *profilev1.UpdateExperience,
	userID, updatedBy uuid.UUID,
	repo repository.Querier,
	serviceName string,
	logger *slog.Logger,
) (uuid.UUID, error) {
	id, idErr := parseUUID(ctx, req.GetId(), serviceName, logger)
	if idErr != nil {
		return id, idErr
	}

	params := &repository.UpdateExperienceParams{
		Title:        req.GetTitle(),
		CompanyName:  req.GetCompanyName(),
		Location:     req.GetLocation(),
		LocationType: enum.LocationType(req.GetLocationType()),
		StartDate:    req.GetStartDate().AsTime(),
		EndDate:      timestampValue(req.GetEndDate()),
		Description:  stringValue(req.GetDescription()),
		UpdatedBy:    updatedBy,
		ID:           id,
		UserID:       userID,
		UpdatedAt:    req.GetUpdatedAt().AsTime(),
	}

	idx, err := repo.UpdateExperience(ctx, params)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			logger.WarnContext(
				ctx,
				"Concurrent experience update detected",
				"service",
				serviceName,
				"userID",
				userID,
				"id",
				id.String(),
			)
			return uuid.Nil(), errs.ErrConflict
		}
		logger.ErrorContext(ctx, "Update Experience Failed", "service", serviceName, "error", err)
		return uuid.Nil(), errs.ErrInternalServer
	}

	return idx, nil
}
