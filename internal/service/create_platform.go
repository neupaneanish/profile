package service

import (
	"context"

	"neupaneanish.com.np/profile/internal/errs"
	rootProfilev1 "neupaneanish.com.np/profile/internal/protobuf/root/profile/v1"
	"neupaneanish.com.np/profile/internal/repository"
	"neupaneanish.com.np/profile/internal/utils"
)

func (s *RootProfileService) CreatePlatform(
	ctx context.Context,
	req *rootProfilev1.CreatePlatformRequest,
) (*rootProfilev1.CreatePlatformResponse, error) {
	serviceName := "CreatePlatform"
	userSession := utils.UserSessionContext(ctx)

	url, logoURL, urlErr := s.platformURLs(ctx, req.GetPlatform().GetUrl(), req.GetPlatform().GetLogoUrl(), serviceName)
	if urlErr != nil {
		return nil, urlErr
	}

	params := &repository.CreatePlatformParams{
		Name:          req.GetPlatform().GetName(),
		Url:           url,
		UrlSuffix:     req.GetPlatform().GetUrlSuffix(),
		LogoUrl:       logoURL,
		LogoUrlSuffix: req.GetPlatform().GetLogoUrlSuffix(),
		LogoUrlPath:   req.GetPlatform().GetLogoUrlPath(),
		Color:         req.GetPlatform().GetColor(),
		CreatedBy:     userSession.UserID,
		UpdatedBy:     userSession.UserID,
	}

	id, err := s.cfg.Repository.CreatePlatform(ctx, params)
	if pErr := s.platformError(
		ctx,
		err,
		serviceName,
		req.GetPlatform().GetName(),
		url,
		logoURL,
		"create",
	); pErr != nil {
		return nil, pErr
	}
	return &rootProfilev1.CreatePlatformResponse{Id: id.String()}, nil
}

func (s *RootProfileService) platformURLs(
	ctx context.Context,
	url, logoURL, serviceName string,
) (string, string, error) {
	pURL, urlErr := utils.ValidateURL(url)
	if urlErr != nil {
		s.cfg.Logger.WarnContext(
			ctx,
			"Invalid URL",
			"service",
			serviceName,
			"error",
			urlErr,
			"url",
			url,
		)
		return "", "", errs.ErrInvalidURL
	}

	pLogoURL, logoURLErr := utils.ValidateURL(logoURL)
	if logoURLErr != nil {
		s.cfg.Logger.ErrorContext(
			ctx,
			"Invalid Logo URL",
			"service",
			serviceName,
			"error",
			logoURLErr,
			"logoURL",
			logoURL,
		)
		return "", "", errs.ErrInvalidURL
	}
	return "https://" + pURL, "https://" + pLogoURL, nil
}
