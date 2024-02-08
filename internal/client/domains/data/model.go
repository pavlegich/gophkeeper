// Package data contains object and methods for
// interacting with data on the client side.
package data

import "context"

// Data contains information about data object.
type Data struct {
	Name     string `json:"name"`
	Type     string `json:"type"`
	Data     []byte `json:"data"`
	Metadata []byte `json:"metadata"`
}

// Service describes methods related with data object.
type Service interface {
	CreateOrUpdate(ctx context.Context) error
	GetValue(ctx context.Context) error
	Delete(ctx context.Context) error
}

// DataReader describes methods related with object,
// that is able to read data from the input, store it
// and return in byte format.
type DataReader interface {
	Read(ctx context.Context) ([]byte, error)
}
