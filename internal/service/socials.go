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

func (s *GatewayProfileService) Socials(
	ctx context.Context,
	_ *gatewayProfilev1.SocialsRequest,
) (*gatewayProfilev1.SocialsResponse, error) {
	serviceName := "GatewaySocials"
	userSession := utils.UserSessionContext(ctx)

	res, err := socials(ctx, userSession.UserID, s.cfg.Repository, s.cfg.Logger, serviceName)
	if err != nil {
		return nil, err
	}

	return &gatewayProfilev1.SocialsResponse{
		Socials: res,
	}, nil
}

func (s *RootProfileService) Socials(
	ctx context.Context,
	req *rootProfilev1.SocialsRequest,
) (*rootProfilev1.SocialsResponse, error) {
	serviceName := "RootSocials"
	userID, userIDErr := parseUUID(ctx, req.GetUserId(), serviceName, s.cfg.Logger)
	if userIDErr != nil {
		return nil, userIDErr
	}

	res, err := socials(ctx, userID, s.cfg.Repository, s.cfg.Logger, serviceName)
	if err != nil {
		return nil, err
	}

	return &rootProfilev1.SocialsResponse{
		Socials: res,
	}, nil
}

func socials(
	ctx context.Context,
	userID uuid.UUID,
	repo repository.Querier,
	logger *slog.Logger,
	serviceName string,
) ([]*profilev1.Socials, error) {
	params := &repository.SocialsParams{
		UserID: userID,
	}

	rows, err := repo.Socials(ctx, params)
	if err != nil {
		logger.ErrorContext(ctx, "Failed to fetch socials", "service", serviceName, "error", err)
		return nil, errs.ErrInternalServer
	}

	res := make([]*profilev1.Socials, len(rows))
	for i, s := range rows {
		res[i] = &profilev1.Socials{
			Id:         s.ID.String(),
			UserId:     s.UserID.String(),
			PlatformId: s.PlatformID.String(),
			Username:   s.Username,
			Name:       s.Name,
			Url:        s.Url,
			Logo:       s.Logo,
			Color:      s.Color,
			CreatedAt:  timestamppb.New(s.CreatedAt),
			CreatedBy:  s.CreatedBy.String(),
			UpdatedAt:  timestamppb.New(s.UpdatedAt),
			UpdatedBy:  s.UpdatedBy.String(),
		}
	}
	return res, nil
}
