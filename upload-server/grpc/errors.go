package grpc

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func ErrUnexpectedly(err error) error {
	return status.Errorf(codes.Internal, "failed unexpectedly: %v", err)
}

func ErrBadFileName() error {
	return status.Errorf(codes.InvalidArgument, "file name not provided in first chunk")
}

func ErrFileNameExist(fileName string) error {
	return status.Errorf(codes.NotFound, "file with the name informed already exists: %s", fileName)
}

func ErrBadFileType() error {
	return status.Errorf(codes.InvalidArgument, "file type not accepted, please send a png, jpg or gif")
}

func ContextError(ctx context.Context) error {
	switch ctx.Err() {
	case context.Canceled:
		return status.Error(codes.Canceled, "request is canceled")
	case context.DeadlineExceeded:
		return status.Error(codes.DeadlineExceeded, "deadline is exceeded")
	default:
		return nil
	}
}
