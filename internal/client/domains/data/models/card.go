// Package models contains objects for different data types
// and methods for interacting with them.
package models

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

type CardDetails struct {
	Number  int       `json:"number"`
	Owner   string    `json:"owner"`
	Expires time.Time `json:"expires"`
	CV      int       `json:"cv"`
}

func ReadCardDetails(ctx context.Context, rw rwmanager.RWService) ([]byte, error) {
	c := &CardDetails{}
	var err error

	// Read card number
	rw.Write(ctx, "Card number: ")
	numberString, err := rw.Read(ctx)
	if err != nil {
		return nil, fmt.Errorf("ReadCardDetails: couldn't read card number %w", err)
	}

	c.Number, err = strconv.Atoi(numberString)
	if err != nil {
		return nil, fmt.Errorf("ReadCardDetails: %w", errs.ErrInvalidCardNumber)
	}
	if !utils.IsValidCardNumber(c.Number) {
		return nil, fmt.Errorf("ReadCardDetails: %w", errs.ErrInvalidCardNumber)
	}

	// Read card expiration date
	rw.Write(ctx, "Card expiration date (MM/YY): ")
	dateString, err := rw.Read(ctx)
	if err != nil {
		return nil, fmt.Errorf("ReadCardDetails: couldn't read card expiration date %w", err)
	}
	c.Expires, err = time.Parse("01/06", dateString)
	if err != nil {
		return nil, fmt.Errorf("ReadCardDetails: %w", errs.ErrInvalidCardDate)
	}

	// Read card owner
	rw.Write(ctx, "Card owner: ")
	c.Owner, err = rw.Read(ctx)
	if err != nil {
		return nil, fmt.Errorf("ReadCardDetails: %w", errs.ErrBadRequest)
	}
	if c.Owner == "" {
		return nil, fmt.Errorf("ReadCardDetails: %w", errs.ErrEmptyInput)
	}

	// Read cv
	rw.Write(ctx, "Card cv: ")
	cvString, err := rw.Read(ctx)
	if err != nil {
		return nil, fmt.Errorf("ReadCardDetails: couldn't read card owner %w", err)
	}
	c.CV, err = strconv.Atoi(cvString)
	if err != nil || c.CV > 999 || c.CV < 100 {
		return nil, fmt.Errorf("ReadCardDetails: %w", errs.ErrInvalidCardCV)
	}

	data, err := json.MarshalIndent(c, "", "   ")
	if err != nil {
		return nil, fmt.Errorf("ReadCardDetails: marshal card details failed %w", err)
	}

	return data, nil
}
