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

func (s *RootProfileService) Nameserver(
	ctx context.Context,
	req *rootProfilev1.NameserverRequest,
) (*rootProfilev1.NameserverResponse, error) {
	serviceName := "RootNameserver"
	userSession := utils.UserSessionContext(ctx)

	id, idErr := utils.ParseUUID(ctx, req.GetId(), serviceName, s.cfg.Logger)
	if idErr != nil {
		return nil, idErr
	}

	params := &repository.NameserverParams{ID: id}

	row, err := s.cfg.Repository.Nameserver(ctx, params)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			s.cfg.Logger.WarnContext(ctx, "Nameserver not found", "service", serviceName, "error", err)
			return nil, errs.ErrNotFound("Nameserver")
		}
		s.cfg.Logger.ErrorContext(ctx, "Nameserver query failed", "service", serviceName, "error", err)
		return nil, errs.ErrInternalServer
	}

	createdUsername, updatedUsername, usernamesErr := utils.GetUsernames(
		ctx,
		row.CreatedBy,
		row.UpdatedBy,
		userSession,
		s.cfg.Client,
		s.cfg.Logger,
	)
	if usernamesErr != nil {
		return nil, usernamesErr
	}

	return &rootProfilev1.NameserverResponse{
		Id:                row.ID.String(),
		Cname:             row.Cname,
		Hostname:          row.Hostname,
		CreatedAt:         timestamppb.New(row.CreatedAt),
		CreatedBy:         row.CreatedBy.String(),
		UpdatedAt:         timestamppb.New(row.UpdatedAt),
		UpdatedBy:         row.UpdatedBy.String(),
		CreatedByUsername: createdUsername,
		UpdatedByUsername: updatedUsername,
	}, nil
}
