package errs

import (
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	ErrInternalServer        = status.Error(codes.Internal, "Internal Server Error")
	ErrCanceled              = status.Error(codes.Canceled, "Request canceled by client")
	ErrRequestTimeout        = status.Error(codes.DeadlineExceeded, "Request timeout exceeded")
	ErrUnauthenticated       = status.Error(codes.Unauthenticated, "Session Expired")
	ErrPermissionDenied      = status.Error(codes.PermissionDenied, "Don't have access")
	ErrProfileAlreadyExists  = status.Error(codes.AlreadyExists, "Profile already exists")
	ErrNotFound              = status.Error(codes.NotFound, "Not found")
	ErrInvalidUserID         = status.Error(codes.InvalidArgument, "Invalid UserID")
	ErrInvalidURL            = status.Error(codes.InvalidArgument, "Invalid URL")
	ErrPlatformNameExists    = status.Error(codes.AlreadyExists, "Platform name already exists")
	ErrPlatformURLExists     = status.Error(codes.AlreadyExists, "Platform url already exists")
	ErrPlatformLogoURLExists = status.Error(codes.AlreadyExists, "Platform logo url already exists")
	ErrConflict              = status.Error(
		codes.FailedPrecondition,
		"Something went wrong, please refresh and try again",
	)
)
