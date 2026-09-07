package service

import (
	"context"

	rootProfilev1 "neupaneanish.com.np/profile/internal/protobuf/root/profile/v1"
	"neupaneanish.com.np/profile/internal/repository"
)

func (s *RootProfileService) DeletePlatform(
	ctx context.Context,
	req *rootProfilev1.DeletePlatformRequest,
) (*rootProfilev1.DeletePlatformResponse, error) {
	serviceName := "DeletePlatform"

	id, idErr := parseUUID(ctx, req.GetId(), serviceName, s.cfg.Logger)
	if idErr != nil {
		return nil, idErr
	}

	params := &repository.DeletePlatformParams{ID: id, UpdatedAt: req.GetUpdatedAt().AsTime()}

	cmdTag, err := s.cfg.Repository.DeletePlatform(ctx, params)
	if deleteErr := deleteDB(
		ctx,
		cmdTag,
		err,
		serviceName,
		"Platform",
		req.GetId(),
		s.cfg.Logger,
	); deleteErr != nil {
		return nil, deleteErr
	}

	return &rootProfilev1.DeletePlatformResponse{}, nil
}
