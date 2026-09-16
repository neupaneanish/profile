package service

import (
	"context"

	"neupaneanish.com.np/profile/internal/errs"
	profilev1 "neupaneanish.com.np/profile/internal/protobuf/gateway/profile/v1"
	"neupaneanish.com.np/profile/internal/repository"
	"neupaneanish.com.np/profile/internal/utils"
)

func (s *GatewayProfileService) DeleteDomain(
	ctx context.Context,
	req *profilev1.DeleteDomainRequest,
) (*profilev1.DeleteDomainResponse, error) {
	serviceName := "DeleteDomain"
	userSession := utils.UserSessionContext(ctx)

	id, idErr := utils.ParseUUID(ctx, req.GetId(), serviceName, s.cfg.Logger)
	if idErr != nil {
		return nil, idErr
	}

	fqdn := req.GetFqdn()

	if err := utils.ValidateHostname(fqdn, false); err != nil {
		s.cfg.Logger.ErrorContext(ctx, "Invalid FQDN", "service", serviceName, "error", err)
		return nil, errs.ErrConflict
	}

	params := &repository.DeleteDomainParams{
		ID:        id,
		UserID:    userSession.UserID,
		Fqdn:      fqdn,
		UpdatedAt: req.GetUpdatedAt().AsTime(),
	}

	affected, err := s.cfg.Repository.DeleteDomain(ctx, params)
	if dbErr := deleteDB(
		ctx,
		affected,
		err,
		serviceName,
		"Domain",
		req.GetId(),
		userSession.Username,
		userSession.UserID,
		userSession.UserID,
		s.cfg.Logger,
		s.cfg.Redpanda,
	); dbErr != nil {
		return nil, dbErr
	}

	cmd := s.cfg.Client.B().Hdel().Key(utils.DomainUserSessionKey).Field(fqdn).Build()

	if vkErr := s.cfg.Client.Do(ctx, cmd).Error(); vkErr != nil {
		s.cfg.Logger.ErrorContext(
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

	payload := utils.RedpandaRootEventNotificationPayload{
		ActorID:  userSession.UserID,
		Username: userSession.Username,
		UserID:   userSession.UserID,
		Message:  "domain deleted",
	}
	s.cfg.Redpanda.Produce(ctx, utils.RedpandaRootEventNotifications, serviceName, payload)

	return &profilev1.DeleteDomainResponse{}, nil
}
