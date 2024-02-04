package models

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/pavlegich/gophkeeper/internal/client/domains/rwmanager"
	"github.com/pavlegich/gophkeeper/internal/client/domains/user"
)

// ReadCredentials reads credentials from the input, returns them in byte format.
func ReadCredentials(ctx context.Context, rw rwmanager.RWService) ([]byte, error) {
	u := &user.User{}
	var err error

	// Read login
	rw.Write(ctx, "Login: ")
	u.Login, err = rw.Read(ctx)
	if err != nil {
		return nil, fmt.Errorf("ReadCredentials: couldn't read login %w", err)
	}

	// Read password
	rw.Write(ctx, "Password: ")
	u.Password, err = rw.Read(ctx)
	if err != nil {
		return nil, fmt.Errorf("ReadCredentials: couldn't read password %w", err)
	}

	data, err := json.MarshalIndent(u, "", "   ")
	if err != nil {
		return nil, fmt.Errorf("ReadCredentials: marshal credentials failed %w", err)
	}

	return data, nil
}
