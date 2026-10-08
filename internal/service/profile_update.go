package service

import (
	"context"

	"google.golang.org/protobuf/types/known/timestamppb"
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
		userSession,
		userSession.UserID,
		req.GetProfile(),
		req.GetUpdatedAt().AsTime(),
		serviceName,
		s.cfg.Repository,
		s.cfg.Client,
		s.cfg.Redpanda,
		s.cfg.Logger,
	)
	if err != nil {
		return nil, err
	}

	return &gatewayProfilev1.UpdateProfileResponse{
		Profile: &gatewayProfilev1.Profile{
			UserId:    res.UserID.String(),
			Name:      res.Name,
			Title:     res.Title,
			Dob:       timestamppb.New(res.Dob),
			UpdatedAt: timestamppb.New(res.UpdatedAt),
		},
	}, nil
}

func (s *RootProfileService) UpdateProfile(
	ctx context.Context,
	req *rootProfilev1.UpdateProfileRequest,
) (*rootProfilev1.UpdateProfileResponse, error) {
	serviceName := "RootUpdateProfile"
	userSession := utils.UserSessionContext(ctx)

	userID, userIDErr := utils.ParseUUID(ctx, req.GetUserId(), serviceName, s.cfg.Logger)
	if userIDErr != nil {
		return nil, userIDErr
	}

	res, err := updateProfile(
		ctx,
		userSession,
		userID,
		req.GetProfile(),
		req.GetUpdatedAt().AsTime(),
		serviceName,
		s.cfg.Repository,
		s.cfg.Client,
		s.cfg.Redpanda,
		s.cfg.Logger,
	)
	if err != nil {
		return nil, err
	}

	createdByUsername, updatedByUsername, usernamesErr := utils.GetUsernames(
		ctx,
		res.CreatedBy,
		res.UpdatedBy,
		userSession,
		s.cfg.Client,
		s.cfg.Logger,
	)
	if usernamesErr != nil {
		return nil, usernamesErr
	}

	return &rootProfilev1.UpdateProfileResponse{
		Profile: &rootProfilev1.Profile{
			UserId:            res.UserID.String(),
			Name:              res.Name,
			Title:             res.Title,
			Dob:               timestamppb.New(res.Dob),
			CreatedAt:         timestamppb.New(res.CreatedAt),
			CreatedBy:         res.CreatedBy.String(),
			UpdatedAt:         timestamppb.New(res.UpdatedAt),
			UpdatedBy:         res.UpdatedBy.String(),
			CreatedByUsername: createdByUsername,
			UpdatedByUsername: updatedByUsername,
		},
	}, nil
}
