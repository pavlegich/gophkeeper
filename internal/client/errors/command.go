package errors

import "errors"

var (
	ErrExit          = errors.New("exit requested")
	ErrUknownCommand = errors.New("unknown command")
	ErrEmptyInput    = errors.New("input is empty")
)
