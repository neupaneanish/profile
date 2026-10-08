package service

import (
	"context"

	"neupaneanish.com.np/profile/internal/errs"
	externalProfilev1 "neupaneanish.com.np/profile/internal/protobuf/external/profile/v1"
	"neupaneanish.com.np/profile/internal/repository"
	"neupaneanish.com.np/profile/internal/utils"
)

func (s *ExternalProfileService) Exists(
	ctx context.Context,
	_ *externalProfilev1.ExistsRequest,
) (*externalProfilev1.ExistsResponse, error) {
	serviceName := "ExternalExists"
	userID := utils.DomainUserSessionContext(ctx)

	params := &repository.ExistsParams{UserID: userID}
	row, err := s.cfg.Repository.Exists(ctx, params)
	if err != nil {
		s.cfg.Logger.ErrorContext(ctx, "Failed to query Exists", "service", serviceName, "error", err)
		return nil, errs.ErrInternalServer
	}

	res := &externalProfilev1.Exists{
		Profile:     row.Profile,
		About:       row.About,
		Educations:  row.Educations,
		Experiences: row.Experiences,
		Socials:     row.Socials,
	}

	return &externalProfilev1.ExistsResponse{Exists: res}, nil
}
