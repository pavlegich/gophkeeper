package repository

import (
	"context"
	"database/sql"

	"github.com/pavlegich/gophkeeper/internal/server/domains/user"
)

type Repository struct {
	db *sql.DB
}

func NewUserRepo(db *sql.DB) *Repository {
	return &Repository{
		db: db,
	}
}

// GetUserByLogin возвращает конкретного пользователя из хранилища
func (r *Repository) GetUserByLogin(ctx context.Context, login string) (*user.User, error) {
	return &user.User{}, nil
}

// CreateUser сохраняет данные пользователя в хранилище
func (r *Repository) CreateUser(ctx context.Context, u *user.User) (*user.User, error) {
	return &user.User{}, nil
}
