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

func (s *GatewayProfileService) Educations(
	ctx context.Context,
	_ *gatewayProfilev1.EducationsRequest,
) (*gatewayProfilev1.EducationsResponse, error) {
	serviceName := "GatewayEducations"
	userSession := utils.UserSessionContext(ctx)

	res, err := internalEducations(ctx, userSession.UserID, s.cfg.Repository, s.cfg.Logger, serviceName)
	if err != nil {
		return nil, err
	}

	return &gatewayProfilev1.EducationsResponse{Educations: res}, nil
}

func (s *RootProfileService) Educations(
	ctx context.Context,
	req *rootProfilev1.EducationsRequest,
) (*rootProfilev1.EducationsResponse, error) {
	serviceName := "RootEducations"

	userID, userIDErr := utils.ParseUUID(ctx, req.GetUserId(), serviceName, s.cfg.Logger)
	if userIDErr != nil {
		return nil, userIDErr
	}

	res, err := internalEducations(ctx, userID, s.cfg.Repository, s.cfg.Logger, serviceName)
	if err != nil {
		return nil, err
	}

	return &rootProfilev1.EducationsResponse{Educations: res}, nil
}

func (s *ExternalProfileService) Educations(
	ctx context.Context,
	_ *externalProfilev1.EducationsRequest,
) (*externalProfilev1.EducationsResponse, error) {
	serviceName := "ExternalEducations"
	userID := utils.DomainUserSessionContext(ctx)

	rows, err := educations(ctx, userID, s.cfg.Repository, s.cfg.Logger, serviceName)
	if err != nil {
		return nil, err
	}

	res := make([]*externalProfilev1.Education, len(rows))
	for i, e := range rows {
		res[i] = &externalProfilev1.Education{
			Id:        e.ID.String(),
			UserId:    e.UserID.String(),
			School:    e.School,
			Degree:    e.Degree,
			StartDate: timestamppb.New(e.StartDate),
			EndDate:   utils.TimestamppbValue(e.EndDate),
		}
	}

	return &externalProfilev1.EducationsResponse{Educations: res}, nil
}

func educations(
	ctx context.Context,
	userID uuid.UUID,
	repo repository.Querier,
	logger *slog.Logger,
	serviceName string,
) ([]*repository.EducationsRow, error) {
	params := &repository.EducationsParams{UserID: userID}

	rows, err := repo.Educations(ctx, params)
	if err != nil {
		logger.ErrorContext(ctx, "Failed to fetch educations", "service", serviceName, "error", err)
		return nil, errs.ErrInternalServer
	}
	return rows, nil
}

func internalEducations(
	ctx context.Context,
	userID uuid.UUID,
	querier repository.Querier,
	logger *slog.Logger,
	serviceName string,
) ([]*profilev1.Education, error) {
	rows, err := educations(ctx, userID, querier, logger, serviceName)
	if err != nil {
		return nil, err
	}

	res := make([]*profilev1.Education, len(rows))
	for i, e := range rows {
		res[i] = &profilev1.Education{
			Id:        e.ID.String(),
			UserId:    e.UserID.String(),
			School:    e.School,
			Degree:    e.Degree,
			StartDate: timestamppb.New(e.StartDate),
			EndDate:   utils.TimestamppbValue(e.EndDate),
		}
	}
	return res, nil
}
