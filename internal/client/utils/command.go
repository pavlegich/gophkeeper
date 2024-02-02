// Package utils contains additional methods for client.
package utils

import (
	"context"
	"errors"
	"syscall"

	"github.com/pavlegich/gophkeeper/internal/client/domains/rwmanager"
	errs "github.com/pavlegich/gophkeeper/internal/client/errors"
)

const (
	Greet          = "Welcome to GophKeeper!"
	Quit           = "quit"
	Success        = "success"
	Exit           = "exit"
	Close          = "close"
	UnexpectedQuit = "unexpected quit"
)

// GetKnownErr checks the error and returns it, if it is known.
func GetKnownErr(err error) error {
	if errors.Is(err, errs.ErrBadRequest) || errors.Is(err, errs.ErrUnknownStatusCode) {
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
	if errors.Is(err, errs.ErrEmptyInput) {
		return errs.ErrEmptyInput
	}
	if errors.Is(err, errs.ErrUnauthorized) {
		return errs.ErrUnauthorized
	}
	if errors.Is(err, syscall.ECONNREFUSED) {
		return errs.ErrConnectionRefused
	}
	if errors.Is(err, errs.ErrInvalidCardNumber) {
		return errs.ErrInvalidCardNumber
	}
	if errors.Is(err, errs.ErrInvalidDataType) {
		return errs.ErrInvalidDataType
	}
	if errors.Is(err, errs.ErrInvalidCardDate) {
		return errs.ErrInvalidCardDate
	}
	if errors.Is(err, errs.ErrInvalidCardCV) {
		return errs.ErrInvalidCardCV
	}
	if errors.Is(err, errs.ErrInvalidMetadata) {
		return errs.ErrInvalidMetadata
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

// DoWithRetryIfEmpty tries to implement function three times, if the input is empty.
func DoWithRetryIfEmpty(ctx context.Context, rw rwmanager.RWService, f func(ctx context.Context) error) error {
	var err error
	for i := 0; i < 3; i++ {
		err = f(ctx)
		if !errors.Is(err, errs.ErrEmptyInput) {
			return err
		}
		rw.WriteString(ctx, errs.ErrEmptyInput.Error())
	}
	return err
}
