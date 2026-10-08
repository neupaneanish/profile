//nolint:dupl // Boilerplate delegating to deleteDatabase
package service

import (
	"context"

	gatewayProfilev1 "neupaneanish.com.np/profile/internal/protobuf/gateway/profile/v1"
	rootProfilev1 "neupaneanish.com.np/profile/internal/protobuf/root/profile/v1"
	"neupaneanish.com.np/profile/internal/utils"
)

func (s *GatewayProfileService) DeleteSocial(
	ctx context.Context,
	req *gatewayProfilev1.DeleteSocialRequest,
) (*gatewayProfilev1.DeleteSocialResponse, error) {
	if err := deleteDatabase(
		ctx,
		req.GetId(),
		"",
		"GatewayDeleteSocial",
		req.GetUpdatedAt().AsTime(),
		utils.DatabaseTableSocial,
		s.cfg.Repository,
		s.cfg.Client,
		s.cfg.Redpanda,
		s.cfg.Logger,
	); err != nil {
		return nil, err
	}

	return &gatewayProfilev1.DeleteSocialResponse{}, nil
}

func (s *RootProfileService) DeleteSocial(
	ctx context.Context,
	req *rootProfilev1.DeleteSocialRequest,
) (*rootProfilev1.DeleteSocialResponse, error) {
	if err := deleteDatabase(
		ctx,
		req.GetId(),
		req.GetUserId(),
		"RootDeleteSocial",
		req.GetUpdatedAt().AsTime(),
		utils.DatabaseTableSocial,
		s.cfg.Repository,
		s.cfg.Client,
		s.cfg.Redpanda,
		s.cfg.Logger,
	); err != nil {
		return nil, err
	}

	return &rootProfilev1.DeleteSocialResponse{}, nil
}
