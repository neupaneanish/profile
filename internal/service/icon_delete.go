package service

import (
	"context"

	rootProfilev1 "neupaneanish.com.np/profile/internal/protobuf/root/profile/v1"
	"neupaneanish.com.np/profile/internal/repository"
	"neupaneanish.com.np/profile/internal/utils"
)

func (s *RootProfileService) DeleteIcon(
	ctx context.Context,
	req *rootProfilev1.DeleteIconRequest,
) (*rootProfilev1.DeleteIconResponse, error) {
	serviceName := "DeleteIcon"
	userSession := utils.UserSessionContext(ctx)

	id, idErr := utils.ParseUUID(ctx, req.GetId(), serviceName, s.cfg.Logger)
	if idErr != nil {
		return nil, idErr
	}

	params := &repository.DeleteIconParams{ID: id, UpdatedAt: req.GetUpdatedAt().AsTime()}

	affected, err := s.cfg.Repository.DeleteIcon(ctx, params)
	if deleteErr := deleteDB(
		ctx,
		affected,
		err,
		serviceName,
		"Icon",
		req.GetId(),
		userSession.Username,
		userSession.UserID,
		userSession.UserID,
		s.cfg.Logger,
		s.cfg.Redpanda,
	); deleteErr != nil {
		return nil, deleteErr
	}

	return &rootProfilev1.DeleteIconResponse{}, nil
}
