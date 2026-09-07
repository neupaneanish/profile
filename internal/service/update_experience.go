package service

import (
	"context"

	gatewayProfilev1 "neupaneanish.com.np/profile/internal/protobuf/gateway/profile/v1"
	rootProfilev1 "neupaneanish.com.np/profile/internal/protobuf/root/profile/v1"
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
