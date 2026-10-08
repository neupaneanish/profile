package service

import (
	"context"
	"errors"

	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
	"neupaneanish.com.np/profile/internal/errs"
	profilev1 "neupaneanish.com.np/profile/internal/protobuf/root/profile/v1"
	"neupaneanish.com.np/profile/internal/redpanda"
	"neupaneanish.com.np/profile/internal/repository"
	"neupaneanish.com.np/profile/internal/utils"
)

func (s *RootProfileService) CreateNameserver(
	ctx context.Context,
	req *profilev1.CreateNameserverRequest,
) (*profilev1.CreateNameserverResponse, error) {
	serviceName := "CreateNameserver"
	userSession := utils.UserSessionContext(ctx)

	hostname := req.GetHostname()

	if err := utils.ValidateHostname(hostname, false); err != nil {
		s.cfg.Logger.WarnContext(ctx, "Invalid hostname", "service", serviceName, "error", err)
		return nil, errs.ErrInvalidHostname
	}

	params := &repository.CreateNameserverParams{
		Hostname:  hostname,
		Cname:     req.GetCname(),
		CreatedBy: userSession.UserID,
		UpdatedBy: userSession.UserID,
	}

	if _, err := s.cfg.Repository.CreateNameserver(ctx, params); err != nil {
		if pgErr, ok := errors.AsType[*pgconn.PgError](err); ok && pgErr.Code == pgerrcode.UniqueViolation {
			s.cfg.Logger.WarnContext(
				ctx,
				"NameServer Already Exists",
				"service", serviceName,
			)
			return nil, errs.ErrUniqueViolation("Nameserver")
		}
		s.cfg.Logger.ErrorContext(
			ctx,
			"Failed to create nameserver",
			"service", serviceName,
			"error", err,
		)
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

	return &profilev1.CreateNameserverResponse{}, nil
}
