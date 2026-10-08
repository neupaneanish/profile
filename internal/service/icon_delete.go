package service

import (
	"context"

	rootProfilev1 "neupaneanish.com.np/profile/internal/protobuf/root/profile/v1"
	"neupaneanish.com.np/profile/internal/utils"
)

func (s *RootProfileService) DeleteIcon(
	ctx context.Context,
	req *rootProfilev1.DeleteIconRequest,
) (*rootProfilev1.DeleteIconResponse, error) {
	if err := deleteDatabase(
		ctx,
		req.GetId(),
		"",
		"RootDeleteIcon",
		req.GetUpdatedAt().AsTime(),
		utils.DatabaseTableIcon,
		s.cfg.Repository,
		s.cfg.Client,
		s.cfg.Redpanda,
		s.cfg.Logger,
	); err != nil {
		return nil, err
	}

	return &rootProfilev1.DeleteIconResponse{}, nil
}
