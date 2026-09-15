package service

import (
	"context"

	"neupaneanish.com.np/profile/internal/errs"
	profilev1 "neupaneanish.com.np/profile/internal/protobuf/gateway/profile/v1"
	"neupaneanish.com.np/profile/internal/redis"
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
	if dbErr := deleteDB(ctx, affected, err, serviceName, "Domain", req.GetId(), s.cfg.Logger); dbErr != nil {
		return nil, dbErr
	}

	if vkErr := redis.HDelete[utils.DomainUser](ctx, utils.DomainUserSessionKey, fqdn, s.cfg.Client); vkErr != nil {
		s.cfg.Logger.ErrorContext(ctx, "Failed to delete fqdn cache", "service", serviceName, "error", vkErr)
	}

	return &profilev1.DeleteDomainResponse{}, nil
}
