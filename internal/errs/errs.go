package errs

import (
	"fmt"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	ErrInternalServer   = status.Error(codes.Internal, "Internal Server Error")
	ErrCanceled         = status.Error(codes.Canceled, "Request canceled by client")
	ErrRequestTimeout   = status.Error(codes.DeadlineExceeded, "Request timeout exceeded")
	ErrUnauthenticated  = status.Error(codes.Unauthenticated, "Session Expired")
	ErrPermissionDenied = status.Error(codes.PermissionDenied, "Don't have access")
	ErrInvalidUserID    = status.Error(codes.InvalidArgument, "Invalid UserID")
	ErrInvalidURL       = status.Error(codes.InvalidArgument, "Invalid URL")
	ErrInvalidHostname  = status.Error(codes.InvalidArgument, "Invalid Hostname")
	ErrConflict         = status.Error(
		codes.FailedPrecondition,
		"Something went wrong, please refresh and try again",
	)
)

func ErrUniqueViolation(message string) error {
	return status.Error(codes.AlreadyExists, fmt.Sprintf("%s already exists", message))
}

func ErrForeignKeyViolation(message string) error {
	return status.Error(codes.InvalidArgument, fmt.Sprintf("%s doesnot exists", message))
}

func ErrNotFound(message string) error {
	return status.Error(codes.NotFound, fmt.Sprintf("%s not found", message))
}
