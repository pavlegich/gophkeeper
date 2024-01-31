// Package utils contains additional methods for client.
package utils

import (
	"errors"

	errs "github.com/pavlegich/gophkeeper/internal/client/errors"
)

const (
	Greet   = "Welcome to GophKeeper!"
	Quit    = "quit"
	Success = "success"
	Exit    = "exit"
)

func IsKnownAndNotExitErr(err error) bool {
	if errors.Is(err, errs.ErrUknownCommand) {
		return true
	}
	return false
}
