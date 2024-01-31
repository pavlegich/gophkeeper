package utils

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"syscall"
	"time"
)

// CheckStatusCode checks status code and returns the answer.
func CheckStatusCode(code int) error {
	switch code {
	case http.StatusOK:
		return nil
	case http.StatusBadRequest, http.StatusUnauthorized:
		return fmt.Errorf("please, check the entry and try again")
	case http.StatusConflict:
		return fmt.Errorf("already exists")
	case http.StatusInternalServerError:
		return fmt.Errorf("server failure, try again")
	case http.StatusNoContent:
		return fmt.Errorf("not exist")
	default:
		return fmt.Errorf("status code: %d", code)
	}
}

// GetRequestWithRetry requests with retries.
// If request is successful, returns response.
func GetRequestWithRetry(ctx context.Context, r *http.Request) (*http.Response, error) {
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
