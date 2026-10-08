package utils

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"time"
	"uuid"

	"github.com/valkey-io/valkey-go"
	"golang.org/x/net/publicsuffix"
	"golang.org/x/sync/errgroup"
	"google.golang.org/protobuf/types/known/timestamppb"
	"google.golang.org/protobuf/types/known/wrapperspb"
	"neupaneanish.com.np/profile/internal/errs"
)

type ContextKey string

const (
	SessionKey       ContextKey = "user_session"
	DomainSessionKey ContextKey = "domain_user_session"

	DatabaseTableProfile    = "profile"
	DatabaseTableAbout      = "about"
	DatabaseTableEducation  = "education"
	DatabaseTableExperience = "experience"
	DatabaseTableIcon       = "icon"
	DatabaseTableSocial     = "social"
	DatabaseTableNameserver = "nameserver"
	DatabaseTableDomain     = "domain"
	DatabaseTableTemplate   = "template"

	DatabaseMethodCreate = "create"
	DatabaseMethodUpdate = "update"
	DatabaseMethodDelete = "delete"

	systemUsername                = "system"
	unknownUsername               = "unknown"
	RedpandaRootNotificationTopic = "root-notification"
)

type UserSession struct {
	UserID   uuid.UUID
	Username string
}

func UserSessionContext(ctx context.Context) *UserSession {
	session, _ := ctx.Value(SessionKey).(*UserSession)
	return session
}

func StringValue(wrapper *wrapperspb.StringValue) *string {
	if wrapper == nil {
		return nil
	}
	value := wrapper.GetValue()
	return &value
}

func TimestampValue(ts *timestamppb.Timestamp) *time.Time {
	if ts == nil {
		return nil
	}
	t := ts.AsTime()
	return &t
}

func TimestamppbValue(t *time.Time) *timestamppb.Timestamp {
	if t == nil {
		return nil
	}
	return timestamppb.New(*t)
}

func StringpbValue(s *string) *wrapperspb.StringValue {
	if s == nil {
		return nil
	}
	return wrapperspb.String(*s)
}

const (
	IconUniqueViolationSiteHostnameSuffix   = "unique_icons_site_hostname_with_suffix"
	IconUniqueViolationSiteHostnameNoSuffix = "unique_icons_site_hostname_no_suffix"
	IconUniqueViolationHostnameSuffix       = "unique_hostname_suffix"
)

func ValidateHostname(hostname string, subDomain bool) error {
	eTLD, icann := publicsuffix.PublicSuffix(hostname)
	if !icann || eTLD == hostname {
		return errors.New("invalid public suffix")
	}

	eTLDPlusOne, err := publicsuffix.EffectiveTLDPlusOne(hostname)
	if err != nil {
		return err
	}

	if subDomain {
		return nil
	}

	if hostname != eTLDPlusOne {
		return errors.New("subdomains not allowed")
	}

	return nil
}

type DomainUser struct {
	Key    string `json:"key"     valkey:",key"`
	UserID string `json:"user_id"`
}

type ExternalUserSession struct {
	UserID uuid.UUID
}

func DomainUserSessionContext(ctx context.Context) uuid.UUID {
	session, _ := ctx.Value(DomainSessionKey).(*ExternalUserSession)
	return session.UserID
}

func ParseUUID(ctx context.Context, userIDStr, serviceName string, logger *slog.Logger) (uuid.UUID, error) {
	userID, uuidErr := uuid.Parse(userIDStr)
	if uuidErr != nil {
		logger.WarnContext(ctx, "Failed to parse userID", "service", serviceName, "userID", userIDStr)
		return uuid.Nil(), errs.ErrInvalidUserID
	}
	return userID, nil
}

type RootNotification struct {
	ActorID       uuid.UUID
	UserID        uuid.UUID
	ActorUsername string
	UserUsername  string
	Table         string
	Method        string
}

func GetUsernames(
	ctx context.Context,
	createdBy, updatedBy uuid.UUID,
	session *UserSession,
	client valkey.Client,
	logger *slog.Logger,
) (string, string, error) {
	if createdBy == updatedBy && createdBy == uuid.Nil() {
		return systemUsername, systemUsername, nil
	}

	if session.UserID == createdBy && createdBy == updatedBy {
		return session.Username, session.Username, nil
	}

	var createdUsername, updateUsername string

	g, gCtx := errgroup.WithContext(ctx)

	g.Go(func() error {
		if session.UserID == createdBy {
			createdUsername = session.Username
			return nil
		}

		username, err := getUsername(gCtx, createdBy, client)
		if err != nil {
			return err
		}
		createdUsername = username
		return nil
	})

	g.Go(func() error {
		if session.UserID == updatedBy {
			updateUsername = session.Username
			return nil
		}

		username, err := getUsername(gCtx, updatedBy, client)
		if err != nil {
			return err
		}
		updateUsername = username
		return nil
	})

	if err := g.Wait(); err != nil {
		logger.ErrorContext(ctx, "Failed to get username", "error", err)
		return unknownUsername, unknownUsername, errs.ErrInternalServer
	}

	return createdUsername, updateUsername, nil
}

func getUsername(ctx context.Context, userID uuid.UUID, client valkey.Client) (string, error) {
	if userID == uuid.Nil() {
		return systemUsername, nil
	}

	key := fmt.Sprintf("username:%s", userID.String())
	cmd := client.B().Get().Key(key).Build()

	value, err := client.Do(ctx, cmd).ToString()
	if err != nil {
		if valkey.IsValkeyNil(err) {
			return "notfound", nil
		}
		return unknownUsername, err
	}
	return value, nil
}
