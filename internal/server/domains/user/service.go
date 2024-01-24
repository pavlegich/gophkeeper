package user

import (
	"context"
	"fmt"

	errs "github.com/pavlegich/gophkeeper/internal/server/errors"
	"golang.org/x/crypto/bcrypt"
)

// UserService contatins objects for user service.
type UserService struct {
	repo Repository
}

// NewUserService returns new user service.
func NewUserService(repo Repository) *UserService {
	return &UserService{
		repo: repo,
	}
}

// Register validates, stores and returns new user.
func (s *UserService) Register(ctx context.Context, user *User) (*User, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("Register: hash generate failed %w", err)
	}
	user.Password = string(hashedPassword)
	user, err = s.repo.CreateUser(ctx, user)
	if err != nil {
		return nil, fmt.Errorf("Register: save user failed %w", err)
	}
	return user, nil
}

// Login validates the obtained user credentials and returns stored user.
func (s *UserService) Login(ctx context.Context, user *User) (*User, error) {
	storedUser, err := s.repo.GetUserByLogin(ctx, user.Login)
	if err != nil {
		return nil, err
	}
	err = bcrypt.CompareHashAndPassword([]byte(storedUser.Password), []byte(user.Password))
	if err != nil {
		return nil, errs.ErrPasswordNotMatch
	}
	return storedUser, nil
}
