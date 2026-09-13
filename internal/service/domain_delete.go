package service

import (
	"context"

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

	id, idErr := parseUUID(ctx, req.GetId(), serviceName, s.cfg.Logger)
	if idErr != nil {
		return nil, idErr
	}

	params := &repository.DeleteDomainParams{
		ID:        id,
		UserID:    userSession.UserID,
		UpdatedAt: req.GetUpdatedAt().AsTime(),
	}

	affected, err := s.cfg.Repository.DeleteDomain(ctx, params)
	if dbErr := deleteDB(ctx, affected, err, serviceName, "Domain", req.GetId(), s.cfg.Logger); dbErr != nil {
		return nil, dbErr
	}

	return &profilev1.DeleteDomainResponse{}, nil
}
