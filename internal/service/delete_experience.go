//nolint:dupl // Clean handler pattern intentionally mirrors Experience Delete
package service

import (
	"context"

	"neupaneanish.com.np/profile/internal/errs"
	gatewayProfilev1 "neupaneanish.com.np/profile/internal/protobuf/gateway/profile/v1"
	"neupaneanish.com.np/profile/internal/repository"
	"neupaneanish.com.np/profile/internal/utils"
)

func (s *GatewayProfileService) DeleteExperience(
	ctx context.Context,
	req *gatewayProfilev1.DeleteExperienceRequest,
) (*gatewayProfilev1.DeleteExperienceResponse, error) {
	serviceName := "DeleteExperience"
	userSession := utils.UserSessionContext(ctx)

	id, idErr := parseUUID(ctx, req.GetId(), serviceName, s.cfg.Logger)
	if idErr != nil {
		return nil, idErr
	}

	params := &repository.DeleteExperienceParams{
		ID:        id,
		UserID:    userSession.UserID,
		UpdatedAt: req.GetUpdatedAt().AsTime(),
	}

	cmdTag, err := s.cfg.Repository.DeleteExperience(ctx, params)
	if err != nil {
		s.cfg.Logger.ErrorContext(ctx, "Delete Experience Failed", "service", serviceName, "error", err)
		return nil, errs.ErrInternalServer
	}

	if cmdTag.RowsAffected() == 0 {
		s.cfg.Logger.WarnContext(
			ctx,
			"Experience record not found or concurrent modification",
			"service", serviceName,
			"id", id.String(),
			"userID", userSession.UserID,
		)
		return nil, errs.ErrConflict
	}

	return &gatewayProfilev1.DeleteExperienceResponse{}, nil
}
