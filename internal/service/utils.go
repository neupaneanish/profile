package service

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net"
	"slices"
	"time"
	"uuid"

	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"google.golang.org/protobuf/types/known/timestamppb"
	"neupaneanish.com.np/profile/internal/config"
	"neupaneanish.com.np/profile/internal/enum"
	"neupaneanish.com.np/profile/internal/errs"
	profilev1 "neupaneanish.com.np/profile/internal/protobuf/common/profile/v1"
	rootProfilev1 "neupaneanish.com.np/profile/internal/protobuf/root/profile/v1"
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
	serviceName, table, id, username string,
	actorID, userID uuid.UUID,
	logger *slog.Logger,
	redpanda *config.Redpanda,
) error {
	if err != nil {
		logger.ErrorContext(ctx, fmt.Sprintf("Delete %s Failed", table), "service", serviceName, "error", err)
		return errs.ErrInternalServer
	}

	if affected != 1 {
		logger.WarnContext(
			ctx,
			fmt.Sprintf("%s record not found or concurrent modification", table),
			"service", serviceName,
			"id", id,
		)
		return errs.ErrConflict
	}
	payload := utils.RedpandaRootEventNotificationPayload{
		ActorID:  actorID,
		Username: username,
		UserID:   userID,
		Message:  fmt.Sprintf("delete %s", table),
	}
	redpanda.Produce(ctx, utils.RedpandaRootEventNotifications, serviceName, payload)

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
			CreatedBy:     userID,
			UpdatedBy:     userID,
		}
		if _, err := repo.CreateEducation(ctx, params); err != nil {
			logger.ErrorContext(ctx, "Create Education Failed", "service", serviceName, "error", err)
			return errs.ErrInternalServer
		}
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
		UpdatedBy:     updatedBy,
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
	return nil
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
			CreatedBy:    userID,
			UpdatedBy:    userID,
		}

		if _, err := repo.CreateExperience(ctx, params); err != nil {
			logger.ErrorContext(ctx, "Create Experience Failed", "service", serviceName, "error", err)
			return errs.ErrInternalServer
		}
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
		UpdatedBy:    updatedBy,
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
	return nil
}

func updateProfile(
	ctx context.Context,
	userID, updatedBy uuid.UUID,
	req *profilev1.CreateUpdateProfile,
	updatedAt time.Time,
	serviceName, username string,
	repo repository.Querier,
	redpanda *config.Redpanda,
	logger *slog.Logger,
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

	payload := utils.RedpandaRootEventNotificationPayload{
		ActorID:  updatedBy,
		Username: username,
		UserID:   userID,
		Message:  "profile updated",
	}
	redpanda.Produce(ctx, utils.RedpandaRootEventNotifications, serviceName, payload)

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
	idx, idxErr := utils.ParseUUID(ctx, req.GetId(), serviceName, logger)
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

	site := req.GetSite()
	siteErr := utils.ValidateHostname(site, false)
	if siteErr != nil {
		s.cfg.Logger.WarnContext(ctx, "Invalid Site", "service", serviceName, "error", siteErr)
		return errs.ErrInvalidURL
	}

	siteSuffix := utils.StringValue(req.GetSiteSuffix())

	url := req.GetUrl()
	urlErr := utils.ValidateHostname(url, true)
	if urlErr != nil {
		s.cfg.Logger.WarnContext(ctx, "Invalid URL", "service", serviceName, "error", urlErr)
		return errs.ErrInvalidURL
	}

	slug := req.GetSlug()

	color := req.GetColor()

	if id == "" {
		params := &repository.CreateIconParams{
			Name:       name,
			Site:       site,
			SiteSuffix: siteSuffix,
			Url:        url,
			Slug:       slug,
			Color:      color,
			CreatedBy:  userSession.UserID,
			UpdatedBy:  userSession.UserID,
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
		Name:       name,
		Site:       site,
		SiteSuffix: siteSuffix,
		Url:        url,
		Slug:       slug,
		Color:      color,
		ID:         idx,
		UpdatedAt:  updatedAt,
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
			case utils.IconUniqueViolationSiteSuffix:
				s.cfg.Logger.WarnContext(
					ctx,
					"Icon Site with suffix already exists",
					"service",
					serviceName,
				)
				return errs.ErrUniqueViolation("Icon Site with suffix")
			case utils.IconUniqueViolationSiteNoSuffix:
				s.cfg.Logger.WarnContext(
					ctx,
					"Icon Site already exists",
					"service",
					serviceName,
				)
				return errs.ErrUniqueViolation("Icon Site")
			case utils.IconUniqueViolationURLSlug:
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

func (s *GatewayProfileService) validateDomain(ctx context.Context, fqdn, txt, ipAdd, serviceName string) error {
	lookupCtx, cancel := context.WithTimeout(ctx, lookupTimeout)
	defer cancel()

	txtRecords, txtRecordsErr := s.cfg.Resolver.LookupTXT(lookupCtx, fqdn)
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
			"domain", fqdn,
			"err", txtRecordsErr,
		)
		return errs.ErrNotFound("TXT")
	}

	hostIps, hostIpsErr := s.cfg.Resolver.LookupIP(lookupCtx, "ip", fqdn)
	if hostIpsErr != nil || len(hostIps) == 0 {
		s.cfg.Logger.WarnContext(
			ctx,
			"domain routing verification failed: host unresolved",
			"service", serviceName,
			"fqdn", fqdn,
			"err", hostIpsErr,
		)
		return errs.ErrNotFound("IP")
	}

	for _, ip := range hostIps {
		if ip.Equal(net.IP(ipAdd)) {
			return nil
		}
	}
	s.cfg.Logger.WarnContext(
		ctx,
		"User domain does not point to our platform IP structure",
		"service", serviceName,
		"fqdn", fqdn,
		"expected_platform_ip", ipAdd,
		"user_resolved_ips", hostIps,
	)
	return errs.ErrNotFound("IP")
}
