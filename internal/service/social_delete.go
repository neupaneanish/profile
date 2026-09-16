package service

import (
	"context"

	profilev1 "neupaneanish.com.np/profile/internal/protobuf/gateway/profile/v1"
	"neupaneanish.com.np/profile/internal/repository"
	"neupaneanish.com.np/profile/internal/utils"
)

func (s *GatewayProfileService) DeleteSocial(
	ctx context.Context,
	req *profilev1.DeleteSocialRequest,
) (*profilev1.DeleteSocialResponse, error) {
	serviceName := "DeleteSocial"
	userSession := utils.UserSessionContext(ctx)

	id, idErr := utils.ParseUUID(ctx, req.GetId(), serviceName, s.cfg.Logger)
	if idErr != nil {
		return nil, idErr
	}

	params := &repository.DeleteSocialParams{
		ID:        id,
		UserID:    userSession.UserID,
		UpdatedAt: req.GetUpdatedAt().AsTime(),
	}

	affected, err := s.cfg.Repository.DeleteSocial(ctx, params)
	if deleteErr := deleteDB(
		ctx,
		affected,
		err,
		serviceName,
		"Social",
		req.GetId(),
		s.cfg.Logger,
	); deleteErr != nil {
		return nil, deleteErr
	}

	payload := utils.RedpandaRootEventNotificationPayload{
		ActorID:  userSession.UserID,
		Username: userSession.Username,
		UserID:   userSession.UserID,
		Message:  "social deleted",
	}
	s.cfg.Redpanda.Produce(ctx, utils.RedpandaRootEventNotifications, serviceName, payload)

	return &profilev1.DeleteSocialResponse{}, nil
}
