package service

import (
	"context"

	rootProfilev1 "neupaneanish.com.np/profile/internal/protobuf/root/profile/v1"
)

func (s *RootProfileService) UpdateIcon(
	ctx context.Context,
	req *rootProfilev1.UpdateIconRequest,
) (*rootProfilev1.UpdateIconResponse, error) {
	serviceName := "UpdatePlatform"

	if err := s.createUpdateIcon(
		ctx,
		req.GetId(),
		req.GetIcon(),
		serviceName,
		req.GetUpdatedAt().AsTime(),
	); err != nil {
		return nil, err
	}
	return &rootProfilev1.UpdateIconResponse{}, nil
}
