package errors

import "errors"

var (
	ErrBadRequest        = errors.New("please, check the entry and try again")
	ErrAlreadyExists     = errors.New("already exists")
	ErrServerInternal    = errors.New("server failure, try again")
	ErrNotExist          = errors.New("not exist")
	ErrUnknownStatusCode = errors.New("status code: ")
	ErrConnectionRefused = errors.New("server do not response, try again")
)
