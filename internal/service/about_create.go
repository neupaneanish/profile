package service

import (
	"context"
	"errors"

	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
	"google.golang.org/protobuf/types/known/timestamppb"
	"neupaneanish.com.np/profile/internal/errs"
	profilev1 "neupaneanish.com.np/profile/internal/protobuf/common/profile/v1"
	gatewayProfilev1 "neupaneanish.com.np/profile/internal/protobuf/gateway/profile/v1"
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

	payload := utils.RedpandaRootEventNotificationPayload{
		ActorID:  userSession.UserID,
		Username: userSession.Username,
		UserID:   userSession.UserID,
		Message:  "created about",
	}
	s.cfg.Redpanda.Produce(ctx, utils.RedpandaRootEventNotifications, serviceName, payload)

	return &gatewayProfilev1.CreateAboutResponse{
		About: &profilev1.About{
			UserId:    row.UserID.String(),
			About:     row.About,
			CreatedAt: timestamppb.New(row.CreatedAt),
			CreatedBy: row.CreatedBy.String(),
			UpdatedAt: timestamppb.New(row.UpdatedAt),
			UpdatedBy: row.UpdatedBy.String(),
		},
	}, nil
}
