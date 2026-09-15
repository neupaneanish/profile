package service

import (
	"context"

	gatewayProfilev1 "neupaneanish.com.np/profile/internal/protobuf/gateway/profile/v1"
	rootProfilev1 "neupaneanish.com.np/profile/internal/protobuf/root/profile/v1"
	"neupaneanish.com.np/profile/internal/utils"
)

func (s *GatewayProfileService) UpdateSocial(
	ctx context.Context,
	req *gatewayProfilev1.UpdateSocialRequest,
) (*gatewayProfilev1.UpdateSocialResponse, error) {
	serviceName := "GatewayUpdateSocial"
	userSession := utils.UserSessionContext(ctx)

	if err := updateSocial(
		ctx,
		userSession.UserID,
		userSession.UserID,
		req.GetSocial(),
		serviceName,
		s.cfg.Repository,
		s.cfg.Logger,
	); err != nil {
		return nil, err
	}

	return &gatewayProfilev1.UpdateSocialResponse{}, nil
}

func (s *RootProfileService) UpdateSocial(
	ctx context.Context,
	req *rootProfilev1.UpdateSocialRequest,
) (*rootProfilev1.UpdateSocialResponse, error) {
	serviceName := "RootUpdateSocial"
	userSession := utils.UserSessionContext(ctx)

	userID, userIDErr := utils.ParseUUID(ctx, req.GetUserId(), serviceName, s.cfg.Logger)
	if userIDErr != nil {
		return nil, userIDErr
	}

	if err := updateSocial(
		ctx,
		userID,
		userSession.UserID,
		req.GetSocial(),
		serviceName,
		s.cfg.Repository,
		s.cfg.Logger,
	); err != nil {
		return nil, err
	}

	return &rootProfilev1.UpdateSocialResponse{}, nil
}
