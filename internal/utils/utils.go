package utils

import (
	"context"
	"errors"
	"net"
	"net/url"
	"strings"
	"time"
	"uuid"

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
	PlatformsNameKey = "platforms_name_key"
	PlatformURLKey   = "platforms_url_key"
)

func ValidateURL(domain string) (string, error) {
	input := strings.ToLower(domain)

	parsed, parsedErr := url.Parse(input)
	if parsedErr != nil {
		return "", errors.New("malformed domain format string")
	}

	host := parsed.Host

	hostname, _, err := net.SplitHostPort(host)
	if err != nil {
		hostname = parsed.Hostname()
	}

	if hostname == "" {
		return "", errors.New("invalid or empty hostname in URL")
	}

	if hostname == "localhost" || strings.HasSuffix(hostname, ".local") || strings.HasSuffix(hostname, ".localhost") {
		return "", errors.New("localhost is not allowed")
	}

	if ip := net.ParseIP(hostname); ip != nil {
		return "", errors.New("IP addresses are not allowed")
	}

	return hostname, nil
}
