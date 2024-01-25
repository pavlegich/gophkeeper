// Package repository contains repository object
// and methods for interaction between service and storage.
package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/pavlegich/gophkeeper/internal/server/domains/user"
	errs "github.com/pavlegich/gophkeeper/internal/server/errors"
)

// Repository contains storage objects.
type Repository struct {
	db *sql.DB
}

// NewUserRepository returns new repository object.
func NewUserRepository(db *sql.DB) *Repository {
	return &Repository{
		db: db,
	}
}

// GetUserByLogin gets by login from the storage and returns user object.
func (r *Repository) GetUserByLogin(ctx context.Context, login string) (*user.User, error) {
	err := r.db.PingContext(ctx)
	if err != nil {
		return nil, fmt.Errorf("GetUserByLogin: connection to database is died %w", err)
	}

	row := r.db.QueryRowContext(ctx, `SELECT id, login, password FROM users WHERE login = $1`, login)

	var storedUser user.User
	err = row.Scan(&storedUser.ID, &storedUser.Login, &storedUser.Password)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, fmt.Errorf("GetUserByLogin: scan row failed %w", errs.ErrUserNotFound)
	}
	if err != nil {
		return nil, fmt.Errorf("GetUserByLogin: scan row failed %w", err)
	}

	err = row.Err()
	if err != nil {
		return nil, fmt.Errorf("GetUserByLogin: row.Err %w", err)
	}

	return &storedUser, nil
}

// CreateUser saves new user data into the storage and returns user object.
func (r *Repository) CreateUser(ctx context.Context, u *user.User) (*user.User, error) {
	err := r.db.PingContext(ctx)
	if err != nil {
		return nil, fmt.Errorf("CreateUser: connection to database in died %w", err)
	}

	row := r.db.QueryRowContext(ctx, `INSERT INTO users (login, password) VALUES ($1, $2) 
	RETURNING id, login, password`, u.Login, u.Password)

	var storedUser user.User
	err = row.Scan(&storedUser.ID, &storedUser.Login, &storedUser.Password)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == pgerrcode.UniqueViolation {
			return nil, fmt.Errorf("CreateUser: %w", errs.ErrLoginBusy)
		}
		return nil, fmt.Errorf("CreateUser: insert into table failed %w", err)
	}

	err = row.Err()
	if err != nil {
		return nil, fmt.Errorf("CreateUser: row.Err %w", err)
	}

	return &storedUser, nil
}
