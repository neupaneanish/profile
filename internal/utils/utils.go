package utils

import (
	"context"
	"errors"
	"net/netip"
	"time"
	"uuid"

	"golang.org/x/net/publicsuffix"
	"google.golang.org/protobuf/types/known/timestamppb"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

type ContextKey string

const SessionKey ContextKey = "user_session"

type UserSession struct {
	UserID uuid.UUID
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

func ValidateURL(domain string, subDomain bool) error {
	eTLD, icann := publicsuffix.PublicSuffix(domain)
	if !icann || eTLD == domain {
		return errors.New("invalid public suffix")
	}

	eTLDPlusOne, err := publicsuffix.EffectiveTLDPlusOne(domain)
	if err != nil {
		return err
	}

	if subDomain {
		return nil
	}

	if domain != eTLDPlusOne {
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
