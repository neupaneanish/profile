package service

import (
	"context"
	"errors"
	"time"
	"uuid"

	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"neupaneanish.com.np/profile/internal/errs"
	profilev1 "neupaneanish.com.np/profile/internal/protobuf/gateway/profile/v1"
	"neupaneanish.com.np/profile/internal/redpanda"
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

	if err := s.validateHostname(ctx, req.GetHostname(), req.GetTxt(), serviceName); err != nil {
		return nil, err
	}

	params := &repository.VerifyDomainParams{
		UpdatedBy: userSession.UserID,
		ID:        id,
		UserID:    userSession.UserID,
		UpdatedAt: req.GetUpdatedAt().AsTime(),
		Txt:       req.GetTxt(),
	}

	row, err := s.cfg.Repository.VerifyDomain(ctx, params)
	if vErr := s.updateDomainError(ctx, err, serviceName); vErr != nil {
		return nil, vErr
	}

	s.updateDomainValkey(ctx, row.UserID, row.TemplateID, userSession, row.Hostname, serviceName)

	return &profilev1.VerifyDomainResponse{}, nil
}

func (s *GatewayProfileService) UpdateDomainTemplate(
	ctx context.Context,
	req *profilev1.UpdateDomainTemplateRequest,
) (*profilev1.UpdateDomainTemplateResponse, error) {
	serviceName := "GatewayDomainTemplateUpdate"
	userSession := utils.UserSessionContext(ctx)

	id, idErr := utils.ParseUUID(ctx, req.GetId(), serviceName, s.cfg.Logger)
	if idErr != nil {
		return nil, idErr
	}

	templateID, templateIDErr := utils.ParseUUID(ctx, req.GetTemplateId(), serviceName, s.cfg.Logger)
	if templateIDErr != nil {
		return nil, templateIDErr
	}

	params := &repository.UpdateDomainTemplateParams{
		TemplateID: templateID,
		UpdatedBy:  userSession.UserID,
		ID:         id,
		UserID:     userSession.UserID,
		UpdatedAt:  req.GetUpdatedAt().AsTime(),
	}

	row, err := s.cfg.Repository.UpdateDomainTemplate(ctx, params)
	if vErr := s.updateDomainError(ctx, err, serviceName); vErr != nil {
		return nil, vErr
	}

	s.updateDomainValkey(ctx, row.UserID, row.TemplateID, userSession, row.Hostname, serviceName)

	return &profilev1.UpdateDomainTemplateResponse{}, nil
}

func (s *GatewayProfileService) updateDomainError(ctx context.Context, err error, serviceName string) error {
	if err != nil {
		if pgxErr, ok := errors.AsType[*pgconn.PgError](err); ok && pgxErr.Code == pgerrcode.ForeignKeyViolation {
			s.cfg.Logger.ErrorContext(
				ctx,
				"template not found",
				"service", serviceName,
				"error", err,
			)
			return errs.ErrForeignKeyViolation("Template ID")
		}
		if errors.Is(err, pgx.ErrNoRows) {
			s.cfg.Logger.WarnContext(ctx,
				"Concurrent domain update detected",
				"service", serviceName,
				"error", err)
			return errs.ErrConflict
		}
		s.cfg.Logger.ErrorContext(ctx, "Verify Domain Failed", "service", serviceName, "error", err)
		return errs.ErrInternalServer
	}
	return nil
}

func (s *GatewayProfileService) updateDomainValkey(
	ctx context.Context,
	userID, templateID uuid.UUID,
	userSession *utils.UserSession,
	hostname, serviceName string,
) {
	cmd := s.cfg.Client.B().Hset().
		Key(hostname).
		FieldValue().
		FieldValue("user_id", userID.String()).
		FieldValue("template_id", templateID.String()).
		Build()

	if vkErr := s.cfg.Client.Do(ctx, cmd).Error(); vkErr != nil {
		s.cfg.Logger.ErrorContext(
			ctx,
			"Failed to set domain template payload in valkey",
			"service", serviceName,
			"hostname", hostname,
			"error", vkErr,
		)
	}

	redpanda.RootNotificationProduce(
		ctx,
		userSession,
		userSession.UserID,
		utils.DatabaseTableDomain,
		utils.DatabaseMethodUpdate,
		serviceName,
		s.cfg.Client,
		s.cfg.Redpanda,
		s.cfg.Logger,
	)
}
