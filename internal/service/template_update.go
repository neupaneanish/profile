package service

import (
	"context"

	rootProfilev1 "neupaneanish.com.np/profile/internal/protobuf/root/profile/v1"
)

func (s *RootProfileService) UpdateTemplate(
	ctx context.Context,
	req *rootProfilev1.UpdateTemplateRequest,
) (*rootProfilev1.UpdateTemplateResponse, error) {
	serviceName := "RootUpdateTemplate"

	if err := s.createUpdateTemplate(
		ctx,
		req.GetTemplate(),
		req.GetId(),
		serviceName,
		req.GetUpdatedAt().AsTime(),
	); err != nil {
		return nil, err
	}

	return &rootProfilev1.UpdateTemplateResponse{}, nil
}
