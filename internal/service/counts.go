package service

import (
	"context"

	"neupaneanish.com.np/profile/internal/errs"
	gatewayProfilev1 "neupaneanish.com.np/profile/internal/protobuf/gateway/profile/v1"
	rootProfilev1 "neupaneanish.com.np/profile/internal/protobuf/root/profile/v1"
	"neupaneanish.com.np/profile/internal/repository"
	"neupaneanish.com.np/profile/internal/utils"
)

func (s *GatewayProfileService) Counts(
	ctx context.Context,
	_ *gatewayProfilev1.CountsRequest,
) (*gatewayProfilev1.CountsResponse, error) {
	serviceName := "Counts"
	userSession := utils.UserSessionContext(ctx)

	params := &repository.ProfileCountsParams{UserID: userSession.UserID}

	row, err := s.cfg.Repository.ProfileCounts(ctx, params)
	if err != nil {
		s.cfg.Logger.ErrorContext(ctx, "Failed to query counts", "service", serviceName, "error", err)
		return nil, errs.ErrInternalServer
	}

	return &gatewayProfilev1.CountsResponse{
		Domains:     row.Domains,
		Educations:  row.Educations,
		Experiences: row.Experiences,
		Socials:     row.Socials,
	}, nil
}

func (s *RootProfileService) Counts(
	ctx context.Context,
	_ *rootProfilev1.CountsRequest,
) (*rootProfilev1.CountsResponse, error) {
	serviceName := "Counts"

	row, err := s.cfg.Repository.RootCounts(ctx)
	if err != nil {
		s.cfg.Logger.ErrorContext(ctx, "Failed to query counts", "service", serviceName, "error", err)
		return nil, errs.ErrInternalServer
	}

	return &rootProfilev1.CountsResponse{
		Profiles:    row.Profiles,
		Abouts:      row.Abouts,
		Domains:     row.Domains,
		Educations:  row.Educations,
		Experiences: row.Experiences,
		Socials:     row.Socials,
		Nameservers: row.Nameservers,
		Icons:       row.Icons,
		Templates:   row.Templates,
	}, nil
}
