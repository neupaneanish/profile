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

func (s *GatewayProfileService) Experience(
	ctx context.Context,
	req *gatewayProfilev1.ExperienceRequest,
) (*gatewayProfilev1.ExperienceResponse, error) {
	serviceName := "GatewayExperience"
	userSession := utils.UserSessionContext(ctx)

	res, err := experience(ctx, userSession.UserID, s.cfg.Repository, s.cfg.Logger, req.GetId(), serviceName)
	if err != nil {
		return nil, err
	}
	return &gatewayProfilev1.ExperienceResponse{
		Experience: res,
	}, nil
}

func (s *RootProfileService) Experience(
	ctx context.Context,
	req *rootProfilev1.ExperienceRequest,
) (*rootProfilev1.ExperienceResponse, error) {
	serviceName := "RootExperience"

	userID, userIDErr := parseUUID(ctx, req.GetUserId(), serviceName, s.cfg.Logger)
	if userIDErr != nil {
		return nil, userIDErr
	}

	res, err := experience(ctx, userID, s.cfg.Repository, s.cfg.Logger, req.GetId(), serviceName)
	if err != nil {
		return nil, err
	}

	return &rootProfilev1.ExperienceResponse{
		Experience: res,
	}, nil
}

func experience(
	ctx context.Context,
	userID uuid.UUID,
	repo repository.Querier,
	logger *slog.Logger,
	id string, serviceName string,
) (*profilev1.Experience, error) {
	idx, idxErr := parseUUID(ctx, id, serviceName, logger)
	if idxErr != nil {
		return nil, idxErr
	}

	params := &repository.ExperienceParams{ID: idx, UserID: userID}

	row, err := repo.Experience(ctx, params)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			logger.WarnContext(ctx, "Experience not found", "service", serviceName)
			return nil, errs.ErrNotFound
		}
		logger.ErrorContext(ctx, "Experience query failed", "service", serviceName, "error", err)
		return nil, errs.ErrInternalServer
	}

	return &profilev1.Experience{
		Id:           row.ID.String(),
		UserId:       row.UserID.String(),
		Title:        row.Title,
		CompanyName:  row.CompanyName,
		Location:     row.Location,
		LocationType: string(row.LocationType),
		StartDate:    timestamppb.New(row.StartDate),
		EndDate:      utils.TimestamppbValue(row.EndDate),
		Description:  utils.StringpbValue(row.Description),
		CreatedAt:    timestamppb.New(row.CreatedAt),
		CreatedBy:    row.CreatedBy.String(),
		UpdatedAt:    timestamppb.New(row.UpdatedAt),
		UpdatedBy:    row.UpdatedBy.String(),
	}, nil
}
