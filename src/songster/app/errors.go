package app

import (
	"errors"
	"fmt"
)

var (
	ErrBadRequest = errors.New("bad request")
	ErrNotFound   = errors.New("not found")
	ErrInternal   = errors.New("internal error")
)

func NewBadRequestError(err error, msg string) error {
	return wrapErrMsg(wrapErr(ErrBadRequest, err), msg)
}

func NewNotFoundError(err error, msg string) error {
	return wrapErrMsg(wrapErr(ErrNotFound, err), msg)
}

func NewInternalError(err error, msg string) error {
	return wrapErrMsg(wrapErr(ErrInternal, err), msg)
}

func IsBadRequestError(err error) bool {
	return errors.Is(err, ErrBadRequest)
}

func IsNotFoundError(err error) bool {
	return errors.Is(err, ErrNotFound)
}

func wrapErr(baseErr, err error) error {
	if err == nil {
		return baseErr
	}
	return fmt.Errorf("%w: %w", baseErr, err)
}

func wrapErrMsg(err error, msg string) error {
	if msg == "" {
		return err
	}
	return fmt.Errorf("%w: %s", err, msg)
}
