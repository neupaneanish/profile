package service

import (
	"context"

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

	iconID, iconIDErr := parseUUID(ctx, req.GetIconId(), serviceName, s.cfg.Logger)
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

	return &profilev1.CreateSocialResponse{}, nil
}
