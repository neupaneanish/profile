package service

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"google.golang.org/protobuf/types/known/timestamppb"
	"neupaneanish.com.np/profile/internal/errs"
	rootProfilev1 "neupaneanish.com.np/profile/internal/protobuf/root/profile/v1"
	"neupaneanish.com.np/profile/internal/repository"
	"neupaneanish.com.np/profile/internal/utils"
)

func (s *RootProfileService) Icon(
	ctx context.Context,
	req *rootProfilev1.IconRequest,
) (*rootProfilev1.IconResponse, error) {
	serviceName := "Platform"

	id, idErr := utils.ParseUUID(ctx, req.GetId(), serviceName, s.cfg.Logger)
	if idErr != nil {
		return nil, idErr
	}

	params := &repository.IconParams{ID: id}

	row, err := s.cfg.Repository.Icon(ctx, params)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			s.cfg.Logger.WarnContext(ctx, "Icon not found", "service", serviceName)
			return nil, errs.ErrNotFound("Icon")
		}
		s.cfg.Logger.ErrorContext(ctx, "Icon query failed", "service", serviceName, "error", err)
		return nil, errs.ErrInternalServer
	}

	return &rootProfilev1.IconResponse{
		Id:         row.ID.String(),
		Name:       row.Name,
		Site:       row.Site,
		SiteSuffix: utils.StringpbValue(row.SiteSuffix),
		Url:        row.Url,
		Slug:       row.Slug,
		Color:      row.Color,
		CreatedAt:  timestamppb.New(row.CreatedAt),
		CreatedBy:  row.CreatedBy.String(),
		UpdatedAt:  timestamppb.New(row.UpdatedAt),
		UpdatedBy:  row.UpdatedBy.String(),
	}, nil
}
