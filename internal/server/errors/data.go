package errors

import "errors"

var (
	ErrDataNotFound      = errors.New("data not found for this user")
	ErrDataAlreadyUpload = errors.New("data already uploaded by this user")
	ErrDataTypeIncorrect = errors.New("incorrect data type")
)
