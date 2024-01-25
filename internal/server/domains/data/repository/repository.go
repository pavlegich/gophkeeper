// Package repository contains repository object
// and methods for interaction between service and storage.
package repository

import (
	"context"
	"database/sql"

	"github.com/pavlegich/gophkeeper/internal/server/domains/data"
)

// Repository contains storage objects.
type Repository struct {
	db *sql.DB
}

// NewUserRepository returns new repository object.
func NewDataRepository(db *sql.DB) *Repository {
	return &Repository{
		db: db,
	}
}

func (r *Repository) GetByName(ctx context.Context, name string) (*data.Data, error) {
	return &data.Data{}, nil
}
func (r *Repository) Create(ctx context.Context, data *data.Data) error {
	return nil
}
func (r *Repository) Update(ctx context.Context, data *data.Data) error {
	return nil
}
func (r *Repository) Delete(ctx context.Context, name string) error {
	return nil
}
