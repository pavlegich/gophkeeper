package readers

import (
	"context"
	"errors"
	"fmt"

	"github.com/pavlegich/gophkeeper/internal/client/domains/rwmanager"
	errs "github.com/pavlegich/gophkeeper/internal/client/errors"
	"github.com/pavlegich/gophkeeper/internal/client/utils"
)

// TextReader contains text reader object data.
type TextReader struct {
	text string
	rw   rwmanager.RWService
}

// NewTextReader creates and returns new text object.
func NewTextReader(ctx context.Context, rw rwmanager.RWService) *TextReader {
	return &TextReader{
		text: "",
		rw:   rw,
	}
}

// Read reads text from the input, returns it in byte format.
func (r *TextReader) Read(ctx context.Context) ([]byte, error) {
	var in string
	var err error

	// Read text
	r.rw.Write(ctx, "Type text, and 'close' to indicate the end: ")
	for in != utils.Close {
		in, err = r.rw.Read(ctx)
		if err != nil && !errors.Is(err, errs.ErrEmptyInput) {
			return nil, fmt.Errorf("Read: couldn't read text %w", err)
		}
		if in != utils.Close {
			r.text += in
		}
	}

	return []byte(r.text), nil
}
