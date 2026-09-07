package service

import (
	"context"

	profilev1 "neupaneanish.com.np/profile/internal/protobuf/root/profile/v1"
	"neupaneanish.com.np/profile/internal/utils"
)

func (s *RootProfileService) UpdateNameServer(
	ctx context.Context,
	req *profilev1.UpdateNameServerRequest,
) (*profilev1.UpdateNameServerResponse, error) {
	serviceName := "UpdateNameServer"
	userSession := utils.UserSessionContext(ctx)

	if _, err := s.createUpdateNameServer(
		ctx,
		req.GetId(),
		req.GetNameserver(),
		userSession.UserID,
		serviceName,
		req.GetActive(),
		req.GetUpdatedAt().AsTime(),
	); err != nil {
		return nil, err
	}

	return &profilev1.UpdateNameServerResponse{Id: req.GetId()}, nil
}
