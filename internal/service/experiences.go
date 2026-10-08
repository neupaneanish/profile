package service

import (
	"context"
	"log/slog"
	"uuid"

	"google.golang.org/protobuf/types/known/timestamppb"
	"neupaneanish.com.np/profile/internal/errs"
	profilev1 "neupaneanish.com.np/profile/internal/protobuf/common/profile/v1"
	externalProfilev1 "neupaneanish.com.np/profile/internal/protobuf/external/profile/v1"
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

	res, err := internalExperiences(ctx, userSession.UserID, s.cfg.Repository, s.cfg.Logger, serviceName)
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

	userID, userIDErr := utils.ParseUUID(ctx, req.GetUserId(), serviceName, s.cfg.Logger)
	if userIDErr != nil {
		return nil, userIDErr
	}

	res, err := internalExperiences(ctx, userID, s.cfg.Repository, s.cfg.Logger, serviceName)
	if err != nil {
		return nil, err
	}

	return &rootProfilev1.ExperiencesResponse{Experiences: res}, nil
}

func (s *ExternalProfileService) Experiences(
	ctx context.Context,
	_ *externalProfilev1.ExperiencesRequest,
) (*externalProfilev1.ExperiencesResponse, error) {
	serviceName := "ExternalExperiences"
	userID := utils.DomainUserSessionContext(ctx)

	rows, err := experiences(ctx, userID, s.cfg.Repository, s.cfg.Logger, serviceName)
	if err != nil {
		return nil, err
	}

	res := make([]*externalProfilev1.Experience, len(rows))
	for i, e := range rows {
		res[i] = &externalProfilev1.Experience{
			Id:          e.ID.String(),
			UserId:      e.UserID.String(),
			Title:       e.Title,
			CompanyName: e.CompanyName,
			StartDate:   timestamppb.New(e.StartDate),
			EndDate:     utils.TimestamppbValue(e.EndDate),
		}
	}

	return &externalProfilev1.ExperiencesResponse{Experiences: res}, nil
}

func experiences(
	ctx context.Context,
	userID uuid.UUID,
	repo repository.Querier,
	logger *slog.Logger,
	serviceName string,
) ([]*repository.ExperiencesRow, error) {
	params := &repository.ExperiencesParams{UserID: userID}

	rows, err := repo.Experiences(ctx, params)
	if err != nil {
		logger.ErrorContext(ctx, "Failed to fetch experiences", "service", serviceName, "error", err)
		return nil, errs.ErrInternalServer
	}

	return rows, nil
}

func internalExperiences(
	ctx context.Context,
	userID uuid.UUID,
	querier repository.Querier,
	logger *slog.Logger,
	serviceName string,
) ([]*profilev1.Experience, error) {
	rows, err := experiences(ctx, userID, querier, logger, serviceName)
	if err != nil {
		return nil, err
	}

	res := make([]*profilev1.Experience, len(rows))
	for i, e := range rows {
		res[i] = &profilev1.Experience{
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
