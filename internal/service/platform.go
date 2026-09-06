package service

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"google.golang.org/protobuf/types/known/timestamppb"
	"neupaneanish.com.np/profile/internal/errs"
	rootProfilev1 "neupaneanish.com.np/profile/internal/protobuf/root/profile/v1"
	"neupaneanish.com.np/profile/internal/repository"
)

func (s *RootProfileService) Platform(
	ctx context.Context,
	req *rootProfilev1.PlatformRequest,
) (*rootProfilev1.PlatformResponse, error) {
	serviceName := "Platform"

	id, idErr := parseUUID(ctx, req.GetId(), serviceName, s.cfg.Logger)
	if idErr != nil {
		return nil, idErr
	}

	params := &repository.PlatformParams{ID: id}

	row, err := s.cfg.Repository.Platform(ctx, params)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			s.cfg.Logger.WarnContext(ctx, "Platform not found", "service", serviceName)
			return nil, errs.ErrNotFound
		}
		s.cfg.Logger.ErrorContext(ctx, "Platform query failed", "service", serviceName, "error", err)
		return nil, errs.ErrInternalServer
	}

	return &rootProfilev1.PlatformResponse{
		Id:            row.ID.String(),
		Name:          row.Name,
		Url:           row.Url,
		UrlSuffix:     row.UrlSuffix,
		LogoUrl:       row.LogoUrl,
		LogoUrlSuffix: row.LogoUrlSuffix,
		LogoUrlPath:   row.LogoUrlPath,
		Color:         row.Color,
		CreatedAt:     timestamppb.New(row.CreatedAt),
		CreatedBy:     row.CreatedBy.String(),
		UpdatedAt:     timestamppb.New(row.UpdatedAt),
		UpdatedBy:     row.UpdatedBy.String(),
	}, nil
}
