package service

import (
	"context"
	"time"

	rootProfilev1 "neupaneanish.com.np/profile/internal/protobuf/root/profile/v1"
)

func (s *RootProfileService) CreateIcon(
	ctx context.Context,
	req *rootProfilev1.CreateIconRequest,
) (*rootProfilev1.CreateIconResponse, error) {
	serviceName := "CreateIcon"

	if err := s.createUpdateIcon(
		ctx,
		"",
		req.GetIcon(),
		serviceName,
		time.Time{},
	); err != nil {
		return nil, err
	}

	return &rootProfilev1.CreateIconResponse{}, nil
}
