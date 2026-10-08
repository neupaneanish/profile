package service

import (
	"context"
	"time"

	rootProfilev1 "neupaneanish.com.np/profile/internal/protobuf/root/profile/v1"
)

func (s *RootProfileService) CreateTemplate(
	ctx context.Context,
	req *rootProfilev1.CreateTemplateRequest,
) (*rootProfilev1.CreateTemplateResponse, error) {
	serviceName := "RootCreateTemplate"

	if err := s.createUpdateTemplate(ctx, req.GetTemplate(), "", serviceName, time.Time{}); err != nil {
		return nil, err
	}

	return &rootProfilev1.CreateTemplateResponse{}, nil
}
