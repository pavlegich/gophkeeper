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
	ContextActionKey contextKey = iota
)

// GetActionFromContext finds and returns user id from the context.
func GetActionFromContext(ctx context.Context) (string, error) {
	ctxValue := ctx.Value(ContextActionKey)
	if ctxValue == nil {
		return "", fmt.Errorf("GetActionFromContext: get context value failed")
	}
	userID, ok := ctxValue.(string)
	if !ok {
		return "", fmt.Errorf("GetActionFromContext: convert context value into string failed")
	}
	return userID, nil
}
