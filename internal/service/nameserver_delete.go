package service

import (
	"context"

	profilev1 "neupaneanish.com.np/profile/internal/protobuf/root/profile/v1"
	"neupaneanish.com.np/profile/internal/repository"
	"neupaneanish.com.np/profile/internal/utils"
)

func (s *RootProfileService) DeleteNameserver(
	ctx context.Context,
	req *profilev1.DeleteNameserverRequest,
) (*profilev1.DeleteNameserverResponse, error) {
	serviceName := "DeleteNameserver"

	id, idErr := utils.ParseUUID(ctx, req.GetId(), serviceName, s.cfg.Logger)
	if idErr != nil {
		return nil, idErr
	}

	params := &repository.DeleteNameserverParams{
		ID:        id,
		UpdatedAt: req.GetUpdatedAt().AsTime(),
	}

	affected, err := s.cfg.Repository.DeleteNameserver(ctx, params)
	if deleteErr := deleteDB(
		ctx,
		affected,
		err,
		serviceName,
		"Nameserver",
		req.GetId(),
		s.cfg.Logger,
	); deleteErr != nil {
		return nil, deleteErr
	}

	return &profilev1.DeleteNameserverResponse{}, nil
}
