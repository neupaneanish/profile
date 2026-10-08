package service

import (
	"context"

	rootProfilev1 "neupaneanish.com.np/profile/internal/protobuf/root/profile/v1"
	"neupaneanish.com.np/profile/internal/utils"
)

func (s *RootProfileService) DeleteNameserver(
	ctx context.Context,
	req *rootProfilev1.DeleteNameserverRequest,
) (*rootProfilev1.DeleteNameserverResponse, error) {
	if err := deleteDatabase(
		ctx,
		req.GetId(),
		"",
		"DeleteNameserver",
		req.GetUpdatedAt().AsTime(),
		utils.DatabaseTableNameserver,
		s.cfg.Repository,
		s.cfg.Client,
		s.cfg.Redpanda,
		s.cfg.Logger,
	); err != nil {
		return nil, err
	}

	return &rootProfilev1.DeleteNameserverResponse{}, nil
}
