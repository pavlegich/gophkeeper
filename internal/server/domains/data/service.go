package data

import "context"

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

func (s *DataService) Create(ctx context.Context, data *Data) (*Data, error) {
	return &Data{}, nil
}
func (s *DataService) Unload(ctx context.Context, name string) (*Data, error) {
	return &Data{}, nil
}
func (s *DataService) Edit(ctx context.Context, data *Data) error {
	return nil
}
func (s *DataService) Delete(ctx context.Context, name string) error {
	return nil
}
