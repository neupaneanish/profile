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

func (s *GatewayProfileService) SocialIcons(
	ctx context.Context,
	_ *gatewayProfilev1.SocialIconsRequest,
) (*gatewayProfilev1.SocialIconsResponse, error) {
	serviceName := "SocialIcons"
	userSession := utils.UserSessionContext(ctx)

	params := &repository.SocialIconsParams{UserID: userSession.UserID}

	rows, err := s.cfg.Repository.SocialIcons(ctx, params)
	if err != nil {
		s.cfg.Logger.ErrorContext(ctx, "Failed to query social icons", "service", serviceName, "error", err)
		return nil, errs.ErrInternalServer
	}

	res := make([]*profilev1.Icons, len(rows))
	for i, p := range rows {
		res[i] = &profilev1.Icons{
			Id:   p.ID.String(),
			Name: p.Name,
			Icon: p.Icon,
		}
	}

	return &gatewayProfilev1.SocialIconsResponse{Icons: res}, nil
}

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
			Id:        s.ID.String(),
			UserId:    s.UserID.String(),
			IconId:    s.IconID.String(),
			Username:  s.Username,
			Name:      s.Name,
			Site:      s.Site,
			Icon:      s.Icon,
			CreatedAt: timestamppb.New(s.CreatedAt),
			CreatedBy: s.CreatedBy.String(),
			UpdatedAt: timestamppb.New(s.UpdatedAt),
			UpdatedBy: s.UpdatedBy.String(),
		}
	}
	return res, nil
}
