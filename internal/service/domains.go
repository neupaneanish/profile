package service

import (
	"context"
	"log/slog"

	"uuid"

	"neupaneanish.com.np/profile/internal/errs"
	profilev1 "neupaneanish.com.np/profile/internal/protobuf/common/profile/v1"
	gatewayProfilev1 "neupaneanish.com.np/profile/internal/protobuf/gateway/profile/v1"
	rootProfilev1 "neupaneanish.com.np/profile/internal/protobuf/root/profile/v1"
	"neupaneanish.com.np/profile/internal/repository"
	"neupaneanish.com.np/profile/internal/utils"
)

func (s *GatewayProfileService) Domains(
	ctx context.Context,
	_ *gatewayProfilev1.DomainsRequest,
) (*gatewayProfilev1.DomainsResponse, error) {
	serviceName := "GatewayDomains"
	userSession := utils.UserSessionContext(ctx)

	res, err := domains(ctx, userSession.UserID, s.cfg.Repository, s.cfg.Logger, serviceName)
	if err != nil {
		return nil, err
	}

	return &gatewayProfilev1.DomainsResponse{Domains: res}, nil
}

func (s *RootProfileService) Domains(
	ctx context.Context,
	req *rootProfilev1.DomainsRequest,
) (*rootProfilev1.DomainsResponse, error) {
	serviceName := "RootDomains"

	userID, userIDErr := utils.ParseUUID(ctx, req.GetUserId(), serviceName, s.cfg.Logger)
	if userIDErr != nil {
		return nil, userIDErr
	}

	res, err := domains(ctx, userID, s.cfg.Repository, s.cfg.Logger, serviceName)
	if err != nil {
		return nil, err
	}

	return &rootProfilev1.DomainsResponse{Domains: res}, nil
}

func domains(
	ctx context.Context,
	userID uuid.UUID,
	repo repository.Querier,
	logger *slog.Logger,
	serviceName string,
) ([]*profilev1.Domain, error) {
	params := &repository.DomainsParams{UserID: userID}

	rows, err := repo.Domains(ctx, params)
	if err != nil {
		logger.ErrorContext(ctx, "Domains query failed", "service", serviceName, "error", err)
		return nil, errs.ErrInternalServer
	}

	res := make([]*profilev1.Domain, len(rows))
	for i, row := range rows {
		res[i] = &profilev1.Domain{
			Id:       row.ID.String(),
			UserId:   row.UserID.String(),
			Hostname: row.Hostname,
			Verified: row.Verified,
		}
	}
	return res, nil
}
