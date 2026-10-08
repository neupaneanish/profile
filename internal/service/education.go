package service

import (
	"context"
	"errors"
	"log/slog"
	"uuid"

	"github.com/jackc/pgx/v5"
	"google.golang.org/protobuf/types/known/timestamppb"
	"neupaneanish.com.np/profile/internal/errs"
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

	row, err := education(
		ctx,
		req.GetId(),
		serviceName,
		userSession.UserID,
		s.cfg.Repository,
		s.cfg.Logger,
	)
	if err != nil {
		return nil, err
	}

	return &gatewayProfilev1.EducationResponse{
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
		Description:   utils.StringpbValue(row.Description),
		UpdatedAt:     timestamppb.New(row.UpdatedAt),
	}, nil
}

func (s *RootProfileService) Education(
	ctx context.Context,
	req *rootProfilev1.EducationRequest,
) (*rootProfilev1.EducationResponse, error) {
	serviceName := "RootEducation"
	userSession := utils.UserSessionContext(ctx)

	userID, userIDErr := utils.ParseUUID(ctx, req.GetUserId(), serviceName, s.cfg.Logger)
	if userIDErr != nil {
		return nil, userIDErr
	}

	row, err := education(
		ctx,
		req.GetId(),
		serviceName,
		userID,
		s.cfg.Repository,
		s.cfg.Logger,
	)
	if err != nil {
		return nil, err
	}

	createdByUsername, updatedByUsername, usernamesErr := utils.GetUsernames(
		ctx,
		row.CreatedBy,
		row.UpdatedBy,
		userSession,
		s.cfg.Client,
		s.cfg.Logger,
	)
	if usernamesErr != nil {
		return nil, usernamesErr
	}

	return &rootProfilev1.EducationResponse{
		Id:                row.ID.String(),
		UserId:            row.UserID.String(),
		School:            row.School,
		Degree:            row.Degree,
		Affiliation:       utils.StringpbValue(row.Affiliation),
		FieldOfStudy:      utils.StringpbValue(row.FieldOfStudy),
		Concentration:     utils.StringpbValue(row.Concentration),
		StartDate:         timestamppb.New(row.StartDate),
		EndDate:           utils.TimestamppbValue(row.EndDate),
		Address:           row.Address,
		Description:       utils.StringpbValue(row.Description),
		CreatedAt:         timestamppb.New(row.CreatedAt),
		CreatedBy:         row.CreatedBy.String(),
		UpdatedAt:         timestamppb.New(row.UpdatedAt),
		UpdatedBy:         row.UpdatedBy.String(),
		CreatedByUsername: createdByUsername,
		UpdatedByUsername: updatedByUsername,
	}, nil
}

func education(
	ctx context.Context,
	id string, serviceName string,
	userID uuid.UUID,
	repo repository.Querier,
	logger *slog.Logger,
) (*repository.Education, error) {
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

	return row, nil
}
