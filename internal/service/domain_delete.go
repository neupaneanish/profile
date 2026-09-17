package service

import (
	"context"
	"log/slog"
	"time"
	"uuid"

	"github.com/valkey-io/valkey-go"
	"neupaneanish.com.np/profile/internal/config"
	"neupaneanish.com.np/profile/internal/enum"
	"neupaneanish.com.np/profile/internal/errs"
	gatewayProfilev1 "neupaneanish.com.np/profile/internal/protobuf/gateway/profile/v1"
	rootProfilev1 "neupaneanish.com.np/profile/internal/protobuf/root/profile/v1"
	"neupaneanish.com.np/profile/internal/repository"
	"neupaneanish.com.np/profile/internal/utils"
)

func (s *GatewayProfileService) DeleteDomain(
	ctx context.Context,
	req *gatewayProfilev1.DeleteDomainRequest,
) (*gatewayProfilev1.DeleteDomainResponse, error) {
	serviceName := "gatewayDeleteDomain"
	userSession := utils.UserSessionContext(ctx)

	if err := deleteDomain(
		ctx,
		req.GetId(),
		userSession.Username,
		req.GetFqdn(),
		serviceName,
		userSession.UserID,
		userSession.UserID,
		req.GetUpdatedAt().AsTime(),
		s.cfg.Repository,
		s.cfg.Client,
		s.cfg.Logger,
		s.cfg.Redpanda,
	); err != nil {
		return nil, err
	}

	return &gatewayProfilev1.DeleteDomainResponse{}, nil
}

func (s *RootProfileService) DeleteDomain(
	ctx context.Context,
	req *rootProfilev1.DeleteDomainRequest,
) (*rootProfilev1.DeleteDomainResponse, error) {
	serviceName := "RootDeleteDomain"
	userSession := utils.UserSessionContext(ctx)

	userID, userIDErr := utils.ParseUUID(ctx, req.GetUserId(), serviceName, s.cfg.Logger)
	if userIDErr != nil {
		return nil, userIDErr
	}

	if err := deleteDomain(
		ctx,
		req.GetId(),
		userSession.Username,
		req.GetFqdn(),
		serviceName,
		userID,
		userSession.UserID,
		req.GetUpdatedAt().AsTime(),
		s.cfg.Repository,
		s.cfg.Client,
		s.cfg.Logger,
		s.cfg.Redpanda,
	); err != nil {
		return nil, err
	}

	return &rootProfilev1.DeleteDomainResponse{}, nil
}

func deleteDomain(
	ctx context.Context,
	id, username, fqdn, serviceName string,
	userID, updatedBy uuid.UUID,
	updatedAt time.Time,
	repo repository.Querier,
	client valkey.Client,
	logger *slog.Logger,
	redpanda *config.Redpanda,
) error {
	idx, idErr := utils.ParseUUID(ctx, id, serviceName, logger)
	if idErr != nil {
		return idErr
	}

	if err := utils.ValidateHostname(fqdn, false); err != nil {
		logger.ErrorContext(ctx, "Invalid FQDN", "service", serviceName, "error", err)
		return errs.ErrConflict
	}

	params := &repository.DeleteDomainParams{
		ID:        idx,
		UserID:    userID,
		Fqdn:      fqdn,
		UpdatedAt: updatedAt,
	}

	affected, err := repo.DeleteDomain(ctx, params)
	if dbErr := deleteDB(
		ctx,
		affected,
		err,
		serviceName,
		id,
		username,
		enum.DBTableDomain,
		updatedBy,
		userID,
		logger,
		redpanda,
	); dbErr != nil {
		return dbErr
	}

	cmd := client.B().Hdel().Key(utils.DomainUserSessionKey).Field(fqdn).Build()

	if vkErr := client.Do(ctx, cmd).Error(); vkErr != nil {
		logger.ErrorContext(
			ctx,
			"Failed to delete domain user session from valkey",
			"service",
			serviceName,
			"domain",
			fqdn,
			"error",
			vkErr,
		)
	}
	return nil
}
