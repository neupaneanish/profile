package service

import (
	"context"
	"errors"
	"time"

	"github.com/jackc/pgx/v5"
	"neupaneanish.com.np/profile/internal/enum"
	"neupaneanish.com.np/profile/internal/errs"
	profilev1 "neupaneanish.com.np/profile/internal/protobuf/gateway/profile/v1"
	"neupaneanish.com.np/profile/internal/repository"
	"neupaneanish.com.np/profile/internal/utils"
)

const (
	lookupTimeout = 2 * time.Second
)

func (s *GatewayProfileService) VerifyDomain(
	ctx context.Context,
	req *profilev1.VerifyDomainRequest,
) (*profilev1.VerifyDomainResponse, error) {
	serviceName := "VerifyDomain"
	userSession := utils.UserSessionContext(ctx)

	id, idErr := utils.ParseUUID(ctx, req.GetId(), serviceName, s.cfg.Logger)
	if idErr != nil {
		return nil, idErr
	}

	domainParams := &repository.DomainParams{
		ID:     id,
		UserID: userSession.UserID,
	}

	domain, domainErr := s.cfg.Repository.Domain(ctx, domainParams)
	if domainErr != nil {
		if errors.Is(domainErr, pgx.ErrNoRows) {
			s.cfg.Logger.WarnContext(ctx, "Domain already verified or not found")
			return nil, errs.ErrConflict
		}
		s.cfg.Logger.ErrorContext(ctx, "Domain query failed")
		return nil, errs.ErrInternalServer
	}

	if err := s.validateDomain(ctx, domain.Fqdn, domain.Txt, domain.Ip, serviceName); err != nil {
		return nil, err
	}

	params := &repository.VerifyDomainParams{
		UpdatedBy: userSession.UserID,
		ID:        id,
		UserID:    userSession.UserID,
		UpdatedAt: req.GetUpdatedAt().AsTime(),
	}

	affected, err := s.cfg.Repository.VerifyDomain(ctx, params)
	if err != nil {
		s.cfg.Logger.ErrorContext(ctx, "Verify Domain Failed", "service", serviceName, "error", err)
		return nil, errs.ErrInternalServer
	}

	if affected != 1 {
		s.cfg.Logger.WarnContext(
			ctx,
			"Record not found or concurrent modification or already verified",
			"service", serviceName,
			"id", id,
		)
		return nil, errs.ErrConflict
	}

	cmd := s.cfg.Client.B().Hset().
		Key(utils.DomainUserSessionKey).
		FieldValue().
		FieldValue(domain.Fqdn, userSession.UserID.String()).
		Build()

	if vkErr := s.cfg.Client.Do(ctx, cmd).Error(); vkErr != nil {
		s.cfg.Logger.ErrorContext(
			ctx,
			"Failed to set domain user session in valkey",
			"service", serviceName,
			"domain", domain.Fqdn,
			"error", vkErr,
		)
	}

	payload := utils.RedpandaRootEventNotificationPayload{
		ActorID:  userSession.UserID,
		Username: userSession.Username,
		UserID:   userSession.UserID,
		Method:   enum.DBMethodUpdate,
		Table:    enum.DBTableDomain,
	}
	s.cfg.Redpanda.Produce(ctx, utils.RedpandaRootDatabaseEventNotifications, serviceName, payload)

	return &profilev1.VerifyDomainResponse{}, nil
}
