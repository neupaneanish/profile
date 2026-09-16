package service

import (
	"context"

	gatewayProfilev1 "neupaneanish.com.np/profile/internal/protobuf/gateway/profile/v1"
	rootProfilev1 "neupaneanish.com.np/profile/internal/protobuf/root/profile/v1"
	"neupaneanish.com.np/profile/internal/utils"
)

func (s *GatewayProfileService) UpdateEducation(
	ctx context.Context,
	req *gatewayProfilev1.UpdateEducationRequest,
) (*gatewayProfilev1.UpdateEducationResponse, error) {
	serviceName := "GatewayUpdateEducation"
	userSession := utils.UserSessionContext(ctx)

	if err := createUpdateEducation(
		ctx,
		req.GetId(),
		userSession.UserID,
		userSession.UserID,
		req.GetEducation(),
		req.GetUpdatedAt().AsTime(),
		s.cfg.Repository,
		s.cfg.Logger,
		serviceName,
	); err != nil {
		return nil, err
	}

	payload := utils.RedpandaRootEventNotificationPayload{
		ActorID:  userSession.UserID,
		Username: userSession.Username,
		UserID:   userSession.UserID,
		Message:  "education updated",
	}
	s.cfg.Redpanda.Produce(ctx, utils.RedpandaRootEventNotifications, serviceName, payload)

	return &gatewayProfilev1.UpdateEducationResponse{}, nil
}

func (s *RootProfileService) UpdateEducation(
	ctx context.Context,
	req *rootProfilev1.UpdateEducationRequest,
) (*rootProfilev1.UpdateEducationResponse, error) {
	serviceName := "RootUpdateEducation"
	userSession := utils.UserSessionContext(ctx)

	userID, userIDErr := utils.ParseUUID(ctx, req.GetUserId(), serviceName, s.cfg.Logger)
	if userIDErr != nil {
		return nil, userIDErr
	}

	if err := createUpdateEducation(
		ctx,
		req.GetId(),
		userID,
		userSession.UserID,
		req.GetEducation(),
		req.GetUpdatedAt().AsTime(),
		s.cfg.Repository,
		s.cfg.Logger,
		serviceName,
	); err != nil {
		return nil, err
	}

	payload := utils.RedpandaRootEventNotificationPayload{
		ActorID:  userSession.UserID,
		Username: userSession.Username,
		UserID:   userID,
		Message:  "education updated",
	}
	s.cfg.Redpanda.Produce(ctx, utils.RedpandaRootEventNotifications, serviceName, payload)

	return &rootProfilev1.UpdateEducationResponse{}, nil
}
