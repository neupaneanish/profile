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

func (s *GatewayProfileService) Experiences(
	ctx context.Context,
	_ *gatewayProfilev1.ExperiencesRequest,
) (*gatewayProfilev1.ExperiencesResponse, error) {
	serviceName := "GatewayExperiences"

	userSession := utils.UserSessionContext(ctx)

	res, err := experiences(ctx, userSession.UserID, s.cfg.Repository, s.cfg.Logger, serviceName)
	if err != nil {
		return nil, err
	}

	return &gatewayProfilev1.ExperiencesResponse{Experiences: res}, nil
}

func (s *RootProfileService) Experiences(
	ctx context.Context,
	req *rootProfilev1.ExperiencesRequest,
) (*rootProfilev1.ExperiencesResponse, error) {
	serviceName := "GatewayExperiences"

	userID, userIDErr := parseUUID(ctx, req.GetUserId(), serviceName, s.cfg.Logger)
	if userIDErr != nil {
		return nil, userIDErr
	}

	res, err := experiences(ctx, userID, s.cfg.Repository, s.cfg.Logger, serviceName)
	if err != nil {
		return nil, err
	}

	return &rootProfilev1.ExperiencesResponse{Experiences: res}, nil
}

func experiences(
	ctx context.Context,
	userID uuid.UUID,
	repo repository.Querier,
	logger *slog.Logger,
	serviceName string,
) ([]*profilev1.Experiences, error) {
	params := &repository.ExperiencesParams{UserID: userID}

	rows, err := repo.Experiences(ctx, params)
	if err != nil {
		logger.ErrorContext(ctx, "Failed to fetch experiences", "service", serviceName, "error", err)
		return nil, errs.ErrInternalServer
	}

	res := make([]*profilev1.Experiences, len(rows))
	for i, e := range rows {
		res[i] = &profilev1.Experiences{
			Id:          e.ID.String(),
			UserId:      e.UserID.String(),
			Title:       e.Title,
			CompanyName: e.CompanyName,
			StartDate:   timestamppb.New(e.StartDate),
			EndDate:     utils.TimestamppbValue(e.EndDate),
		}
	}
	return res, nil
}
