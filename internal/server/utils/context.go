// Package utils contains additional methods for server.
package utils

import (
	"context"
	"fmt"
)

type contextKey int

const (
	// List of const variables contains variables for
	// put values into and get values from the context.
	ContextIDKey contextKey = iota
)

// GetUserIDFromContext finds and returns user id from the context.
func GetUserIDFromContext(ctx context.Context) (int, error) {
	ctxValue := ctx.Value(ContextIDKey)
	if ctxValue == nil {
		return -1, fmt.Errorf("GetUserIDFromContext: get context value failed")
	}
	userID, ok := ctxValue.(int)
	if !ok {
		return -1, fmt.Errorf("GetUserIDFromContext: convert context value into integer failed")
	}
	return userID, nil
}
