package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
	"neupaneanish.com.np/profile/internal/errs"
	rootProfilev1 "neupaneanish.com.np/profile/internal/protobuf/root/profile/v1"
	"neupaneanish.com.np/profile/internal/repository"
	"neupaneanish.com.np/profile/internal/utils"
)

func (s *RootProfileService) UpdatePlatform(
	ctx context.Context,
	req *rootProfilev1.UpdatePlatformRequest,
) (*rootProfilev1.UpdatePlatformResponse, error) {
	serviceName := "UpdatePlatform"
	userSession := utils.UserSessionContext(ctx)

	url, logoURL, urlErr := s.platformURLs(ctx, req.GetPlatform().GetUrl(), req.GetPlatform().GetLogoUrl(), serviceName)
	if urlErr != nil {
		return nil, urlErr
	}

	id, idErr := parseUUID(ctx, req.GetId(), serviceName, s.cfg.Logger)
	if idErr != nil {
		return nil, idErr
	}

	params := &repository.UpdatePlatformParams{
		Name:          req.GetPlatform().GetName(),
		Url:           url,
		UrlSuffix:     req.GetPlatform().GetUrlSuffix(),
		LogoUrl:       logoURL,
		LogoUrlSuffix: req.GetPlatform().GetLogoUrlSuffix(),
		LogoUrlPath:   req.GetPlatform().GetLogoUrlPath(),
		Color:         req.GetPlatform().GetColor(),
		UpdatedBy:     userSession.UserID,
		ID:            id,
		UpdatedAt:     req.GetUpdatedAt().AsTime(),
	}

	cmdTag, err := s.cfg.Repository.UpdatePlatform(ctx, params)
	if pErr := s.platformError(
		ctx,
		err,
		serviceName,
		req.GetPlatform().GetName(),
		url,
		logoURL,
		"update",
	); pErr != nil {
		return nil, pErr
	}

	if cmdTag.RowsAffected() == 0 {
		s.cfg.Logger.WarnContext(
			ctx,
			"Platform record not found or concurrent modification",
			"service", serviceName,
			"id", id.String(),
			"userID", userSession.UserID,
		)
		return nil, errs.ErrConflict
	}
	return &rootProfilev1.UpdatePlatformResponse{Id: id.String()}, nil
}

func (s *RootProfileService) platformError(
	ctx context.Context,
	err error,
	serviceName, name, url, logoURL, method string,
) error {
	if err != nil {
		if pgxErr, ok := errors.AsType[*pgconn.PgError](err); ok && pgxErr.Code == pgerrcode.UniqueViolation {
			switch pgxErr.ConstraintName {
			case utils.PlatformsNameKey:
				s.cfg.Logger.WarnContext(
					ctx,
					"Name already exists",
					"service",
					serviceName,
					"name",
					name,
				)
				return errs.ErrPlatformNameExists
			case utils.PlatformURLKey:
				s.cfg.Logger.WarnContext(
					ctx,
					"URL already exists",
					"service",
					serviceName,
					"url", url,
				)
				return errs.ErrPlatformURLExists
			default:
				s.cfg.Logger.WarnContext(
					ctx,
					"Logo URL already exists",
					"service",
					serviceName,
					"logoURL", logoURL,
				)
				return errs.ErrPlatformLogoURLExists
			}
		}
		s.cfg.Logger.ErrorContext(
			ctx,
			fmt.Sprintf("Failed to %s platform", method),
			"service",
			serviceName,
			"error",
			err,
		)
		return errs.ErrInternalServer
	}
	return nil
}
