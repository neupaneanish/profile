package service

import (
	"context"

	"neupaneanish.com.np/profile/internal/errs"
	profilev1 "neupaneanish.com.np/profile/internal/protobuf/gateway/profile/v1"
	"neupaneanish.com.np/profile/internal/repository"
	"neupaneanish.com.np/profile/internal/utils"
)

func (s *GatewayProfileService) DeleteSocial(
	ctx context.Context,
	req *profilev1.DeleteSocialRequest,
) (*profilev1.DeleteSocialResponse, error) {
	serviceName := "DeleteSocial"
	userSession := utils.UserSessionContext(ctx)

	id, idErr := parseUUID(ctx, req.GetId(), serviceName, s.cfg.Logger)
	if idErr != nil {
		return nil, idErr
	}

	params := &repository.DeleteSocialParams{
		ID:        id,
		UserID:    userSession.UserID,
		UpdatedAt: req.GetUpdatedAt().AsTime(),
	}

	cmdTag, err := s.cfg.Repository.DeleteSocial(ctx, params)
	if err != nil {
		s.cfg.Logger.ErrorContext(ctx, "Delete social Failed", "service", serviceName, "error", err)
		return nil, errs.ErrInternalServer
	}

	if cmdTag.RowsAffected() == 0 {
		s.cfg.Logger.WarnContext(
			ctx,
			"Social record not found or concurrent modification",
			"service", serviceName,
			"id", req.GetId(),
		)
		return nil, errs.ErrConflict
	}

	return &profilev1.DeleteSocialResponse{}, nil
}
