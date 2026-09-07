package service

import (
	"context"
	"log/slog"
	"uuid"

	"neupaneanish.com.np/profile/internal/errs"
	profilev1 "neupaneanish.com.np/profile/internal/protobuf/common/profile/v1"
	gatewayProfilev1 "neupaneanish.com.np/profile/internal/protobuf/gateway/profile/v1"
	rootProfilev1 "neupaneanish.com.np/profile/internal/protobuf/root/profile/v1"
	"neupaneanish.com.np/profile/internal/repository"
	"neupaneanish.com.np/profile/internal/utils"
)

func (s *GatewayProfileService) UpdateSocial(
	ctx context.Context,
	req *gatewayProfilev1.UpdateSocialRequest,
) (*gatewayProfilev1.UpdateSocialResponse, error) {
	serviceName := "GatewayUpdateSocial"
	userSession := utils.UserSessionContext(ctx)

	if err := updateSocial(
		ctx,
		userSession.UserID,
		userSession.UserID,
		req.GetSocial(),
		serviceName,
		s.cfg.Repository,
		s.cfg.Logger,
	); err != nil {
		return nil, err
	}

	return &gatewayProfilev1.UpdateSocialResponse{Id: req.GetSocial().GetId()}, nil
}

func (s *RootProfileService) UpdateSocial(
	ctx context.Context,
	req *rootProfilev1.UpdateSocialRequest,
) (*rootProfilev1.UpdateSocialResponse, error) {
	serviceName := "RootUpdateSocial"
	userSession := utils.UserSessionContext(ctx)

	userID, userIDErr := parseUUID(ctx, req.GetUserId(), serviceName, s.cfg.Logger)
	if userIDErr != nil {
		return nil, userIDErr
	}

	if err := updateSocial(
		ctx,
		userID,
		userSession.UserID,
		req.GetSocial(),
		serviceName,
		s.cfg.Repository,
		s.cfg.Logger,
	); err != nil {
		return nil, err
	}

	return &rootProfilev1.UpdateSocialResponse{Id: req.GetSocial().GetId()}, nil
}

func updateSocial(
	ctx context.Context,
	userID uuid.UUID,
	updatedBy uuid.UUID,
	req *profilev1.UpdateSocial,
	serviceName string,
	repo repository.Querier,
	logger *slog.Logger,
) error {
	idx, idxErr := parseUUID(ctx, req.GetId(), serviceName, logger)
	if idxErr != nil {
		return idxErr
	}
	params := &repository.UpdateSocialParams{
		Username:  req.GetUsername(),
		UpdatedBy: updatedBy,
		ID:        idx,
		UserID:    userID,
		UpdatedAt: req.GetUpdatedAt().AsTime(),
	}

	cmdTag, err := repo.UpdateSocial(ctx, params)
	if uErr := socialError(ctx, err, serviceName, "update", logger); uErr != nil {
		return uErr
	}

	if cmdTag.RowsAffected() == 0 {
		logger.WarnContext(
			ctx,
			"Social record not found or concurrent modification",
			"service", serviceName,
			"id", idx.String(),
			"userID", userID.String(),
		)
		return errs.ErrConflict
	}
	return nil
}
