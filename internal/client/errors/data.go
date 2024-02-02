// Package errors contains error decriptions for client.
package errors

import "errors"

var (
	ErrInvalidDataType   = errors.New("invalid data type")
	ErrInvalidCardNumber = errors.New("invalid card number")
	ErrInvalidCardDate   = errors.New("invalid card expiration date")
	ErrInvalidCardCV     = errors.New("invalid card cv")
)
