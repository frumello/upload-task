package grpc

import (
	"github.com/pkg/errors"
)

func ErrUnexpectedly(err error) error {
	return errors.Wrapf(err, "failed unexpectedly while reading chunks from stream")
}

func ErrBadFileName(err error) error {
	return errors.Wrapf(err, "file name not provided in first chunk")
}

func ErrFileNameExist(err error, fileName string) error {
	return errors.Wrapf(err, "file with the name informed already exists: %s", fileName)
}

func ErrCreatingFile(err error) error {
	return errors.Wrapf(err, "failed unexpectedly while creating file")
}

func ErrFileType(err error) error {
	return errors.Wrapf(err, "failed unexpectedly while reading file type")
}

func ErrBadFileType(err error) error {
	return errors.Wrapf(err, "file type not accepted, please send a png, jpg or gif")
}
