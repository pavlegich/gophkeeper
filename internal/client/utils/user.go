package utils

import (
	"context"
	"errors"

	errs "github.com/pavlegich/gophkeeper/internal/client/errors"
)

// DoWithRetryIfEmpty tries to implement function three times, if the input is empty.
func DoWithRetryIfEmpty(ctx context.Context, f func(ctx context.Context) error) error {
	var err error
	for i := 0; i < 3; i++ {
		err = f(ctx)
		if !errors.Is(err, errs.ErrEmptyInput) {
			return err
		}
	}
	return err
}
