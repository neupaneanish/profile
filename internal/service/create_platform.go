package service

import (
	"context"

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
