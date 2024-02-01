// Package utils contains additional methods for client.
package utils

import (
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

func GetKnownErr(err error) error {
	if errors.Is(err, errs.ErrUknownCommand) {
		return errs.ErrUknownCommand
	}
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
	if errors.Is(err, syscall.ECONNREFUSED) {
		return errs.ErrConnectionRefused
	}
	return nil
}
