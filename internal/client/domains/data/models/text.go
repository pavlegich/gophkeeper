package models

import (
	"context"
	"errors"
	"fmt"

	"github.com/pavlegich/gophkeeper/internal/client/domains/rwmanager"
	errs "github.com/pavlegich/gophkeeper/internal/client/errors"
	"github.com/pavlegich/gophkeeper/internal/client/utils"
)

// ReadText reads text from the input, returns it in byte format.
func ReadText(ctx context.Context, rw rwmanager.RWService) ([]byte, error) {
	res := ""
	var in string
	var err error

	// Read text
	rw.Write(ctx, "Type text, and 'close' to indicate the end: ")
	for in != utils.Close {
		in, err = rw.Read(ctx)
		if err != nil && !errors.Is(err, errs.ErrEmptyInput) {
			return nil, fmt.Errorf("ReadText: couldn't read text %w", err)
		}
		if in != utils.Close {
			res += in
		}
	}

	return []byte(res), nil
}
