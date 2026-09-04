package errs

import (
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	ErrInternalServer       = status.Error(codes.Internal, "Internal Server Error")
	ErrCanceled             = status.Error(codes.Canceled, "Request canceled by client")
	ErrRequestTimeout       = status.Error(codes.DeadlineExceeded, "Request timeout exceeded")
	ErrUnauthenticated      = status.Error(codes.Unauthenticated, "Session Expired")
	ErrPermissionDenied     = status.Error(codes.PermissionDenied, "Don't have access")
	ErrProfileAlreadyExists = status.Error(codes.AlreadyExists, "Profile already exists")
	ErrProfileNotFound      = status.Error(codes.AlreadyExists, "Profile not found")
	ErrInvalidUserID        = status.Error(codes.AlreadyExists, "Invalid UserID")
	ErrProfileConflict      = status.Error(
		codes.FailedPrecondition,
		"Something went wrong, please refresh and try again",
	)
)
