package service

import (
	"context"

	"neupaneanish.com.np/profile/internal/errs"
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
	if err != nil {
		s.cfg.Logger.ErrorContext(ctx, "Delete Platform Failed", "service", serviceName, "error", err)
		return nil, errs.ErrInternalServer
	}

	if cmdTag.RowsAffected() == 0 {
		s.cfg.Logger.WarnContext(
			ctx,
			"Platform record not found or concurrent modification",
			"service", serviceName,
			"id", id.String(),
		)
		return nil, errs.ErrConflict
	}

	return &rootProfilev1.DeletePlatformResponse{}, nil
}
