// Package data contains data object,
// service and repository for interacting between
// handlers and storage.
package data

import (
	"context"
	"time"
)

type Data struct {
	ID        int       `db:"id" json:"id"`
	UserID    int       `db:"user_id" json:"user_id"`
	Name      string    `db:"name" json:"name"`
	Type      string    `db:"type" json:"type"`
	Data      string    `db:"data" json:"data"`
	Metadata  string    `db:"metadata" json:"metadata"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
}

type Service interface {
	Create(ctx context.Context, data *Data) (*Data, error)
	Unload(ctx context.Context, name string) (*Data, error)
	Edit(ctx context.Context, data *Data) error
	Delete(ctx context.Context, name string) error
}

type Repository interface {
	GetByName(ctx context.Context, name string) (*Data, error)
	Create(ctx context.Context, data *Data) error
	Update(ctx context.Context, data *Data) error
	Delete(ctx context.Context, name string) error
}
