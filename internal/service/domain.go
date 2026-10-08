package service

import (
	"context"
	"errors"
	"log/slog"
	"uuid"

	"github.com/jackc/pgx/v5"
	"google.golang.org/protobuf/types/known/timestamppb"
	"neupaneanish.com.np/profile/internal/errs"
	gatewayProfilev1 "neupaneanish.com.np/profile/internal/protobuf/gateway/profile/v1"
	rootProfilev1 "neupaneanish.com.np/profile/internal/protobuf/root/profile/v1"
	"neupaneanish.com.np/profile/internal/repository"
	"neupaneanish.com.np/profile/internal/utils"
)

func (s *GatewayProfileService) Domain(
	ctx context.Context,
	req *gatewayProfilev1.DomainRequest,
) (*gatewayProfilev1.DomainResponse, error) {
	serviceName := "GatewayDomain"
	userSession := utils.UserSessionContext(ctx)

	row, err := domain(ctx, req.GetId(), serviceName, userSession.UserID, s.cfg.Repository, s.cfg.Logger)
	if err != nil {
		return nil, err
	}

	return &gatewayProfilev1.DomainResponse{
		Id:           row.ID.String(),
		UserId:       row.UserID.String(),
		TemplateId:   row.TemplateID.String(),
		Template:     row.Template,
		TemplateIcon: row.TemplateIcon,
		Hostname:     row.Hostname,
		Txt:          row.Txt,
		Nameserver:   row.Nameserver,
		Verified:     row.Verified,
		UpdatedAt:    timestamppb.New(row.UpdatedAt),
	}, nil
}

func (s *RootProfileService) Domain(
	ctx context.Context,
	req *rootProfilev1.DomainRequest,
) (*rootProfilev1.DomainResponse, error) {
	serviceName := "RootDomain"
	userSession := utils.UserSessionContext(ctx)

	userID, userIDErr := utils.ParseUUID(ctx, req.GetUserId(), serviceName, s.cfg.Logger)
	if userIDErr != nil {
		return nil, userIDErr
	}

	row, err := domain(ctx, req.GetId(), serviceName, userID, s.cfg.Repository, s.cfg.Logger)
	if err != nil {
		return nil, err
	}

	createdByUsername, updatedByUsername, usernamesErr := utils.GetUsernames(
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

	return &rootProfilev1.DomainResponse{
		Id:                row.ID.String(),
		UserId:            row.UserID.String(),
		TemplateId:        row.TemplateID.String(),
		Template:          row.Template,
		TemplateIcon:      row.TemplateIcon,
		Hostname:          row.Hostname,
		Txt:               row.Txt,
		Nameserver:        row.Nameserver,
		NameserverId:      row.NameserverID.String(),
		VerifiedAt:        utils.TimestamppbValue(row.VerifiedAt),
		Verified:          row.Verified,
		CreatedAt:         timestamppb.New(row.CreatedAt),
		CreatedBy:         row.CreatedBy.String(),
		UpdatedAt:         timestamppb.New(row.UpdatedAt),
		UpdatedBy:         row.UpdatedBy.String(),
		CreatedByUsername: createdByUsername,
		UpdatedByUsername: updatedByUsername,
	}, nil
}

func domain(
	ctx context.Context,
	id, serviceName string,
	userID uuid.UUID,
	querier repository.Querier,
	logger *slog.Logger,
) (*repository.DomainRow, error) {
	idx, idxErr := utils.ParseUUID(ctx, id, serviceName, logger)
	if idxErr != nil {
		return nil, idxErr
	}

	params := &repository.DomainParams{
		ID:     idx,
		UserID: userID,
	}

	row, err := querier.Domain(ctx, params)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			logger.WarnContext(ctx, "Not found", "service", serviceName, "error", err)
			return nil, errs.ErrNotFound("Domain")
		}
		logger.ErrorContext(ctx, "Failed to query domain", "service", serviceName, "error", err)
		return nil, errs.ErrInternalServer
	}

	return row, nil
}
