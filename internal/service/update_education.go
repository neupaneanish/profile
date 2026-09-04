package service

import (
	"context"
	"errors"
	"log/slog"
	"uuid"

	"github.com/jackc/pgx/v5"
	"neupaneanish.com.np/profile/internal/errs"
	profilev1 "neupaneanish.com.np/profile/internal/protobuf/common/profile/v1"
	gatewayProfilev1 "neupaneanish.com.np/profile/internal/protobuf/gateway/profile/v1"
	rootProfilev1 "neupaneanish.com.np/profile/internal/protobuf/root/profile/v1"
	"neupaneanish.com.np/profile/internal/repository"
	"neupaneanish.com.np/profile/internal/utils"
)

func (s *GatewayProfileService) UpdateEducation(
	ctx context.Context,
	req *gatewayProfilev1.UpdateEducationRequest,
) (*gatewayProfilev1.UpdateEducationResponse, error) {
	serviceName := "GatewayUpdateEducation"
	userSession := utils.UserSessionContext(ctx)

	id, err := updateEducation(
		ctx,
		userSession.UserID,
		userSession.UserID,
		req.GetEducation(),
		s.cfg.Repository,
		s.cfg.Logger,
		serviceName,
	)
	if err != nil {
		return nil, err
	}
	return &gatewayProfilev1.UpdateEducationResponse{Id: id.String()}, nil
}

func (s *RootProfileService) UpdateEducation(
	ctx context.Context,
	req *rootProfilev1.UpdateEducationRequest,
) (*rootProfilev1.UpdateEducationResponse, error) {
	serviceName := "RootUpdateEducation"
	userSession := utils.UserSessionContext(ctx)

	userID, userIDErr := parseUUID(ctx, req.GetUserId(), serviceName, s.cfg.Logger)
	if userIDErr != nil {
		return nil, userIDErr
	}

	id, err := updateEducation(
		ctx,
		userID,
		userSession.UserID,
		req.GetEducation(),
		s.cfg.Repository,
		s.cfg.Logger,
		serviceName,
	)
	if err != nil {
		return nil, err
	}
	return &rootProfilev1.UpdateEducationResponse{Id: id.String()}, nil
}

func updateEducation(
	ctx context.Context,
	userID, updatedBy uuid.UUID,
	req *profilev1.UpdateEducation,
	repo repository.Querier,
	logger *slog.Logger,
	serviceName string,
) (uuid.UUID, error) {
	id, idErr := parseUUID(ctx, req.GetId(), serviceName, logger)
	if idErr != nil {
		return id, idErr
	}
	params := &repository.UpdateEducationParams{
		School:        req.GetSchool(),
		Degree:        req.GetDegree(),
		Affiliation:   stringValue(req.GetAffiliation()),
		FieldOfStudy:  stringValue(req.GetFieldOfStudy()),
		Concentration: stringValue(req.GetConcentration()),
		StartDate:     req.GetStartDate().AsTime(),
		EndDate:       timestampValue(req.GetEndDate()),
		Address:       req.GetAddress(),
		Description:   stringValue(req.GetDescription()),
		UpdatedBy:     updatedBy,
		ID:            id,
		UserID:        userID,
		UpdatedAt:     req.GetUpdatedAt().AsTime(),
	}
	eduID, err := repo.UpdateEducation(ctx, params)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			logger.WarnContext(ctx, "Concurrent education update detected", "service", serviceName, "userID", userID, "id", id.String())
			return uuid.Nil(), errs.ErrConflict
		}
		logger.ErrorContext(ctx, "Update Education Failed", "service", serviceName, "error", err)
		return uuid.Nil(), errs.ErrInternalServer
	}
	return eduID, nil
}
