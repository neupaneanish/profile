package service

import (
	"context"

	profilev1 "neupaneanish.com.np/profile/internal/protobuf/root/profile/v1"
	"neupaneanish.com.np/profile/internal/repository"
	"neupaneanish.com.np/profile/internal/utils"
)

func (s *RootProfileService) DeleteNameserver(
	ctx context.Context,
	req *profilev1.DeleteNameserverRequest,
) (*profilev1.DeleteNameserverResponse, error) {
	serviceName := "DeleteNameserver"
	userSession := utils.UserSessionContext(ctx)

	id, idErr := utils.ParseUUID(ctx, req.GetId(), serviceName, s.cfg.Logger)
	if idErr != nil {
		return nil, idErr
	}

	params := &repository.DeleteNameserverParams{
		ID:        id,
		UpdatedAt: req.GetUpdatedAt().AsTime(),
	}

	affected, err := s.cfg.Repository.DeleteNameserver(ctx, params)
	if deleteErr := deleteDB(
		ctx,
		affected,
		err,
		serviceName,
		"Nameserver",
		req.GetId(),
		userSession.Username,
		userSession.UserID,
		userSession.UserID,
		s.cfg.Logger,
		s.cfg.Redpanda,
	); deleteErr != nil {
		return nil, deleteErr
	}

	payload := utils.RedpandaRootEventNotificationPayload{
		ActorID:  userSession.UserID,
		Username: userSession.Username,
		UserID:   userSession.UserID,
		Message:  "nameserver deleted",
	}
	s.cfg.Redpanda.Produce(ctx, utils.RedpandaRootEventNotifications, serviceName, payload)

	return &profilev1.DeleteNameserverResponse{}, nil
}
