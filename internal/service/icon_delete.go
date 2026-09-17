package service

import (
	"context"

	"neupaneanish.com.np/profile/internal/enum"
	rootProfilev1 "neupaneanish.com.np/profile/internal/protobuf/root/profile/v1"
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
		enum.DBTableIcon,
		s.cfg.Repository,
		s.cfg.Logger,
		s.cfg.Redpanda,
	); err != nil {
		return nil, err
	}

	return &rootProfilev1.DeleteIconResponse{}, nil
}
