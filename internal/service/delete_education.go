package service

import (
	"context"

	"neupaneanish.com.np/profile/internal/errs"
	gatewayProfilev1 "neupaneanish.com.np/profile/internal/protobuf/gateway/profile/v1"
	"neupaneanish.com.np/profile/internal/repository"
	"neupaneanish.com.np/profile/internal/utils"
)

func (s *GatewayProfileService) DeleteEducation(
	ctx context.Context,
	req *gatewayProfilev1.DeleteEducationRequest,
) (*gatewayProfilev1.DeleteEducationResponse, error) {
	serviceName := "DeleteEducation"
	userSession := utils.UserSessionContext(ctx)

	id, idErr := parseUUID(ctx, req.GetId(), serviceName, s.cfg.Logger)
	if idErr != nil {
		return nil, idErr
	}

	params := &repository.DeleteEducationParams{
		ID:        id,
		UserID:    userSession.UserID,
		UpdatedAt: req.GetUpdatedAt().AsTime(),
	}

	cmdTag, err := s.cfg.Repository.DeleteEducation(ctx, params)
	if err != nil {
		s.cfg.Logger.ErrorContext(ctx, "Delete Education Failed", "service", serviceName, "error", err)
		return nil, errs.ErrInternalServer
	}

	if cmdTag.RowsAffected() == 0 {
		s.cfg.Logger.WarnContext(
			ctx,
			"Education record not found or concurrent modification",
			"service", serviceName,
			"id", id.String(),
			"userID", userSession.UserID,
		)
		return nil, errs.ErrConflict
	}

	return &gatewayProfilev1.DeleteEducationResponse{}, nil
}
