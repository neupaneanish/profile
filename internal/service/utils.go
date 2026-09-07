package service

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"time"
	"uuid"

	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"google.golang.org/protobuf/types/known/timestamppb"
	"neupaneanish.com.np/profile/internal/enum"
	"neupaneanish.com.np/profile/internal/errs"
	profilev1 "neupaneanish.com.np/profile/internal/protobuf/common/profile/v1"
	rootProfilev1 "neupaneanish.com.np/profile/internal/protobuf/root/profile/v1"
	"neupaneanish.com.np/profile/internal/repository"
	"neupaneanish.com.np/profile/internal/utils"
)

func parseUUID(ctx context.Context, userIDStr, serviceName string, logger *slog.Logger) (uuid.UUID, error) {
	userID, uuidErr := uuid.Parse(userIDStr)
	if uuidErr != nil {
		logger.WarnContext(ctx, "Failed to parse userID", "service", serviceName, "userID", userIDStr)
		return uuid.Nil(), errs.ErrInvalidUserID
	}
	return userID, nil
}

func (s *RootProfileService) nameServerError(ctx context.Context, err error, serviceName, method string) error {
	if err != nil {
		if pgErr, ok := errors.AsType[*pgconn.PgError](err); ok && pgErr.Code == pgerrcode.UniqueViolation {
			s.cfg.Logger.WarnContext(
				ctx,
				"NameServer Already Exists",
				"service", serviceName,
			)
			return errs.ErrUniqueViolation("Nameserver")
		}
		s.cfg.Logger.ErrorContext(
			ctx,
			fmt.Sprintf("Failed to %s name server", method),
			"service", serviceName,
			"error", err,
		)
		return errs.ErrInternalServer
	}
	return nil
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

func socialError(
	ctx context.Context,
	err error,
	serviceName, method string,
	logger *slog.Logger,
) error {
	if err != nil {
		if pgErr, ok := errors.AsType[*pgconn.PgError](err); ok {
			switch pgErr.Code {
			case pgerrcode.UniqueViolation:
				logger.WarnContext(
					ctx,
					"Already Exists",
					"service",
					serviceName,
				)
				return errs.ErrUniqueViolation("Social")
			case pgerrcode.ForeignKeyViolation:
				logger.ErrorContext(
					ctx,
					"Platform does not exists",
					"service",
					serviceName,
				)
				return errs.ErrForeignKeyViolation("Platform")
			}
		}
		logger.ErrorContext(
			ctx,
			fmt.Sprintf("Failed to %s social", method),
			"service",
			serviceName,
			"error",
			err,
		)
		return errs.ErrInternalServer
	}

	return nil
}

func deleteDB(
	ctx context.Context,
	cmdTag pgconn.CommandTag,
	err error,
	serviceName, table, id string,
	logger *slog.Logger,
) error {
	if err != nil {
		logger.ErrorContext(ctx, fmt.Sprintf("Delete %s Failed", table), "service", serviceName, "error", err)
		return errs.ErrInternalServer
	}

	if cmdTag.RowsAffected() == 0 {
		logger.WarnContext(
			ctx,
			fmt.Sprintf("%s record not found or concurrent modification", table),
			"service", serviceName,
			"id", id,
		)
		return errs.ErrConflict
	}
	return nil
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

func createUpdateExperience(
	ctx context.Context,
	id string,
	userID, updatedBy uuid.UUID,
	req *profilev1.CreateUpdateExperience,
	updatedAt time.Time,
	repo repository.Querier,
	logger *slog.Logger,
	serviceName string,
) (uuid.UUID, error) {
	if id == "" {
		params := &repository.CreateExperienceParams{
			UserID:       userID,
			Title:        req.GetTitle(),
			CompanyName:  req.GetCompanyName(),
			Location:     req.GetLocation(),
			LocationType: enum.LocationType(req.GetLocationType()),
			StartDate:    req.GetStartDate().AsTime(),
			EndDate:      utils.TimestampValue(req.GetEndDate()),
			Description:  utils.StringValue(req.GetDescription()),
			CreatedBy:    userID,
			UpdatedBy:    userID,
		}

		idx, err := repo.CreateExperience(ctx, params)
		if err != nil {
			logger.ErrorContext(ctx, "Create Experience Failed", "service", serviceName, "error", err)
			return uuid.Nil(), errs.ErrInternalServer
		}
		return idx, nil
	}

	idx, idErr := parseUUID(ctx, id, serviceName, logger)
	if idErr != nil {
		return uuid.Nil(), idErr
	}

	params := &repository.UpdateExperienceParams{
		Title:        req.GetTitle(),
		CompanyName:  req.GetCompanyName(),
		Location:     req.GetLocation(),
		LocationType: enum.LocationType(req.GetLocationType()),
		StartDate:    req.GetStartDate().AsTime(),
		EndDate:      utils.TimestampValue(req.GetEndDate()),
		Description:  utils.StringValue(req.GetDescription()),
		UpdatedBy:    updatedBy,
		ID:           idx,
		UserID:       userID,
		UpdatedAt:    updatedAt,
	}

	cmdTag, err := repo.UpdateExperience(ctx, params)
	if err != nil {
		logger.ErrorContext(ctx, "Update Education Failed", "service", serviceName, "error", err)
		return uuid.Nil(), errs.ErrInternalServer
	}

	if cmdTag.RowsAffected() == 0 {
		logger.WarnContext(
			ctx,
			"Concurrent experience update detected",
			"service",
			serviceName,
			"userID", userID,
			"id", id,
		)
		return uuid.Nil(), errs.ErrConflict
	}
	return idx, nil
}

func (s *RootProfileService) createUpdateNameServer(
	ctx context.Context,
	id string,
	req *rootProfilev1.CreateUpdateNameServer,
	updatedBy uuid.UUID,
	serviceName string,
	active bool,
	updatedAt time.Time,
) (uuid.UUID, error) {
	domain, domainErr := utils.ValidateURL(req.GetDomain())
	if domainErr != nil {
		s.cfg.Logger.ErrorContext(ctx, "Invalid Domain", "service", serviceName, "error", domainErr)
		return uuid.Nil(), errs.ErrInvalidURL
	}

	cname := req.GetCname()

	if id == "" {
		params := &repository.CreateNameServerParams{
			Domain:    domain,
			Cname:     cname,
			CreatedBy: updatedBy,
			UpdatedBy: updatedBy,
		}

		idx, err := s.cfg.Repository.CreateNameServer(ctx, params)

		if nsErr := s.nameServerError(ctx, err, serviceName, "create"); nsErr != nil {
			return uuid.Nil(), nsErr
		}

		return idx, nil
	}

	idx, idxErr := parseUUID(ctx, id, serviceName, s.cfg.Logger)
	if idxErr != nil {
		return idx, idxErr
	}

	params := &repository.UpdateNameServerParams{
		Domain:    domain,
		Cname:     cname,
		Active:    active,
		ID:        idx,
		UpdatedAt: updatedAt,
	}

	cmdTag, err := s.cfg.Repository.UpdateNameServer(ctx, params)
	if nsErr := s.nameServerError(ctx, err, serviceName, "update"); nsErr != nil {
		return uuid.Nil(), nsErr
	}

	if cmdTag.RowsAffected() == 0 {
		s.cfg.Logger.WarnContext(
			ctx,
			"Concurrent nameserver update detected",
			"service",
			serviceName,
			"id", id,
		)
		return uuid.Nil(), errs.ErrConflict
	}
	return idx, nil
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
				return errs.ErrUniqueViolation("Platform")
			case utils.PlatformURLKey:
				s.cfg.Logger.WarnContext(
					ctx,
					"URL already exists",
					"service",
					serviceName,
					"url", url,
				)
				return errs.ErrUniqueViolation("URL")
			default:
				s.cfg.Logger.WarnContext(
					ctx,
					"Logo URL already exists",
					"service",
					serviceName,
					"logoURL", logoURL,
				)
				return errs.ErrUniqueViolation("Logo URL")
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

func updateProfile(
	ctx context.Context,
	userID, updatedBy uuid.UUID,
	req *profilev1.CreateUpdateProfile,
	updatedAt time.Time,
	repo repository.Querier,
	logger *slog.Logger,
	serviceName string,
) (*profilev1.Profile, error) {
	params := &repository.UpdateProfileParams{
		Name:      req.GetName(),
		Title:     req.GetTitle(),
		UpdatedBy: updatedBy,
		UserID:    userID,
		UpdatedAt: updatedAt,
	}

	row, rowErr := repo.UpdateProfile(ctx, params)
	if rowErr != nil {
		if errors.Is(rowErr, pgx.ErrNoRows) {
			logger.WarnContext(ctx, "Concurrent profile update detected", "service", serviceName, "userID", userID)
			return nil, errs.ErrConflict
		}
		logger.ErrorContext(ctx, "Update Profile Failed", "service", serviceName, "error", rowErr)
		return nil, errs.ErrInternalServer
	}
	return &profilev1.Profile{
		UserId:    row.UserID.String(),
		Name:      row.Name,
		Title:     row.Title,
		Dob:       timestamppb.New(row.Dob),
		CreatedAt: timestamppb.New(row.CreatedAt),
		CreatedBy: row.CreatedBy.String(),
		UpdatedAt: timestamppb.New(row.UpdatedAt),
		UpdatedBy: row.UpdatedBy.String(),
	}, nil
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
