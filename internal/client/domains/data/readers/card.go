// Package readers contains objects for different data types
// and methods for interacting with them.
package readers

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/pavlegich/gophkeeper/internal/client/domains/rwmanager"
	errs "github.com/pavlegich/gophkeeper/internal/client/errors"
	"github.com/pavlegich/gophkeeper/internal/client/utils"
)

// CardDetails contains client card details.
type CardDetails struct {
	Number  int       `json:"number"`
	Owner   string    `json:"owner"`
	Expires time.Time `json:"expires"`
	CV      int       `json:"cv"`
}

// CardReader contains data for card reader object.
type CardReader struct {
	details *CardDetails
	rw      rwmanager.RWService
}

// NewCardReader creates and returns new card reader object.
func NewCardReader(ctx context.Context, rw rwmanager.RWService) *CardReader {
	return &CardReader{
		details: &CardDetails{},
		rw:      rw,
	}
}

// Read reads card details from the input, returns them in byte format.
func (r *CardReader) Read(ctx context.Context) ([]byte, error) {
	var err error

	// Read card number
	r.rw.Write(ctx, "Card number: ")
	numberString, err := r.rw.Read(ctx)
	if err != nil {
		return nil, fmt.Errorf("Read: couldn't read card number %w", err)
	}

	r.details.Number, err = strconv.Atoi(numberString)
	if err != nil {
		return nil, fmt.Errorf("Read: %w", errs.ErrInvalidCardNumber)
	}
	if !utils.IsValidCardNumber(r.details.Number) {
		return nil, fmt.Errorf("Read: %w", errs.ErrInvalidCardNumber)
	}

	// Read card expiration date
	r.rw.Write(ctx, "Card expiration date (MM/YY): ")
	dateString, err := r.rw.Read(ctx)
	if err != nil {
		return nil, fmt.Errorf("Read: couldn't read card expiration date %w", err)
	}
	r.details.Expires, err = time.Parse("01/06", dateString)
	if err != nil {
		return nil, fmt.Errorf("Read: %w", errs.ErrInvalidCardDate)
	}

	// Read card owner
	r.rw.Write(ctx, "Card owner: ")
	r.details.Owner, err = r.rw.Read(ctx)
	if err != nil {
		return nil, fmt.Errorf("Read: couldn't read card owner %w", err)
	}

	// Read CV
	r.rw.Write(ctx, "Card CV: ")
	cvString, err := r.rw.Read(ctx)
	if err != nil {
		return nil, fmt.Errorf("Read: couldn't read card CV %w", err)
	}
	r.details.CV, err = strconv.Atoi(cvString)
	if err != nil || r.details.CV > 999 || r.details.CV < 100 {
		return nil, fmt.Errorf("Read: %w", errs.ErrInvalidCardCV)
	}

	data, err := json.MarshalIndent(r.details, "", "   ")
	if err != nil {
		return nil, fmt.Errorf("Read: marshal card details failed %w", err)
	}

	return data, nil
}
