//nolint:dupl // Boilerplate delegating to deleteDatabase
package service

import (
	"context"

	"neupaneanish.com.np/profile/internal/enum"
	gatewayProfilev1 "neupaneanish.com.np/profile/internal/protobuf/gateway/profile/v1"
	rootProfilev1 "neupaneanish.com.np/profile/internal/protobuf/root/profile/v1"
)

func (s *GatewayProfileService) DeleteExperience(
	ctx context.Context,
	req *gatewayProfilev1.DeleteExperienceRequest,
) (*gatewayProfilev1.DeleteExperienceResponse, error) {
	if err := deleteDatabase(
		ctx,
		req.GetId(),
		"",
		"GatewayDeleteExperience",
		req.GetUpdatedAt().AsTime(),
		enum.DBTableExperience,
		s.cfg.Repository,
		s.cfg.Logger,
		s.cfg.Redpanda,
	); err != nil {
		return nil, err
	}

	return &gatewayProfilev1.DeleteExperienceResponse{}, nil
}

func (s *RootProfileService) DeleteExperience(
	ctx context.Context,
	req *rootProfilev1.DeleteExperienceRequest,
) (*rootProfilev1.DeleteExperienceResponse, error) {
	if err := deleteDatabase(
		ctx,
		req.GetId(),
		req.GetUserId(),
		"RootDeleteExperience",
		req.GetUpdatedAt().AsTime(),
		enum.DBTableExperience,
		s.cfg.Repository,
		s.cfg.Logger,
		s.cfg.Redpanda,
	); err != nil {
		return nil, err
	}

	return &rootProfilev1.DeleteExperienceResponse{}, nil
}
