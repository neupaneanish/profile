package service

import (
	"context"

	profilev1 "neupaneanish.com.np/profile/internal/protobuf/root/profile/v1"
	"neupaneanish.com.np/profile/internal/repository"
)

func (s *RootProfileService) DeleteNameServer(
	ctx context.Context,
	req *profilev1.DeleteNameServerRequest,
) (*profilev1.DeleteNameServerResponse, error) {
	serviceName := "DeleteNameserver"

	id, idErr := parseUUID(ctx, req.GetId(), serviceName, s.cfg.Logger)
	if idErr != nil {
		return nil, idErr
	}

	params := &repository.DeleteNameServerParams{
		ID:        id,
		UpdatedAt: req.GetUpdatedAt().AsTime(),
	}

	cmdTag, err := s.cfg.Repository.DeleteNameServer(ctx, params)
	if deleteErr := deleteDB(
		ctx,
		cmdTag,
		err,
		serviceName,
		"Nameserver",
		req.GetId(),
		s.cfg.Logger,
	); deleteErr != nil {
		return nil, deleteErr
	}

	return &profilev1.DeleteNameServerResponse{}, nil
}
