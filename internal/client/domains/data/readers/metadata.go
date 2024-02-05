package readers

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

// MetadataReader contains metadata reader object data.
type MetadataReader struct {
	data map[string]string
	rw   rwmanager.RWService
}

// NewMetadataReader creates and returns new metadata object.
func NewMetadataReader(ctx context.Context, rw rwmanager.RWService) *MetadataReader {
	return &MetadataReader{
		data: make(map[string]string),
		rw:   rw,
	}
}

// Read reads metadata from the input, returns it in byte format.
func (r *MetadataReader) Read(ctx context.Context) ([]byte, error) {
	var in string
	var err error

	// Read text
	r.rw.Write(ctx, "Type metadata ('key : value'), and 'close' to indicate the end: ")
	for in != utils.Close {
		in, err = r.rw.Read(ctx)
		if errors.Is(err, errs.ErrEmptyInput) {
			break
		}
		if err != nil {
			res, _ := json.MarshalIndent(r.data, "", "   ")
			return res, fmt.Errorf("Read: couldn't read metadata key %w", err)
		}
		m := strings.Split(in, " : ")
		if in != utils.Close {
			if len(m) != 2 {
				res, _ := json.MarshalIndent(r.data, "", "   ")
				return res, fmt.Errorf("Read: %w", errs.ErrInvalidMetadata)
			}
			r.data[m[0]] = m[1]
		}
	}

	res, _ := json.MarshalIndent(r.data, "", "   ")

	return res, nil
}
