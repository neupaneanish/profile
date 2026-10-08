//nolint:dupl // Boilerplate delegating to deleteDatabase
package service

import (
	"context"

	gatewayProfilev1 "neupaneanish.com.np/profile/internal/protobuf/gateway/profile/v1"
	rootProfilev1 "neupaneanish.com.np/profile/internal/protobuf/root/profile/v1"
	"neupaneanish.com.np/profile/internal/utils"
)

func (s *GatewayProfileService) DeleteEducation(
	ctx context.Context,
	req *gatewayProfilev1.DeleteEducationRequest,
) (*gatewayProfilev1.DeleteEducationResponse, error) {
	if err := deleteDatabase(
		ctx,
		req.GetId(),
		"",
		"GatewayDeleteEducation",
		req.GetUpdatedAt().AsTime(),
		utils.DatabaseTableEducation,
		s.cfg.Repository,
		s.cfg.Client,
		s.cfg.Redpanda,
		s.cfg.Logger,
	); err != nil {
		return nil, err
	}

	return &gatewayProfilev1.DeleteEducationResponse{}, nil
}

func (s *RootProfileService) DeleteEducation(
	ctx context.Context,
	req *rootProfilev1.DeleteEducationRequest,
) (*rootProfilev1.DeleteEducationResponse, error) {
	if err := deleteDatabase(
		ctx,
		req.GetId(),
		req.GetUserId(),
		"RootDeleteEducation",
		req.GetUpdatedAt().AsTime(),
		utils.DatabaseTableEducation,
		s.cfg.Repository,
		s.cfg.Client,
		s.cfg.Redpanda,
		s.cfg.Logger,
	); err != nil {
		return nil, err
	}

	return &rootProfilev1.DeleteEducationResponse{}, nil
}
