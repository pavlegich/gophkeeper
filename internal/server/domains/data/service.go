package data

import (
	"context"
	"fmt"

	errs "github.com/pavlegich/gophkeeper/internal/server/errors"
	"github.com/pavlegich/gophkeeper/internal/utils"
)

// DataService contatins objects for user service.
type DataService struct {
	repo Repository
}

// NewDataService returns new data service.
func NewDataService(repo Repository) *DataService {
	return &DataService{
		repo: repo,
	}
}

// Create upload new data into the storage.
func (s *DataService) Create(ctx context.Context, data *Data) error {
	if !utils.IsCorrectDataType(data.Type) {
		return fmt.Errorf("Create: %w", errs.ErrDataTypeIncorrect)
	}
	err := s.repo.CreateData(ctx, data)
	if err != nil {
		return fmt.Errorf("Create: create data failed %w", err)
	}
	return nil
}

// Unload unloads data by name and returns data object.
func (s *DataService) Unload(ctx context.Context, name string) (*Data, error) {
	d, err := s.repo.GetDataByName(ctx, name)
	if err != nil {
		return nil, fmt.Errorf("Unload: get data failed %w", err)
	}
	return d, nil
}

// Edit updates requested user's data in storage.
func (s *DataService) Edit(ctx context.Context, data *Data) error {
	err := s.repo.UpdateData(ctx, data)
	if err != nil {
		return fmt.Errorf("Edit: edit data failed %w", err)
	}
	return nil
}

// Delete deletes requested user's data from the storage.
func (s *DataService) Delete(ctx context.Context, name string) error {
	err := s.repo.DeleteDataByName(ctx, name)
	if err != nil {
		return fmt.Errorf("Delete: delete data failed %w", err)
	}
	return nil
}
