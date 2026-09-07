package service

import (
	"context"

	gatewayProfilev1 "neupaneanish.com.np/profile/internal/protobuf/gateway/profile/v1"
	rootProfilev1 "neupaneanish.com.np/profile/internal/protobuf/root/profile/v1"
	"neupaneanish.com.np/profile/internal/utils"
)

func (s *GatewayProfileService) UpdateProfile(
	ctx context.Context,
	req *gatewayProfilev1.UpdateProfileRequest,
) (*gatewayProfilev1.UpdateProfileResponse, error) {
	serviceName := "GatewayUpdateProfile"
	userSession := utils.UserSessionContext(ctx)

	res, err := updateProfile(
		ctx,
		userSession.UserID,
		userSession.UserID,
		req.GetProfile(),
		req.GetUpdatedAt().AsTime(),
		s.cfg.Repository,
		s.cfg.Logger,
		serviceName,
	)
	if err != nil {
		return nil, err
	}

	return &gatewayProfilev1.UpdateProfileResponse{
		Profile: res,
	}, nil
}

func (s *RootProfileService) UpdateProfile(
	ctx context.Context,
	req *rootProfilev1.UpdateProfileRequest,
) (*rootProfilev1.UpdateProfileResponse, error) {
	serviceName := "RootUpdateProfile"
	userSession := utils.UserSessionContext(ctx)

	userID, userIDErr := parseUUID(ctx, req.GetUserId(), serviceName, s.cfg.Logger)
	if userIDErr != nil {
		return nil, userIDErr
	}

	res, err := updateProfile(
		ctx,
		userID,
		userSession.UserID,
		req.GetProfile(),
		req.GetUpdatedAt().AsTime(),
		s.cfg.Repository,
		s.cfg.Logger,
		serviceName,
	)
	if err != nil {
		return nil, err
	}

	return &rootProfilev1.UpdateProfileResponse{
		Profile: res,
	}, nil
}
