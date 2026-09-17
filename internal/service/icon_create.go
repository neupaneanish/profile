package service

import (
	"context"
	"time"

	"neupaneanish.com.np/profile/internal/enum"
	rootProfilev1 "neupaneanish.com.np/profile/internal/protobuf/root/profile/v1"
	"neupaneanish.com.np/profile/internal/utils"
)

func (s *RootProfileService) CreateIcon(
	ctx context.Context,
	req *rootProfilev1.CreateIconRequest,
) (*rootProfilev1.CreateIconResponse, error) {
	serviceName := "CreateIcon"
	userSession := utils.UserSessionContext(ctx)

	if err := s.createUpdateIcon(
		ctx,
		"",
		req.GetIcon(),
		serviceName,
		time.Time{},
	); err != nil {
		return nil, err
	}

	payload := utils.RedpandaRootEventNotificationPayload{
		ActorID:  userSession.UserID,
		Username: userSession.Username,
		UserID:   userSession.UserID,
		Method:   enum.DBMethodCreate,
		Table:    enum.DBTableIcon,
	}
	s.cfg.Redpanda.Produce(ctx, utils.RedpandaRootDatabaseEventNotifications, serviceName, payload)

	return &rootProfilev1.CreateIconResponse{}, nil
}
