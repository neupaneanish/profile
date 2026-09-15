package utils

import (
	"context"
	"errors"
	"log/slog"
	"net/netip"
	"time"
	"uuid"

	"golang.org/x/net/publicsuffix"
	"google.golang.org/protobuf/types/known/timestamppb"
	"google.golang.org/protobuf/types/known/wrapperspb"
	"neupaneanish.com.np/profile/internal/errs"
)

type ContextKey string

const (
	SessionKey       ContextKey = "user_session"
	DomainSessionKey ContextKey = "domain_user_session"
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
	IconUniqueViolationSiteSuffix   = "unique_icons_site_with_suffix"
	IconUniqueViolationSiteNoSuffix = "unique_icons_site_no_suffix"
	IconUniqueViolationURLSlug      = "unique_url_slug"
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

func ValidateIP(ip string) error {
	addr, err := netip.ParseAddr(ip)
	if err != nil {
		return err
	}

	if addr.IsPrivate() || addr.IsLoopback() || addr.IsLinkLocalUnicast() {
		return errors.New("private IP not allowed")
	}

	return nil
}

type DomainUser struct {
	Key    string `json:"key"     valkey:",key"`
	UserID string `json:"user_id"`
}

const (
	DomainUserSessionKey = "domain:user:session"
)

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
