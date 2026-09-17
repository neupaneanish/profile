package service

import (
	"context"

	"neupaneanish.com.np/profile/internal/enum"
	profilev1 "neupaneanish.com.np/profile/internal/protobuf/gateway/profile/v1"
	"neupaneanish.com.np/profile/internal/repository"
	"neupaneanish.com.np/profile/internal/utils"
)

func (s *GatewayProfileService) CreateSocial(
	ctx context.Context,
	req *profilev1.CreateSocialRequest,
) (*profilev1.CreateSocialResponse, error) {
	serviceName := "CreateSocial"
	userSession := utils.UserSessionContext(ctx)

	iconID, iconIDErr := utils.ParseUUID(ctx, req.GetIconId(), serviceName, s.cfg.Logger)
	if iconIDErr != nil {
		return nil, iconIDErr
	}

	params := &repository.CreateSocialParams{
		UserID:    userSession.UserID,
		IconID:    iconID,
		Username:  req.GetUsername(),
		CreatedBy: userSession.UserID,
		UpdatedBy: userSession.UserID,
	}

	_, err := s.cfg.Repository.CreateSocial(ctx, params)
	if sErr := socialError(ctx, err, serviceName, "create", s.cfg.Logger); sErr != nil {
		return nil, sErr
	}

	payload := utils.RedpandaRootEventNotificationPayload{
		ActorID:  userSession.UserID,
		Username: userSession.Username,
		UserID:   userSession.UserID,
		Method:   enum.DBMethodCreate,
		Table:    enum.DBTableSocial,
	}
	s.cfg.Redpanda.Produce(ctx, utils.RedpandaRootDatabaseEventNotifications, serviceName, payload)

	return &profilev1.CreateSocialResponse{}, nil
}
