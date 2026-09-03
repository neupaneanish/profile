package service

import (
	"context"
	"errors"

	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
	"neupaneanish.com.np/profile/internal/errs"
	gatewayProfilev1 "neupaneanish.com.np/profile/internal/protobuf/gateway/profile/v1"
	"neupaneanish.com.np/profile/internal/repository"
	"neupaneanish.com.np/profile/internal/utils"
)

func (s *GatewayProfileService) CreateProfile(
	ctx context.Context,
	req *gatewayProfilev1.CreateProfileRequest,
) (*gatewayProfilev1.CreateProfileResponse, error) {
	serviceName := "CreateProfile"
	userSession := utils.UserSessionContext(ctx)

	params := &repository.CreateProfileParams{
		UserID:    userSession.UserID,
		Name:      req.GetName(),
		Title:     req.GetTitle(),
		Dob:       req.GetDob().AsTime(),
		CreatedBy: userSession.UserID,
		UpdatedBy: userSession.UserID,
	}

	_, cmdTagErr := s.cfg.Repository.CreateProfile(ctx, params)
	if cmdTagErr != nil {
		if pgxErr, ok := errors.AsType[*pgconn.PgError](cmdTagErr); ok && pgxErr.Code == pgerrcode.UniqueViolation {
			s.cfg.Logger.WarnContext(ctx, "Profile already created", "service", serviceName)
			return nil, errs.ErrProfileAlreadyExists
		}
		s.cfg.Logger.ErrorContext(ctx, "Failed to insert profile", "service", serviceName, "error", cmdTagErr)
		return nil, errs.ErrInternalServer
	}

	return &gatewayProfilev1.CreateProfileResponse{}, nil
}
