package service

import (
	"context"
	"log/slog"
	"time"
	"uuid"

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

	id, err := createUpdateEducation(
		ctx,
		req.GetId(),
		userSession.UserID,
		userSession.UserID,
		req.GetEducation(),
		req.GetUpdatedAt().AsTime(),
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

	idx, err := createUpdateEducation(
		ctx,
		req.GetId(),
		userID,
		userSession.UserID,
		req.GetEducation(),
		req.GetUpdatedAt().AsTime(),
		s.cfg.Repository,
		s.cfg.Logger,
		serviceName,
	)
	if err != nil {
		return nil, err
	}
	return &rootProfilev1.UpdateEducationResponse{Id: idx.String()}, nil
}

func createUpdateEducation(
	ctx context.Context,
	id string,
	userID, updatedBy uuid.UUID,
	req *profilev1.CreateUpdateEducation,
	updatedAt time.Time,
	repo repository.Querier,
	logger *slog.Logger,
	serviceName string,
) (uuid.UUID, error) {
	if id == "" {
		params := &repository.CreateEducationParams{
			UserID:        userID,
			School:        req.GetSchool(),
			Degree:        req.GetDegree(),
			Affiliation:   utils.StringValue(req.GetAffiliation()),
			FieldOfStudy:  utils.StringValue(req.GetFieldOfStudy()),
			Concentration: utils.StringValue(req.GetConcentration()),
			StartDate:     req.GetStartDate().AsTime(),
			EndDate:       utils.TimestampValue(req.GetEndDate()),
			Address:       req.GetAddress(),
			Description:   utils.StringValue(req.GetDescription()),
			CreatedBy:     userID,
			UpdatedBy:     userID,
		}
		idx, err := repo.CreateEducation(ctx, params)
		if err != nil {
			logger.ErrorContext(ctx, "Create Education Failed", "service", serviceName, "error", err)
			return uuid.Nil(), errs.ErrInternalServer
		}
		return idx, nil
	}

	idx, updateIDErr := parseUUID(ctx, id, serviceName, logger)
	if updateIDErr != nil {
		return uuid.Nil(), updateIDErr
	}

	params := &repository.UpdateEducationParams{
		School:        req.GetSchool(),
		Degree:        req.GetDegree(),
		Affiliation:   utils.StringValue(req.GetAffiliation()),
		FieldOfStudy:  utils.StringValue(req.GetFieldOfStudy()),
		Concentration: utils.StringValue(req.GetConcentration()),
		StartDate:     req.GetStartDate().AsTime(),
		EndDate:       utils.TimestampValue(req.GetEndDate()),
		Address:       req.GetAddress(),
		Description:   utils.StringValue(req.GetDescription()),
		UpdatedBy:     updatedBy,
		ID:            idx,
		UserID:        userID,
		UpdatedAt:     updatedAt,
	}

	cmdTag, err := repo.UpdateEducation(ctx, params)
	if err != nil {
		logger.ErrorContext(ctx, "Update Education Failed", "service", serviceName, "error", err)
		return uuid.Nil(), errs.ErrInternalServer
	}

	if cmdTag.RowsAffected() == 0 {
		logger.WarnContext(
			ctx,
			"Concurrent education update detected",
			"service",
			serviceName,
			"userID", userID,
			"id", id,
		)
		return uuid.Nil(), errs.ErrConflict
	}
	return idx, nil
}
