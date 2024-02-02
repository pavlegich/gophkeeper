// Package utils contains additional methods for client.
package utils

import (
	"context"
	"errors"
	"syscall"

	errs "github.com/pavlegich/gophkeeper/internal/client/errors"
)

const (
	Greet          = "Welcome to GophKeeper!"
	Quit           = "quit"
	Success        = "success"
	Exit           = "exit"
	UnexpectedQuit = "unexpected quit"
)

// GetKnownErr checks the error and returns it, if it is known.
func GetKnownErr(err error) error {
	if errors.Is(err, errs.ErrBadRequest) {
		return errs.ErrBadRequest
	}
	if errors.Is(err, errs.ErrAlreadyExists) {
		return errs.ErrAlreadyExists
	}
	if errors.Is(err, errs.ErrServerInternal) {
		return errs.ErrServerInternal
	}
	if errors.Is(err, errs.ErrNotExist) {
		return errs.ErrNotExist
	}
	if errors.Is(err, errs.ErrUnknownStatusCode) {
		return errs.ErrUnknownStatusCode
	}
	if errors.Is(err, errs.ErrEmptyInput) {
		return errs.ErrEmptyInput
	}
	if errors.Is(err, errs.ErrUnauthorized) {
		return errs.ErrUnauthorized
	}
	if errors.Is(err, syscall.ECONNREFUSED) {
		return errs.ErrConnectionRefused
	}
	return nil
}

// DoWithRetryIfUnknown tries to implement function three times, if the input is empty.
func DoWithRetryIfUnknown(ctx context.Context, f func(ctx context.Context) error) error {
	var err error
	for i := 0; i < 3; i++ {
		err = f(ctx)
		if !errors.Is(err, errs.ErrUnknownCommand) {
			return err
		}
	}
	return err
}
