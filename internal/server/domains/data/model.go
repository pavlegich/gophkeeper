// Package data contains data object,
// service and repository for interacting between
// handlers and storage.
package data

import (
	"context"
	"time"
)

// Data contains data of data object.
type Data struct {
	ID        int       `db:"id" json:"id"`
	UserID    int       `db:"user_id" json:"user_id"`
	Name      string    `db:"name" json:"name"`
	Type      string    `db:"type" json:"type"`
	Data      []byte    `db:"data" json:"data"`
	Metadata  []byte    `db:"metadata" json:"metadata"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
}

// Service describes methods related with data object
// for communication between handlers and repositories.
type Service interface {
	Create(ctx context.Context, data *Data) error
	Unload(ctx context.Context, dType string, name string) (*Data, error)
	Edit(ctx context.Context, data *Data) error
	Delete(ctx context.Context, dType string, name string) error
}

// Repository describes methods related with data object
// for communication between services and database.
type Repository interface {
	GetDataByName(ctx context.Context, dType string, name string) (*Data, error)
	CreateData(ctx context.Context, data *Data) error
	UpdateData(ctx context.Context, data *Data) error
	DeleteDataByName(ctx context.Context, dType string, name string) error
}
