//nolint:dupl // Boilerplate delegating to deleteDatabase
package service

import (
	"context"

	"neupaneanish.com.np/profile/internal/enum"
	gatewayProfilev1 "neupaneanish.com.np/profile/internal/protobuf/gateway/profile/v1"
	rootProfilev1 "neupaneanish.com.np/profile/internal/protobuf/root/profile/v1"
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
		enum.DBTableEducation,
		s.cfg.Repository,
		s.cfg.Logger,
		s.cfg.Redpanda,
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
		enum.DBTableEducation,
		s.cfg.Repository,
		s.cfg.Logger,
		s.cfg.Redpanda,
	); err != nil {
		return nil, err
	}

	return &rootProfilev1.DeleteEducationResponse{}, nil
}
