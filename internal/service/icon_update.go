package service

import (
	"context"

	rootProfilev1 "neupaneanish.com.np/profile/internal/protobuf/root/profile/v1"
	"neupaneanish.com.np/profile/internal/redpanda"
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

	return &rootProfilev1.UpdateIconResponse{}, nil
}
