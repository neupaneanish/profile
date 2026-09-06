package service

import (
	"context"
	"log/slog"
	"time"
	"uuid"

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

	id, err := createUpdateExperience(
		ctx,
		req.GetId(),
		userSession.UserID,
		userSession.UserID,
		req.GetExperience(),
		req.GetUpdatedAt().AsTime(),
		s.cfg.Repository,
		s.cfg.Logger,
		serviceName,
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

	id, err := createUpdateExperience(
		ctx,
		req.GetId(),
		userID,
		userSession.UserID,
		req.GetExperience(),
		req.GetUpdatedAt().AsTime(),
		s.cfg.Repository,
		s.cfg.Logger,
		serviceName,
	)
	if err != nil {
		return nil, err
	}

	return &rootProfilev1.UpdateExperienceResponse{Id: id.String()}, nil
}

func createUpdateExperience(
	ctx context.Context,
	id string,
	userID, updatedBy uuid.UUID,
	req *profilev1.CreateUpdateExperience,
	updatedAt time.Time,
	repo repository.Querier,
	logger *slog.Logger,
	serviceName string,
) (uuid.UUID, error) {
	if id == "" {
		params := &repository.CreateExperienceParams{
			UserID:       userID,
			Title:        req.GetTitle(),
			CompanyName:  req.GetCompanyName(),
			Location:     req.GetLocation(),
			LocationType: enum.LocationType(req.GetLocationType()),
			StartDate:    req.GetStartDate().AsTime(),
			EndDate:      utils.TimestampValue(req.GetEndDate()),
			Description:  utils.StringValue(req.GetDescription()),
			CreatedBy:    userID,
			UpdatedBy:    userID,
		}

		idx, err := repo.CreateExperience(ctx, params)
		if err != nil {
			logger.ErrorContext(ctx, "Create Experience Failed", "service", serviceName, "error", err)
			return uuid.Nil(), errs.ErrInternalServer
		}
		return idx, nil
	}

	idx, idErr := parseUUID(ctx, id, serviceName, logger)
	if idErr != nil {
		return uuid.Nil(), idErr
	}

	params := &repository.UpdateExperienceParams{
		Title:        req.GetTitle(),
		CompanyName:  req.GetCompanyName(),
		Location:     req.GetLocation(),
		LocationType: enum.LocationType(req.GetLocationType()),
		StartDate:    req.GetStartDate().AsTime(),
		EndDate:      utils.TimestampValue(req.GetEndDate()),
		Description:  utils.StringValue(req.GetDescription()),
		UpdatedBy:    updatedBy,
		ID:           idx,
		UserID:       userID,
		UpdatedAt:    updatedAt,
	}

	cmdTag, err := repo.UpdateExperience(ctx, params)
	if err != nil {
		logger.ErrorContext(ctx, "Update Education Failed", "service", serviceName, "error", err)
		return uuid.Nil(), errs.ErrInternalServer
	}

	if cmdTag.RowsAffected() == 0 {
		logger.WarnContext(
			ctx,
			"Concurrent experience update detected",
			"service",
			serviceName,
			"userID", userID,
			"id", id,
		)
		return uuid.Nil(), errs.ErrConflict
	}
	return idx, nil
}
