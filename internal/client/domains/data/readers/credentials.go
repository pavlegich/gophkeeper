package readers

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/pavlegich/gophkeeper/internal/client/domains/rwmanager"
	"github.com/pavlegich/gophkeeper/internal/client/domains/user"
)

// CredentialsReader contains data for credentials reader object.
type CredentialsReader struct {
	data *user.User
	rw   rwmanager.RWService
}

// NewCredentialsReader creates and returns new credentials reader.
func NewCredentialsReader(ctx context.Context, rw rwmanager.RWService) *CredentialsReader {
	return &CredentialsReader{
		data: &user.User{},
		rw:   rw,
	}
}

// Read reads credentials from the input, returns them in byte format.
func (r *CredentialsReader) Read(ctx context.Context) ([]byte, error) {
	var err error

	// Read login
	r.rw.Write(ctx, "Login: ")
	r.data.Login, err = r.rw.Read(ctx)
	if err != nil {
		return nil, fmt.Errorf("Read: couldn't read login %w", err)
	}

	// Read password
	r.rw.Write(ctx, "Password: ")
	r.data.Password, err = r.rw.Read(ctx)
	if err != nil {
		return nil, fmt.Errorf("Read: couldn't read password %w", err)
	}

	data, err := json.MarshalIndent(r.data, "", "   ")
	if err != nil {
		return nil, fmt.Errorf("Read: marshal credentials failed %w", err)
	}

	return data, nil
}
