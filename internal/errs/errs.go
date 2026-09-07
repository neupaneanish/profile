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
	ErrConflict         = status.Error(
		codes.FailedPrecondition,
		"Something went wrong, please refresh and try again",
	)
)

func ErrUniqueViolation(name string) error {
	return status.Error(codes.AlreadyExists, fmt.Sprintf("%s already exists", name))
}

func ErrForeignKeyViolation(name string) error {
	return status.Error(codes.InvalidArgument, fmt.Sprintf("%s doesnot exists", name))
}

func ErrNotFound(name string) error {
	return status.Error(codes.NotFound, fmt.Sprintf("%s not found", name))
}
