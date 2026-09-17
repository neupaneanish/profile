package service

import (
	"context"
	"crypto/rand"
	"errors"
	"fmt"

	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"neupaneanish.com.np/profile/internal/enum"
	"neupaneanish.com.np/profile/internal/errs"
	profilev1 "neupaneanish.com.np/profile/internal/protobuf/gateway/profile/v1"
	"neupaneanish.com.np/profile/internal/repository"
	"neupaneanish.com.np/profile/internal/utils"
)

func (s *GatewayProfileService) CreateDomain(
	ctx context.Context,
	req *profilev1.CreateDomainRequest,
) (*profilev1.CreateDomainResponse, error) {
	serviceName := "CreateDomain"
	userSession := utils.UserSessionContext(ctx)

	if err := utils.ValidateHostname(req.GetUrl(), false); err != nil {
		s.cfg.Logger.WarnContext(ctx, "Invalid Domain", "service", serviceName, "error", err)
		return nil, errs.ErrInvalidURL
	}

	nameserverID, nameserverIDErr := s.cfg.Repository.Nameserver(ctx)
	if nameserverIDErr != nil {
		if errors.Is(nameserverIDErr, pgx.ErrNoRows) {
			s.cfg.Logger.ErrorContext(ctx, "No nameserver found", "service", serviceName)
			return nil, errs.ErrInternalServer
		}
		s.cfg.Logger.ErrorContext(
			ctx,
			"Domain Nameserver query failed",
			"service",
			serviceName,
			"error",
			nameserverIDErr,
		)
		return nil, errs.ErrInternalServer
	}

	txt := fmt.Sprintf("tuin-verify=%s", rand.Text())

	params := &repository.CreateDomainParams{
		UserID:       userSession.UserID,
		NameserverID: nameserverID,
		Fqdn:         req.GetUrl(),
		Txt:          txt,
		CreatedBy:    userSession.UserID,
		UpdatedBy:    userSession.UserID,
	}

	if _, err := s.cfg.Repository.CreateDomain(ctx, params); err != nil {
		if pgxErr, ok := errors.AsType[*pgconn.PgError](err); ok && pgxErr.Code == pgerrcode.UniqueViolation {
			s.cfg.Logger.WarnContext(
				ctx,
				"Domain already exists",
				"service", serviceName,
				"domain", req.GetUrl(),
			)
			return nil, errs.ErrUniqueViolation("Domain")
		}
		s.cfg.Logger.ErrorContext(ctx, "Domain create failed", "service", serviceName, "error", err)
		return nil, errs.ErrInternalServer
	}
	payload := utils.RedpandaRootEventNotificationPayload{
		ActorID:  userSession.UserID,
		Username: userSession.Username,
		UserID:   userSession.UserID,
		Method:   enum.DBMethodCreate,
		Table:    enum.DBTableDomain,
	}

	s.cfg.Redpanda.Produce(ctx, utils.RedpandaRootDatabaseEventNotifications, serviceName, payload)
	return &profilev1.CreateDomainResponse{}, nil
}
