//nolint:dupl // Clean handler pattern intentionally mirrors Educations
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

func (s *GatewayProfileService) Educations(
	ctx context.Context,
	_ *gatewayProfilev1.EducationsRequest,
) (*gatewayProfilev1.EducationsResponse, error) {
	serviceName := "GatewayEducations"
	userSession := utils.UserSessionContext(ctx)

	res, err := educations(ctx, userSession.UserID, s.cfg.Repository, s.cfg.Logger, serviceName)
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

	userID, userIDErr := parseUUID(ctx, req.GetUserId(), serviceName, s.cfg.Logger)
	if userIDErr != nil {
		return nil, userIDErr
	}

	res, err := educations(ctx, userID, s.cfg.Repository, s.cfg.Logger, serviceName)
	if err != nil {
		return nil, err
	}

	return &rootProfilev1.EducationsResponse{Educations: res}, nil
}

func educations(
	ctx context.Context,
	userID uuid.UUID,
	repo repository.Querier,
	logger *slog.Logger,
	serviceName string,
) ([]*profilev1.Educations, error) {
	params := &repository.EducationsParams{UserID: userID}

	rows, err := repo.Educations(ctx, params)
	if err != nil {
		logger.ErrorContext(ctx, "Failed to fetch educations", "service", serviceName, "error", err)
		return nil, errs.ErrInternalServer
	}

	res := make([]*profilev1.Educations, len(rows))
	for i, e := range rows {
		res[i] = &profilev1.Educations{
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
