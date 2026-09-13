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
		Name:      req.GetProfile().GetName(),
		Title:     req.GetProfile().GetTitle(),
		Dob:       req.GetDob().AsTime(),
		CreatedBy: userSession.UserID,
		UpdatedBy: userSession.UserID,
	}

	if err := s.cfg.Repository.CreateProfile(ctx, params); err != nil {
		if pgxErr, ok := errors.AsType[*pgconn.PgError](err); ok && pgxErr.Code == pgerrcode.UniqueViolation {
			s.cfg.Logger.WarnContext(ctx, "Profile already created", "service", serviceName)
			return nil, errs.ErrUniqueViolation("Profile")
		}
		s.cfg.Logger.ErrorContext(ctx, "Failed to insert profile", "service", serviceName, "error", err)
		return nil, errs.ErrInternalServer
	}

	return &gatewayProfilev1.CreateProfileResponse{}, nil
}
