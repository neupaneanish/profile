package service

import (
	"context"
	"errors"

	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
	"google.golang.org/protobuf/types/known/timestamppb"
	"neupaneanish.com.np/profile/internal/errs"
	gatewayProfilev1 "neupaneanish.com.np/profile/internal/protobuf/gateway/profile/v1"
	"neupaneanish.com.np/profile/internal/redpanda"
	"neupaneanish.com.np/profile/internal/repository"
	"neupaneanish.com.np/profile/internal/utils"
)

func (s *GatewayProfileService) CreateAbout(
	ctx context.Context,
	req *gatewayProfilev1.CreateAboutRequest,
) (*gatewayProfilev1.CreateAboutResponse, error) {
	serviceName := "CreateAbout"
	userSession := utils.UserSessionContext(ctx)

	params := &repository.CreateAboutParams{
		UserID:    userSession.UserID,
		About:     req.GetAbout(),
		CreatedBy: userSession.UserID,
		UpdatedBy: userSession.UserID,
	}

	row, err := s.cfg.Repository.CreateAbout(ctx, params)
	if err != nil {
		if pgxErr, ok := errors.AsType[*pgconn.PgError](err); ok && pgxErr.Code == pgerrcode.UniqueViolation {
			s.cfg.Logger.WarnContext(ctx, "About already created", "service", serviceName)
			return nil, errs.ErrUniqueViolation("About")
		}
		s.cfg.Logger.ErrorContext(ctx, "Failed to create about", "service", serviceName, "error", err)
		return nil, errs.ErrInternalServer
	}

	redpanda.RootNotificationProduce(
		ctx,
		userSession,
		userSession.UserID,
		utils.DatabaseTableAbout,
		utils.DatabaseMethodCreate,
		serviceName,
		s.cfg.Client,
		s.cfg.Redpanda,
		s.cfg.Logger,
	)

	return &gatewayProfilev1.CreateAboutResponse{
		About: &gatewayProfilev1.About{
			UserId:    row.UserID.String(),
			About:     row.About,
			UpdatedAt: timestamppb.New(row.UpdatedAt),
		},
	}, nil
}
