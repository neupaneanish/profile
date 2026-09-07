package service

import (
	"context"

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
	if deleteErr := deleteDB(
		ctx,
		cmdTag,
		err,
		serviceName,
		"Education",
		req.GetId(),
		s.cfg.Logger,
	); deleteErr != nil {
		return nil, deleteErr
	}

	return &gatewayProfilev1.DeleteEducationResponse{}, nil
}
