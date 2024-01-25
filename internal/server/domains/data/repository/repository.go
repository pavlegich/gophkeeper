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
	"github.com/pavlegich/gophkeeper/internal/server/domains/data"
	errs "github.com/pavlegich/gophkeeper/internal/server/errors"
	"github.com/pavlegich/gophkeeper/internal/utils"
)

// Repository contains storage objects.
type Repository struct {
	db *sql.DB
}

// NewDataRepository returns new repository object.
func NewDataRepository(db *sql.DB) *Repository {
	return &Repository{
		db: db,
	}
}

// GetDataByName gets data by name from the storage and returns data object.
func (r *Repository) GetDataByName(ctx context.Context, name string) (*data.Data, error) {
	err := r.db.PingContext(ctx)
	if err != nil {
		return nil, fmt.Errorf("GetDataByName: connection to database is died %w", err)
	}

	userID, err := utils.GetUserIDFromContext(ctx)
	if err != nil {
		return nil, fmt.Errorf("GetDataByName: couldn't read user id from the context %w", err)
	}
	dType, err := utils.GetTypeFromContext(ctx)
	if err != nil {
		return nil, fmt.Errorf("GetDataByName: couldn't read data type from the context %w", err)
	}

	row := r.db.QueryRowContext(ctx, `SELECT id, user_id, name, data_type, data, created_at, metadata 
	FROM data WHERE user_id = $1, data_type = $2, name = $3`, userID, dType, name)

	var storedData data.Data
	err = row.Scan(&storedData.ID, &storedData.UserID, &storedData.Name, &storedData.Type,
		&storedData.Data, &storedData.CreatedAt, &storedData.Metadata)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, fmt.Errorf("GetDataByName: scan row failed %w", errs.ErrDataNotFound)
	}
	if err != nil {
		return nil, fmt.Errorf("GetDataByName: scan row failed %w", err)
	}

	err = row.Err()
	if err != nil {
		return nil, fmt.Errorf("GetUserByLogin: row.Err %w", err)
	}

	return &storedData, nil
}

// CreateData saves new data object into the storage.
func (r *Repository) CreateData(ctx context.Context, data *data.Data) error {
	err := r.db.PingContext(ctx)
	if err != nil {
		return fmt.Errorf("CreateData: connection to database in died %w", err)
	}

	tx, err := r.db.Begin()
	if err != nil {
		return fmt.Errorf("CreateData: begin transaction failed %w", err)
	}
	defer tx.Rollback()

	row := tx.QueryRowContext(ctx, `SELECT id FROM data WHERE user_id = $1 AND name = $2 
	AND data_type = $3`, data.UserID, data.Name, data.Type)
	var id int
	if err := row.Scan(&id); !errors.Is(err, sql.ErrNoRows) {
		if err == nil {
			return fmt.Errorf("CreateData: %w", errs.ErrDataAlreadyUpload)
		}
		return fmt.Errorf("CreateData: scan data row with id failed %w", err)
	}

	_, err = tx.ExecContext(ctx, `INSERT INTO data (user_id, name, data_type, data, metadata) 
	VALUES ($1, $2, $3, $4, $5)`, data.UserID, data.Name, data.Type, data.Data, data.Metadata)

	if err != nil {
		var pgErr *pgconn.PgError
		fmt.Println(pgErr.Code)
		if errors.As(err, &pgErr) && pgErr.Code == pgerrcode.UniqueViolation {
			return fmt.Errorf("CreateUser: %w", errs.ErrDataAlreadyUpload)
		}
		return fmt.Errorf("CreateData: insert data failed %w", err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("CreateData: commit transaction failed %w", err)
	}

	return nil
}

// UpdateData updates user data in storage.
func (r *Repository) UpdateData(ctx context.Context, data *data.Data) error {
	err := r.db.PingContext(ctx)
	if err != nil {
		return fmt.Errorf("UpdateData: connection to database in died %w", err)
	}

	_, err = r.db.ExecContext(ctx, `UPDATE data SET data = $1, metadata = $2 
	WHERE user_id = $3, name = $4, data_type = $5`,
		data.Data, data.Metadata, data.UserID, data.Name, data.Type)
	if err != nil {
		return fmt.Errorf("UpdateData: update table failed %w", err)
	}

	if errors.Is(err, sql.ErrNoRows) {
		return fmt.Errorf("UpdateData: nothing to update, %w", errs.ErrDataNotFound)
	}

	return nil
}

// DeleteDataByName deletes requested data by it's name.
func (r *Repository) DeleteDataByName(ctx context.Context, name string) error {
	err := r.db.PingContext(ctx)
	if err != nil {
		return fmt.Errorf("DeleteData: connection to database in died %w", err)
	}

	userID, err := utils.GetUserIDFromContext(ctx)
	if err != nil {
		return fmt.Errorf("DeleteDataByName: couldn't read user id from the context %w", err)
	}
	dType, err := utils.GetTypeFromContext(ctx)
	if err != nil {
		return fmt.Errorf("DeleteDataByName: couldn't read data type from the context %w", err)
	}

	_, err = r.db.ExecContext(ctx, `DELETE FROM data WHERE user_id = $1, data_type = $2, name = $3`,
		userID, dType, name)
	if err != nil {
		return fmt.Errorf("DeleteDataByName: couldn't delete data from the storage %w", err)
	}

	return nil
}
