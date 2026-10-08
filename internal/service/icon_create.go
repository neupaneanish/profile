package service

import (
	"context"
	"time"

	rootProfilev1 "neupaneanish.com.np/profile/internal/protobuf/root/profile/v1"
	"neupaneanish.com.np/profile/internal/redpanda"
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

	redpanda.RootNotificationProduce(
		ctx,
		userSession,
		userSession.UserID,
		utils.DatabaseTableDomain,
		utils.DatabaseMethodCreate,
		serviceName,
		s.cfg.Client,
		s.cfg.Redpanda,
		s.cfg.Logger,
	)

	return &rootProfilev1.CreateIconResponse{}, nil
}
