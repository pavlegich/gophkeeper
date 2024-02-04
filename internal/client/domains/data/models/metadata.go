package models

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/pavlegich/gophkeeper/internal/client/domains/rwmanager"
	errs "github.com/pavlegich/gophkeeper/internal/client/errors"
	"github.com/pavlegich/gophkeeper/internal/client/utils"
)

// ReadMetadata reads metadata from the input, returns it in byte format.
func ReadMetadata(ctx context.Context, rw rwmanager.RWService) ([]byte, error) {
	meta := map[string]string{}
	var in string
	var err error

	// Read text
	rw.Write(ctx, "Type metadata ('key : value'), and 'close' to indicate the end: ")
	for in != utils.Close {
		in, err = rw.Read(ctx)
		if errors.Is(err, errs.ErrEmptyInput) {
			break
		}
		if err != nil {
			res, _ := json.MarshalIndent(meta, "", "   ")
			return res, fmt.Errorf("ReadMetadata: couldn't read metadata key %w", err)
		}
		m := strings.Split(in, " : ")
		if in != utils.Close {
			if len(m) != 2 {
				res, _ := json.MarshalIndent(meta, "", "   ")
				return res, fmt.Errorf("ReadMetadata: %w", errs.ErrInvalidMetadata)
			}
			meta[m[0]] = m[1]
		}
	}

	res, _ := json.MarshalIndent(meta, "", "   ")

	return res, nil
}
