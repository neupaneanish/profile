package service

import (
	"context"

	"neupaneanish.com.np/profile/internal/enum"
	rootProfilev1 "neupaneanish.com.np/profile/internal/protobuf/root/profile/v1"
	"neupaneanish.com.np/profile/internal/utils"
)

func (s *RootProfileService) UpdateIcon(
	ctx context.Context,
	req *rootProfilev1.UpdateIconRequest,
) (*rootProfilev1.UpdateIconResponse, error) {
	serviceName := "UpdatePlatform"
	userSession := utils.UserSessionContext(ctx)

	if err := s.createUpdateIcon(
		ctx,
		req.GetId(),
		req.GetIcon(),
		serviceName,
		req.GetUpdatedAt().AsTime(),
	); err != nil {
		return nil, err
	}

	payload := utils.RedpandaRootEventNotificationPayload{
		ActorID:  userSession.UserID,
		Username: userSession.Username,
		UserID:   userSession.UserID,
		Method:   enum.DBMethodUpdate,
		Table:    enum.DBTableIcon,
	}
	s.cfg.Redpanda.Produce(ctx, utils.RedpandaRootDatabaseEventNotifications, serviceName, payload)

	return &rootProfilev1.UpdateIconResponse{}, nil
}
