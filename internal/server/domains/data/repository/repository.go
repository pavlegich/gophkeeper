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
	"github.com/pavlegich/gophkeeper/internal/server/utils"
)

// Repository contains storage objects.
type Repository struct {
	db *sql.DB
}

// NewDataRepository returns new repository object.
func NewDataRepository(ctx context.Context, db *sql.DB) *Repository {
	return &Repository{
		db: db,
	}
}

// GetDataByName gets data by name from the storage and returns data object.
func (r *Repository) GetDataByName(ctx context.Context, dType string, name string) (*data.Data, error) {
	userID, err := utils.GetUserIDFromContext(ctx)
	if err != nil {
		return nil, fmt.Errorf("GetDataByName: couldn't read user id from the context %w", err)
	}

	row := r.db.QueryRowContext(ctx, `SELECT id, user_id, name, data_type, data, created_at, metadata 
	FROM data WHERE user_id = $1 AND data_type = $2 AND name = $3`, userID, dType, name)

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
		return nil, fmt.Errorf("GetDataByName: row.Err %w", err)
	}

	return &storedData, nil
}

// CreateData saves new data object into the storage.
func (r *Repository) CreateData(ctx context.Context, d *data.Data) error {
	tx, err := r.db.Begin()
	if err != nil {
		return fmt.Errorf("CreateData: begin transaction failed %w", err)
	}
	defer tx.Rollback()

	row := tx.QueryRowContext(ctx, `SELECT id FROM data WHERE user_id = $1 AND name = $2 
	AND data_type = $3`, d.UserID, d.Name, d.Type)
	var id int
	if err := row.Scan(&id); !errors.Is(err, sql.ErrNoRows) {
		if err == nil {
			return fmt.Errorf("CreateData: %w", errs.ErrDataAlreadyUpload)
		}
		return fmt.Errorf("CreateData: scan data row with id failed %w", err)
	}

	_, err = tx.ExecContext(ctx, `INSERT INTO data (user_id, name, data_type, data, metadata) 
	VALUES ($1, $2, $3, $4, $5)`, d.UserID, d.Name, d.Type, d.Data, d.Metadata)

	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == pgerrcode.UniqueViolation {
			return fmt.Errorf("CreateData: %w", errs.ErrDataAlreadyUpload)
		}
		return fmt.Errorf("CreateData: insert data failed %w", err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("CreateData: commit transaction failed %w", err)
	}

	return nil
}

// UpdateData updates user data in storage.
func (r *Repository) UpdateData(ctx context.Context, d *data.Data) error {
	res, err := r.db.ExecContext(ctx, `UPDATE data SET data = $1, metadata = $2 
	WHERE user_id = $3 AND name = $4 AND data_type = $5`,
		d.Data, d.Metadata, d.UserID, d.Name, d.Type)
	if err != nil {
		return fmt.Errorf("UpdateData: update table failed %w", err)
	}

	rowsCount, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("UpdateData: couldn't get rows affected %w", err)
	}
	if rowsCount == 0 {
		return fmt.Errorf("UpdateData: nothing to update, %w", errs.ErrDataNotFound)
	}

	return nil
}

// DeleteDataByName deletes requested data by it's name.
func (r *Repository) DeleteDataByName(ctx context.Context, dType string, name string) error {
	userID, err := utils.GetUserIDFromContext(ctx)
	if err != nil {
		return fmt.Errorf("DeleteDataByName: couldn't read user id from the context %w", err)
	}

	res, err := r.db.ExecContext(ctx, `DELETE FROM data WHERE user_id = $1 AND data_type = $2 AND name = $3`,
		userID, dType, name)
	if err != nil {
		return fmt.Errorf("DeleteDataByName: couldn't delete data from the storage %w", err)
	}

	rowsCount, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("DeleteDataByName: couldn't get rows affected %w", err)
	}
	if rowsCount == 0 {
		return fmt.Errorf("DeleteDataByName: nothing to delete, %w", errs.ErrDataNotFound)
	}

	return nil
}
