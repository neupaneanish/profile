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

	platformID, platformIDErr := parseUUID(ctx, req.GetPlatformId(), serviceName, s.cfg.Logger)
	if platformIDErr != nil {
		return nil, platformIDErr
	}

	params := &repository.CreateSocialParams{
		UserID:     userSession.UserID,
		PlatformID: platformID,
		Username:   req.GetUsername(),
		CreatedBy:  userSession.UserID,
		UpdatedBy:  userSession.UserID,
	}

	id, err := s.cfg.Repository.CreateSocial(ctx, params)
	if sErr := socialError(ctx, err, serviceName, "create", s.cfg.Logger); sErr != nil {
		return nil, sErr
	}

	return &profilev1.CreateSocialResponse{Id: id.String()}, nil
}
