package service

import (
	"context"
	"log/slog"

	"uuid"

	"google.golang.org/protobuf/types/known/timestamppb"
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

	userID, userIDErr := parseUUID(ctx, req.GetUserId(), serviceName, s.cfg.Logger)
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
) ([]*profilev1.Domains, error) {
	params := &repository.DomainsParams{UserID: userID}

	rows, err := repo.Domains(ctx, params)
	if err != nil {
		logger.ErrorContext(ctx, "Domains query failed", "service", serviceName, "error", err)
		return nil, errs.ErrInternalServer
	}

	res := make([]*profilev1.Domains, len(rows))

	for i, d := range rows {
		res[i] = &profilev1.Domains{
			Id:        d.ID.String(),
			UserId:    d.UserID.String(),
			IpType:    d.IpType,
			Ip:        d.Ip,
			Fqdn:      d.Fqdn,
			Txt:       d.Txt,
			Verified:  d.Verified,
			CreatedAt: timestamppb.New(d.CreatedAt),
			CreatedBy: d.CreatedBy.String(),
			UpdatedAt: timestamppb.New(d.UpdatedAt),
			UpdatedBy: d.UpdatedBy.String(),
		}
	}

	return res, nil
}
