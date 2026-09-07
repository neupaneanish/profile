package service

import (
	"context"
	"time"

	profilev1 "neupaneanish.com.np/profile/internal/protobuf/root/profile/v1"
	"neupaneanish.com.np/profile/internal/utils"
)

func (s *RootProfileService) CreateNameServer(
	ctx context.Context,
	req *profilev1.CreateNameServerRequest,
) (*profilev1.CreateNameServerResponse, error) {
	serviceName := "CreateNameServer"
	userSession := utils.UserSessionContext(ctx)

	id, err := s.createUpdateNameServer(
		ctx,
		"",
		req.GetNameserver(),
		userSession.UserID,
		serviceName,
		false,
		time.Time{},
	)
	if err != nil {
		return nil, err
	}

	return &profilev1.CreateNameServerResponse{Id: id.String()}, nil
}
