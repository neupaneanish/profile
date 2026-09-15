package service

import (
	"context"
	"errors"
	"log/slog"
	"uuid"

	"github.com/jackc/pgx/v5"
	"google.golang.org/protobuf/types/known/timestamppb"
	"neupaneanish.com.np/profile/internal/errs"
	profilev1 "neupaneanish.com.np/profile/internal/protobuf/common/profile/v1"
	gatewayProfilev1 "neupaneanish.com.np/profile/internal/protobuf/gateway/profile/v1"
	rootProfilev1 "neupaneanish.com.np/profile/internal/protobuf/root/profile/v1"
	"neupaneanish.com.np/profile/internal/repository"
	"neupaneanish.com.np/profile/internal/utils"
)

func (s *GatewayProfileService) Education(
	ctx context.Context,
	req *gatewayProfilev1.EducationRequest,
) (*gatewayProfilev1.EducationResponse, error) {
	serviceName := "GatewayEducation"
	userSession := utils.UserSessionContext(ctx)

	res, err := education(ctx, userSession.UserID, s.cfg.Repository, s.cfg.Logger, req.GetId(), serviceName)
	if err != nil {
		return nil, err
	}
	return &gatewayProfilev1.EducationResponse{
		Education: res,
	}, nil
}

func (s *RootProfileService) Education(
	ctx context.Context,
	req *rootProfilev1.EducationRequest,
) (*rootProfilev1.EducationResponse, error) {
	serviceName := "RootEducation"

	userID, userIDErr := utils.ParseUUID(ctx, req.GetUserId(), serviceName, s.cfg.Logger)
	if userIDErr != nil {
		return nil, userIDErr
	}

	res, err := education(ctx, userID, s.cfg.Repository, s.cfg.Logger, req.GetId(), serviceName)
	if err != nil {
		return nil, err
	}

	return &rootProfilev1.EducationResponse{
		Education: res,
	}, nil
}

func education(
	ctx context.Context,
	userID uuid.UUID,
	repo repository.Querier,
	logger *slog.Logger,
	id string, serviceName string,
) (*profilev1.Education, error) {
	idx, idxErr := utils.ParseUUID(ctx, id, serviceName, logger)
	if idxErr != nil {
		return nil, idxErr
	}

	params := &repository.EducationParams{ID: idx, UserID: userID}

	row, err := repo.Education(ctx, params)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			logger.WarnContext(ctx, "Education not found", "service", serviceName)
			return nil, errs.ErrNotFound("Education")
		}
		logger.ErrorContext(ctx, "Education query failed", "service", serviceName, "error", err)
		return nil, errs.ErrInternalServer
	}

	return &profilev1.Education{
		Id:            row.ID.String(),
		UserId:        row.UserID.String(),
		School:        row.School,
		Degree:        row.Degree,
		Affiliation:   utils.StringpbValue(row.Affiliation),
		FieldOfStudy:  utils.StringpbValue(row.FieldOfStudy),
		Concentration: utils.StringpbValue(row.Concentration),
		StartDate:     timestamppb.New(row.StartDate),
		EndDate:       utils.TimestamppbValue(row.EndDate),
		Address:       row.Address,
		Description:   utils.StringpbValue(row.Concentration),
		CreatedAt:     timestamppb.New(row.CreatedAt),
		CreatedBy:     row.CreatedBy.String(),
		UpdatedAt:     timestamppb.New(row.UpdatedAt),
		UpdatedBy:     row.UpdatedBy.String(),
	}, nil
}
