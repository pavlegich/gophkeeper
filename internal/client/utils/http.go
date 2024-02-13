package utils

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"syscall"
	"time"

	errs "github.com/pavlegich/gophkeeper/internal/client/errors"
)

// CheckStatusCode checks status code and returns the answer.
func CheckStatusCode(code int) error {
	switch code {
	case http.StatusOK:
		return nil
	case http.StatusBadRequest:
		return errs.ErrBadRequest
	case http.StatusUnauthorized:
		return errs.ErrUnauthorized
	case http.StatusConflict:
		return errs.ErrAlreadyExists
	case http.StatusInternalServerError:
		return errs.ErrServerInternal
	case http.StatusNoContent:
		return errs.ErrNotExist
	default:
		return fmt.Errorf("%w%d", errs.ErrUnknownStatusCode, code)
	}
}

// DoRequestWithRetry requests with retries.
// If request is successful, returns response.
func DoRequestWithRetry(ctx context.Context, r *http.Request) (*http.Response, error) {
	var err error = nil
	var resp *http.Response

	intervals := []time.Duration{0, time.Second, 3 * time.Second, 5 * time.Second}
	for _, interval := range intervals {
		time.Sleep(interval)
		resp, err = http.DefaultClient.Do(r)
		if !errors.Is(err, syscall.ECONNREFUSED) {
			break
		}
	}
	if err != nil {
		return nil, fmt.Errorf("GetRequestWithRetry: request failed %w", err)
	}
	return resp, nil
}
