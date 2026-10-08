package service

import (
	"context"

	rootProfilev1 "neupaneanish.com.np/profile/internal/protobuf/root/profile/v1"
	"neupaneanish.com.np/profile/internal/utils"
)

func (s *RootProfileService) DeleteTemplate(
	ctx context.Context,
	req *rootProfilev1.DeleteTemplateRequest,
) (*rootProfilev1.DeleteTemplateResponse, error) {
	serviceName := "RootDeleteTemplate"
	if err := deleteDatabase(
		ctx,
		req.GetId(),
		"",
		"RootDeleteTemplate",
		req.GetUpdatedAt().AsTime(),
		utils.DatabaseTableTemplate,
		s.cfg.Repository,
		s.cfg.Client,
		s.cfg.Redpanda,
		s.cfg.Logger,
	); err != nil {
		return nil, err
	}

	cmd := s.cfg.Client.B().Hdel().Key(req.GetId()).Field("name").Build()

	if vkErr := s.cfg.Client.Do(ctx, cmd).Error(); vkErr != nil {
		s.cfg.Logger.ErrorContext(
			ctx,
			"Failed to delete template payload from valkey",
			"service", serviceName,
			"error", vkErr,
		)
	}

	return &rootProfilev1.DeleteTemplateResponse{}, nil
}
