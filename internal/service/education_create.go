package service

import (
	"context"
	"time"

	gatewayProfilev1 "neupaneanish.com.np/profile/internal/protobuf/gateway/profile/v1"
	"neupaneanish.com.np/profile/internal/utils"
)

func (s *GatewayProfileService) CreateEducation(
	ctx context.Context, req *gatewayProfilev1.CreateEducationRequest,
) (*gatewayProfilev1.CreateEducationResponse, error) {
	serviceName := "CreateEducation"
	userSession := utils.UserSessionContext(ctx)

	if err := createUpdateEducation(
		ctx,
		"",
		userSession.UserID,
		userSession.UserID,
		req.GetEducation(),
		time.Time{},
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
		Message:  "education created",
	}
	s.cfg.Redpanda.Produce(ctx, utils.RedpandaRootEventNotifications, serviceName, payload)

	return &gatewayProfilev1.CreateEducationResponse{}, nil
}
