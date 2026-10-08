package service

import (
	"context"
	"crypto/rand"
	"errors"
	"fmt"

	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
	"neupaneanish.com.np/profile/internal/errs"
	profilev1 "neupaneanish.com.np/profile/internal/protobuf/gateway/profile/v1"
	"neupaneanish.com.np/profile/internal/redpanda"
	"neupaneanish.com.np/profile/internal/repository"
	"neupaneanish.com.np/profile/internal/utils"
)

func (s *GatewayProfileService) CreateDomain(
	ctx context.Context,
	req *profilev1.CreateDomainRequest,
) (*profilev1.CreateDomainResponse, error) {
	serviceName := "CreateDomain"
	userSession := utils.UserSessionContext(ctx)

	hostname := req.GetHostname()

	if err := utils.ValidateHostname(hostname, false); err != nil {
		s.cfg.Logger.WarnContext(ctx, "Invalid Domain", "service", serviceName, "error", err)
		return nil, errs.ErrInvalidURL
	}

	txt := fmt.Sprintf("tuin-verify=%s", rand.Text())

	params := &repository.CreateDomainParams{
		UserID:    userSession.UserID,
		Hostname:  hostname,
		Txt:       txt,
		CreatedBy: userSession.UserID,
		UpdatedBy: userSession.UserID,
	}

	if _, err := s.cfg.Repository.CreateDomain(ctx, params); err != nil {
		if pgxErr, ok := errors.AsType[*pgconn.PgError](err); ok {
			switch pgxErr.Code {
			case pgerrcode.UniqueViolation:
				s.cfg.Logger.WarnContext(
					ctx,
					"hostname already exists",
					"service", serviceName,
					"hostname", hostname,
				)
				return nil, errs.ErrUniqueViolation("Domain")
			case pgerrcode.ForeignKeyViolation:
				s.cfg.Logger.ErrorContext(
					ctx,
					"nameserver / template not found",
					"service", serviceName,
					"hostname", hostname,
				)
				return nil, errs.ErrInternalServer
			}
		}
		s.cfg.Logger.ErrorContext(ctx, "Domain create failed", "service", serviceName, "error", err)
		return nil, errs.ErrInternalServer
	}
	redpanda.RootNotificationProduce(
		ctx,
		userSession,
		userSession.UserID,
		utils.DatabaseTableDomain,
		utils.DatabaseMethodCreate,
		serviceName,
		s.cfg.Client,
		s.cfg.Redpanda,
		s.cfg.Logger,
	)
	return &profilev1.CreateDomainResponse{}, nil
}
