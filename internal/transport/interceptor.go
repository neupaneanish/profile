package transport

import (
	"context"
	"errors"
	"log/slog"
	"time"

	"uuid"

	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"neupaneanish.com.np/profile/internal/enum"

	"neupaneanish.com.np/profile/internal/errs"
	"neupaneanish.com.np/profile/internal/utils"
)

func LoggerInterceptor(logger *slog.Logger) logging.Logger {
	return logging.LoggerFunc(func(ctx context.Context, level logging.Level, msg string, fields ...any) {
		logger.Log(ctx, slog.Level(level), msg, fields...)
	})
}

func UnaryTimeoutInterceptor(defaultTimeout time.Duration) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, _ *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
		if _, ok := ctx.Deadline(); ok {
			return handler(ctx, req)
		}

		newCtx, cancel := context.WithTimeout(ctx, defaultTimeout)
		defer cancel()

		resp, err := handler(newCtx, req)

		if err != nil && newCtx.Err() != nil {
			if errors.Is(newCtx.Err(), context.DeadlineExceeded) {
				return nil, errs.ErrRequestTimeout
			}
			if errors.Is(newCtx.Err(), context.Canceled) {
				return nil, errs.ErrCanceled
			}
		}
		return resp, err
	}
}

func StreamTimeoutInterceptor(maxDuration time.Duration) grpc.StreamServerInterceptor {
	return func(srv any, ss grpc.ServerStream, _ *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		ctx := ss.Context()

		if _, ok := ctx.Deadline(); !ok {
			var cancel context.CancelFunc
			ctx, cancel = context.WithTimeout(ctx, maxDuration)
			defer cancel()

			ss = &WrappedTimeoutStream{
				ServerStream:  ss,
				StreamContext: ctx,
			}
		}

		err := handler(srv, ss)

		if ctx.Err() != nil {
			if errors.Is(ctx.Err(), context.DeadlineExceeded) {
				return errs.ErrRequestTimeout
			}
			if errors.Is(ctx.Err(), context.Canceled) {
				return errs.ErrCanceled
			}
		}
		return err
	}
}

type WrappedTimeoutStream struct {
	grpc.ServerStream

	StreamContext context.Context
}

func (w *WrappedTimeoutStream) Context() context.Context {
	return w.StreamContext
}

func AuthInterceptor(
	ctx context.Context,
	logger *slog.Logger,
	external, gateway, root map[string]struct{},
) (context.Context, error) {
	serviceName := "Interceptor"
	fullMethod, ok := grpc.Method(ctx)
	if !ok {
		return ctx, errs.ErrInternalServer
	}

	_, isExternalEndpoint := external[fullMethod]
	_, isGatewayEndpoint := gateway[fullMethod]
	_, isRootEndpoint := root[fullMethod]

	switch {
	case isGatewayEndpoint:
		return authContext(ctx, logger, enum.UserRoleUser, serviceName)
	case isRootEndpoint:
		return authContext(ctx, logger, enum.UserRoleRoot, serviceName)
	case isExternalEndpoint:
		return externalContext(ctx, logger, serviceName)
	default:
		return ctx, errs.ErrPermissionDenied
	}
}

func metadataDetail(ctx context.Context, header string) string {
	value := metadata.ValueFromIncomingContext(ctx, header)
	if len(value) == 0 {
		return ""
	}
	return value[0]
}

func authContext(
	ctx context.Context,
	logger *slog.Logger,
	role enum.UserRole,
	serviceName string,
) (context.Context, error) {
	xUserID := metadataDetail(ctx, "x-user-id")
	xRole := metadataDetail(ctx, "x-role")
	xJti := metadataDetail(ctx, "x-jti")
	xUsername := metadataDetail(ctx, "x-username")

	hasUserDetail := xUserID != "" && xRole != "" && xJti != "" && xUsername != ""
	if !hasUserDetail {
		logger.WarnContext(ctx,
			"missing required gateway headers",
			"service", serviceName,
			"user_id", xUserID,
			"role", xRole,
			"username", xUsername,
		)
		return ctx, errs.ErrUnauthenticated
	}

	userID, userIDErr := uuid.Parse(xUserID)
	if userIDErr != nil {
		logger.ErrorContext(ctx, "Invalid userID", "service", serviceName, "userID", xUserID, "error", userIDErr)
		return ctx, errs.ErrUnauthenticated
	}

	if !enum.UserRole(xRole).Valid() {
		logger.WarnContext(ctx, "Invalid role", "service", serviceName, "role", xRole)
		return ctx, errs.ErrUnauthenticated
	}

	if enum.UserRole(xRole) != role {
		logger.WarnContext(ctx, "Missing permission", "service", serviceName, "role", xRole)
		return ctx, errs.ErrPermissionDenied
	}

	ctx = logging.InjectFields(ctx, logging.Fields{
		"user_id", userID.String(),
		"role", xRole,
		"jti", xJti,
		"username", xUsername,
	})

	return context.WithValue(
		ctx, utils.SessionKey,
		&utils.UserSession{
			UserID:   userID,
			Username: xUsername,
		},
	), nil
}

func externalContext(ctx context.Context, logger *slog.Logger, serviceName string) (context.Context, error) {
	xUserID := metadataDetail(ctx, "x-user-id")
	xHostname := metadataDetail(ctx, "x-hostname")
	if xUserID == "" || xHostname == "" {
		logger.WarnContext(ctx,
			"missing required gateway headers",
			"service", serviceName,
			"userID", xUserID,
			"hostname", xHostname,
		)
		return ctx, errs.ErrNotFound("Host")
	}

	userID, userIDErr := utils.ParseUUID(ctx, xUserID, serviceName, logger)
	if userIDErr != nil {
		return ctx, errs.ErrNotFound("Host")
	}

	ctx = logging.InjectFields(ctx, logging.Fields{
		"user_id", xUserID,
		"hostname", xHostname,
	})

	return context.WithValue(
		ctx, utils.DomainSessionKey,
		&utils.ExternalUserSession{
			UserID: userID,
		},
	), nil
}
