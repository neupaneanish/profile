package service

import (
	"context"

	"neupaneanish.com.np/profile/internal/enum"
	rootProfilev1 "neupaneanish.com.np/profile/internal/protobuf/root/profile/v1"
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
		enum.DBTableNameserver,
		s.cfg.Repository,
		s.cfg.Logger,
		s.cfg.Redpanda,
	); err != nil {
		return nil, err
	}

	return &rootProfilev1.DeleteNameserverResponse{}, nil
}
