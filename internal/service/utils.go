package service

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"slices"
	"strings"
	"time"
	"uuid"

	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/twmb/franz-go/pkg/kgo"
	"github.com/valkey-io/valkey-go"
	"neupaneanish.com.np/profile/internal/enum"
	"neupaneanish.com.np/profile/internal/errs"
	profilev1 "neupaneanish.com.np/profile/internal/protobuf/common/profile/v1"
	rootProfilev1 "neupaneanish.com.np/profile/internal/protobuf/root/profile/v1"
	"neupaneanish.com.np/profile/internal/redpanda"
	"neupaneanish.com.np/profile/internal/repository"
	"neupaneanish.com.np/profile/internal/utils"
)

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
					"Icon does not exists",
					"service",
					serviceName,
				)
				return errs.ErrForeignKeyViolation("Icon")
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
	affected int64,
	err error,
	serviceName, id, table string,
	session *utils.UserSession,
	userID uuid.UUID,
	client valkey.Client,
	rpClient *kgo.Client,
	logger *slog.Logger,
) error {
	if err != nil {
		logger.ErrorContext(ctx, "delete database record failed", "service", serviceName, "table", table, "error", err)
		return errs.ErrInternalServer
	}

	if affected != 1 {
		logger.WarnContext(
			ctx,
			"record not found or concurrent modification",
			"service", serviceName,
			"id", id,
			"table", table,
		)
		return errs.ErrConflict
	}
	redpanda.RootNotificationProduce(
		ctx,
		session,
		userID,
		table,
		utils.DatabaseMethodDelete,
		serviceName,
		client,
		rpClient,
		logger,
	)

	return nil
}

//nolint:funlen
func createUpdateEducation(
	ctx context.Context,
	id, serviceName string,
	session *utils.UserSession,
	userID uuid.UUID,
	req *profilev1.CreateUpdateEducation,
	updatedAt time.Time,
	repo repository.Querier,
	client valkey.Client,
	rpClient *kgo.Client,
	logger *slog.Logger,
) error {
	school := req.GetSchool()
	degree := req.GetDegree()
	affiliation := utils.StringValue(req.GetAffiliation())
	fieldOfStudy := utils.StringValue(req.GetFieldOfStudy())
	concentration := utils.StringValue(req.GetConcentration())
	startDate := req.GetStartDate().AsTime()
	endDate := utils.TimestampValue(req.GetEndDate())
	address := req.GetAddress()
	description := utils.StringValue(req.GetDescription())

	if id == "" {
		params := &repository.CreateEducationParams{
			UserID:        userID,
			School:        school,
			Degree:        degree,
			Affiliation:   affiliation,
			FieldOfStudy:  fieldOfStudy,
			Concentration: concentration,
			StartDate:     startDate,
			EndDate:       endDate,
			Address:       address,
			Description:   description,
			CreatedBy:     session.UserID,
			UpdatedBy:     session.UserID,
		}
		if _, err := repo.CreateEducation(ctx, params); err != nil {
			logger.ErrorContext(ctx, "Create Education Failed", "service", serviceName, "error", err)
			return errs.ErrInternalServer
		}
		redpanda.RootNotificationProduce(
			ctx,
			session,
			session.UserID,
			utils.DatabaseTableEducation,
			utils.DatabaseMethodCreate,
			serviceName,
			client,
			rpClient,
			logger,
		)
		return nil
	}

	idx, updateIDErr := utils.ParseUUID(ctx, id, serviceName, logger)
	if updateIDErr != nil {
		return updateIDErr
	}

	params := &repository.UpdateEducationParams{
		School:        school,
		Degree:        degree,
		Affiliation:   affiliation,
		FieldOfStudy:  fieldOfStudy,
		Concentration: concentration,
		StartDate:     startDate,
		EndDate:       endDate,
		Address:       address,
		Description:   description,
		UpdatedBy:     session.UserID,
		ID:            idx,
		UserID:        userID,
		UpdatedAt:     updatedAt,
	}

	affected, err := repo.UpdateEducation(ctx, params)
	if err != nil {
		logger.ErrorContext(ctx, "Update Education Failed", "service", serviceName, "error", err)
		return errs.ErrInternalServer
	}

	if affected != 1 {
		logger.WarnContext(
			ctx,
			"Concurrent education update detected",
			"service",
			serviceName,
			"userID", userID,
			"id", id,
		)
		return errs.ErrConflict
	}
	redpanda.RootNotificationProduce(
		ctx,
		session,
		session.UserID,
		utils.DatabaseTableEducation,
		utils.DatabaseMethodUpdate,
		serviceName,
		client,
		rpClient,
		logger,
	)

	return nil
}

//nolint:funlen
func createUpdateExperience(
	ctx context.Context,
	id, serviceName string,
	session *utils.UserSession,
	userID uuid.UUID,
	req *profilev1.CreateUpdateExperience,
	updatedAt time.Time,
	repo repository.Querier,
	client valkey.Client,
	rpClient *kgo.Client,
	logger *slog.Logger,
) error {
	title := req.GetTitle()
	companyName := req.GetCompanyName()
	location := req.GetLocation()
	locationType := enum.LocationType(req.GetLocationType())
	startDate := req.GetStartDate().AsTime()
	endDate := utils.TimestampValue(req.GetEndDate())
	description := utils.StringValue(req.GetDescription())

	if id == "" {
		params := &repository.CreateExperienceParams{
			UserID:       userID,
			Title:        title,
			CompanyName:  companyName,
			Location:     location,
			LocationType: locationType,
			StartDate:    startDate,
			EndDate:      endDate,
			Description:  description,
			CreatedBy:    session.UserID,
			UpdatedBy:    session.UserID,
		}

		if _, err := repo.CreateExperience(ctx, params); err != nil {
			logger.ErrorContext(ctx, "Create Experience Failed", "service", serviceName, "error", err)
			return errs.ErrInternalServer
		}
		redpanda.RootNotificationProduce(
			ctx,
			session,
			session.UserID,
			utils.DatabaseTableExperience,
			utils.DatabaseMethodCreate,
			serviceName,
			client,
			rpClient,
			logger,
		)

		return nil
	}

	idx, idErr := utils.ParseUUID(ctx, id, serviceName, logger)
	if idErr != nil {
		return idErr
	}

	params := &repository.UpdateExperienceParams{
		Title:        title,
		CompanyName:  companyName,
		Location:     location,
		LocationType: locationType,
		StartDate:    startDate,
		EndDate:      endDate,
		Description:  description,
		UpdatedBy:    session.UserID,
		ID:           idx,
		UserID:       userID,
		UpdatedAt:    updatedAt,
	}

	affected, err := repo.UpdateExperience(ctx, params)
	if err != nil {
		logger.ErrorContext(ctx, "Update Education Failed", "service", serviceName, "error", err)
		return errs.ErrInternalServer
	}

	if affected != 1 {
		logger.WarnContext(
			ctx,
			"Concurrent experience update detected",
			"service",
			serviceName,
			"userID", userID,
			"id", id,
		)
		return errs.ErrConflict
	}

	redpanda.RootNotificationProduce(
		ctx,
		session,
		session.UserID,
		utils.DatabaseTableExperience,
		utils.DatabaseMethodUpdate,
		serviceName,
		client,
		rpClient,
		logger,
	)

	return nil
}

func updateProfile(
	ctx context.Context,
	session *utils.UserSession,
	userID uuid.UUID,
	req *profilev1.CreateUpdateProfile,
	updatedAt time.Time,
	serviceName string,
	repo repository.Querier,
	client valkey.Client,
	rpClient *kgo.Client,
	logger *slog.Logger,
) (*repository.Profile, error) {
	params := &repository.UpdateProfileParams{
		Name:      req.GetName(),
		Title:     req.GetTitle(),
		UpdatedBy: session.UserID,
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

	redpanda.RootNotificationProduce(
		ctx,
		session,
		session.UserID,
		utils.DatabaseTableProfile,
		utils.DatabaseMethodUpdate,
		serviceName,
		client,
		rpClient,
		logger,
	)

	return row, nil
}

func updateSocial(
	ctx context.Context,
	session *utils.UserSession,
	userID uuid.UUID,
	req *profilev1.UpdateSocial,
	serviceName string,
	repo repository.Querier,
	client valkey.Client,
	rpClient *kgo.Client,
	logger *slog.Logger,
) error {
	idx, idxErr := utils.ParseUUID(ctx, req.GetId(), serviceName, logger)
	if idxErr != nil {
		return idxErr
	}
	params := &repository.UpdateSocialParams{
		Username:  req.GetUsername(),
		UpdatedBy: session.UserID,
		ID:        idx,
		UserID:    userID,
		UpdatedAt: req.GetUpdatedAt().AsTime(),
	}

	affected, err := repo.UpdateSocial(ctx, params)
	if uErr := socialError(ctx, err, serviceName, "update", logger); uErr != nil {
		return uErr
	}

	if affected != 1 {
		logger.WarnContext(
			ctx,
			"Social record not found or concurrent modification",
			"service", serviceName,
			"id", idx.String(),
			"userID", userID.String(),
		)
		return errs.ErrConflict
	}
	redpanda.RootNotificationProduce(
		ctx,
		session,
		session.UserID,
		utils.DatabaseTableSocial,
		utils.DatabaseMethodUpdate,
		serviceName,
		client,
		rpClient,
		logger,
	)

	return nil
}

func (s *RootProfileService) createUpdateIcon(
	ctx context.Context,
	id string,
	req *rootProfilev1.CreateUpdateIcon,
	serviceName string,
	updatedAt time.Time,
) error {
	userSession := utils.UserSessionContext(ctx)

	name := req.GetName()

	siteHostname := req.GetSiteHostname()
	siteHostnameErr := utils.ValidateHostname(siteHostname, false)
	if siteHostnameErr != nil {
		s.cfg.Logger.WarnContext(ctx, "Invalid Site", "service", serviceName, "error", siteHostnameErr)
		return errs.ErrInvalidURL
	}

	siteSuffix := utils.StringValue(req.GetSiteSuffix())

	hostname := req.GetHostname()
	hostnameErr := utils.ValidateHostname(hostname, true)
	if hostnameErr != nil {
		s.cfg.Logger.WarnContext(ctx, "Invalid URL", "service", serviceName, "error", hostnameErr)
		return errs.ErrInvalidURL
	}

	suffix := req.GetSuffix()

	color := req.GetColor()

	if id == "" {
		params := &repository.CreateIconParams{
			Name:         name,
			SiteHostname: siteHostname,
			SiteSuffix:   siteSuffix,
			Hostname:     hostname,
			Suffix:       suffix,
			Color:        color,
			CreatedBy:    userSession.UserID,
			UpdatedBy:    userSession.UserID,
		}

		if _, err := s.cfg.Repository.CreateIcon(ctx, params); err != nil {
			if icErr := s.iconError(ctx, err, serviceName, "create"); icErr != nil {
				return icErr
			}
		}
		return nil
	}

	idx, idxErr := utils.ParseUUID(ctx, id, serviceName, s.cfg.Logger)
	if idxErr != nil {
		return idxErr
	}

	params := &repository.UpdateIconParams{
		Name:         name,
		SiteHostname: siteHostname,
		SiteSuffix:   siteSuffix,
		Hostname:     hostname,
		Suffix:       suffix,
		Color:        color,
		ID:           idx,
		UpdatedAt:    updatedAt,
	}

	affected, err := s.cfg.Repository.UpdateIcon(ctx, params)
	if icErr := s.iconError(ctx, err, serviceName, "update"); icErr != nil {
		return icErr
	}

	if affected != 1 {
		s.cfg.Logger.WarnContext(
			ctx,
			"Record not found or concurrent modification or already verified",
			"service", serviceName,
			"id", id,
		)
		return errs.ErrConflict
	}

	return nil
}

func (s *RootProfileService) iconError(
	ctx context.Context,
	err error,
	serviceName,
	method string,
) error {
	if err != nil {
		if pgErr, ok := errors.AsType[*pgconn.PgError](err); ok && pgErr.Code == pgerrcode.UniqueViolation {
			switch pgErr.ConstraintName {
			case utils.IconUniqueViolationSiteHostnameSuffix:
				s.cfg.Logger.WarnContext(
					ctx,
					"Icon Site with suffix already exists",
					"service",
					serviceName,
				)
				return errs.ErrUniqueViolation("Icon Site with suffix")
			case utils.IconUniqueViolationSiteHostnameNoSuffix:
				s.cfg.Logger.WarnContext(
					ctx,
					"Icon Site already exists",
					"service",
					serviceName,
				)
				return errs.ErrUniqueViolation("Icon Site")
			case utils.IconUniqueViolationHostnameSuffix:
				s.cfg.Logger.WarnContext(
					ctx,
					"Icon URL already exists",
					"service",
					serviceName,
				)
				return errs.ErrUniqueViolation("Icon URL")
			default:
				s.cfg.Logger.WarnContext(
					ctx,
					"Icon Name already exists",
					"service",
					serviceName,
				)
				return errs.ErrUniqueViolation("Icon Name")
			}
		}
		s.cfg.Logger.WarnContext(
			ctx,
			fmt.Sprintf("Failed to %s icon", method),
			"service",
			serviceName,
		)
		return errs.ErrInternalServer
	}
	return nil
}

func (s *GatewayProfileService) validateHostname(ctx context.Context, hostname, txt, serviceName string) error {
	lookupCtx, cancel := context.WithTimeout(ctx, lookupTimeout)
	defer cancel()

	txtRecords, txtRecordsErr := s.cfg.Resolver.LookupTXT(lookupCtx, hostname)
	if txtRecordsErr != nil || len(txtRecords) == 0 {
		s.cfg.Logger.WarnContext(
			ctx,
			"domain verification failed: host unreachable or no public TXT records found",
			"service",
			serviceName,
			"err",
			txtRecordsErr,
		)
		return errs.ErrNotFound("TXT")
	}

	if !slices.Contains(txtRecords, txt) {
		s.cfg.Logger.WarnContext(
			ctx,
			"domain ownership verification failed: matching token not found in TXT pool",
			"service", serviceName,
			"hostname", hostname,
			"err", txtRecordsErr,
		)
		return errs.ErrNotFound("TXT")
	}
	return nil
}

func deleteDatabase(
	ctx context.Context,
	id, userIDStr, serviceName string,
	updatedAt time.Time,
	table string,
	repo repository.Querier,
	client valkey.Client,
	rpClient *kgo.Client,
	logger *slog.Logger,
) error {
	userSession := utils.UserSessionContext(ctx)

	idx, idErr := utils.ParseUUID(ctx, id, serviceName, logger)
	if idErr != nil {
		return idErr
	}

	userID := userSession.UserID

	if userIDStr != "" {
		parsedID, parsedIDErr := utils.ParseUUID(ctx, userIDStr, serviceName, logger)
		if parsedIDErr != nil {
			return parsedIDErr
		}
		userID = parsedID
	}

	var affected int64
	var err error

	switch table {
	case utils.DatabaseTableSocial:
		params := &repository.DeleteSocialParams{
			ID:        idx,
			UserID:    userID,
			UpdatedAt: updatedAt,
		}

		affected, err = repo.DeleteSocial(ctx, params)

	case utils.DatabaseTableExperience:
		params := &repository.DeleteExperienceParams{
			ID:        idx,
			UserID:    userID,
			UpdatedAt: updatedAt,
		}

		affected, err = repo.DeleteExperience(ctx, params)

	case utils.DatabaseTableEducation:
		params := &repository.DeleteEducationParams{
			ID:        idx,
			UserID:    userID,
			UpdatedAt: updatedAt,
		}

		affected, err = repo.DeleteEducation(ctx, params)

	case utils.DatabaseTableIcon:
		params := &repository.DeleteIconParams{
			ID:        idx,
			UpdatedAt: updatedAt,
		}

		affected, err = repo.DeleteIcon(ctx, params)

	case utils.DatabaseTableNameserver:
		params := &repository.DeleteNameserverParams{
			ID:        idx,
			UpdatedAt: updatedAt,
		}

		affected, err = repo.DeleteNameserver(ctx, params)
	case utils.DatabaseTableTemplate:
		params := &repository.DeleteTemplateParams{
			ID:        idx,
			UpdatedAt: updatedAt,
		}

		affected, err = repo.DeleteTemplate(ctx, params)

	default:
		logger.WarnContext(ctx, "Invalid table", "service", serviceName, "table", table)
		return errs.ErrInternalServer
	}

	return deleteDB(
		ctx,
		affected,
		err,
		serviceName,
		id,
		table,
		userSession,
		userSession.UserID,
		client,
		rpClient,
		logger,
	)
}

func (s *RootProfileService) createUpdateTemplate(
	ctx context.Context,
	req *rootProfilev1.Template,
	id, serviceName string,
	updatedAt time.Time,
) error {
	userSession := utils.UserSessionContext(ctx)
	iconID, iconIDErr := utils.ParseUUID(ctx, req.GetIconId(), serviceName, s.cfg.Logger)
	if iconIDErr != nil {
		return iconIDErr
	}

	if id == "" {
		params := &repository.CreateTemplateParams{
			IconID:      iconID,
			Name:        req.GetName(),
			Description: req.GetDescription(),
			CreatedBy:   userSession.UserID,
			UpdatedBy:   userSession.UserID,
		}

		idx, err := s.cfg.Repository.CreateTemplate(ctx, params)
		if err != nil {
			if tErr := s.templateError(ctx, err, serviceName); tErr != nil {
				return tErr
			}
		}

		s.updateTemplateValkey(ctx, userSession, idx, req.GetName(), serviceName)

		return nil
	}

	idx, idxErr := utils.ParseUUID(ctx, id, serviceName, s.cfg.Logger)
	if idxErr != nil {
		return idxErr
	}

	params := &repository.UpdateTemplateParams{
		Name:        req.GetName(),
		IconID:      iconID,
		Description: req.GetDescription(),
		UpdatedBy:   userSession.UserID,
		ID:          idx,
		UpdatedAt:   updatedAt,
	}

	row, affectedErr := s.cfg.Repository.UpdateTemplate(ctx, params)
	if err := s.templateError(ctx, affectedErr, serviceName); err != nil {
		return err
	}

	s.updateTemplateValkey(ctx, userSession, row.ID, row.Name, serviceName)

	return nil
}

func (s *RootProfileService) templateError(ctx context.Context, err error, serviceName string) error {
	if err != nil {
		if pgErr, ok := errors.AsType[*pgconn.PgError](err); ok {
			switch pgErr.Code {
			case pgerrcode.UniqueViolation:
				s.cfg.Logger.WarnContext(ctx, "Template name already exists", "service", serviceName, "error", err)
				return errs.ErrUniqueViolation("Name")
			default:
				s.cfg.Logger.WarnContext(ctx, "Icon ID not exists", "service", serviceName, "error", err)
				return errs.ErrForeignKeyViolation("Icon ID")
			}
		}

		if errors.Is(err, pgx.ErrNoRows) {
			s.cfg.Logger.WarnContext(ctx, "Concurrent template update detected", "service", serviceName, "error", err)
			return errs.ErrConflict
		}
		s.cfg.Logger.WarnContext(
			ctx,
			"Failed to create template",
			"service",
			serviceName,
		)
		return errs.ErrInternalServer
	}
	return nil
}

func (s *RootProfileService) updateTemplateValkey(
	ctx context.Context,
	userSession *utils.UserSession,
	id uuid.UUID,
	name, serviceName string,
) {
	cmd := s.cfg.Client.B().Hset().
		Key(id.String()).
		FieldValue().
		FieldValue("name", strings.ToLower(name)).
		Build()

	if vkErr := s.cfg.Client.Do(ctx, cmd).Error(); vkErr != nil {
		s.cfg.Logger.ErrorContext(
			ctx,
			"Failed to set template payload in valkey",
			"service", serviceName,
			"name", name,
			"error", vkErr,
		)
	}

	redpanda.RootNotificationProduce(
		ctx,
		userSession,
		userSession.UserID,
		utils.DatabaseTableTemplate,
		utils.DatabaseMethodUpdate,
		serviceName,
		s.cfg.Client,
		s.cfg.Redpanda,
		s.cfg.Logger,
	)
}
