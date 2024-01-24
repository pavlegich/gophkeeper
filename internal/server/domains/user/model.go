// Package user contains objects and methods
// for user registration, login nad authorization.
package user

import "context"

// User contains information about user.
type User struct {
	ID       int    `db:"id" json:"id"`
	Login    string `db:"login" json:"login"`
	Password string `db:"password" json:"password"`
}

// Service describes methods related with user
// for communication between handlers and repositories.
type Service interface {
	Register(ctx context.Context, user *User) (*User, error)
	Login(ctx context.Context, user *User) (*User, error)
}

// Service describes methods related with user
// for communication between services and database.
type Repository interface {
	GetUserByLogin(ctx context.Context, login string) (*User, error)
	CreateUser(ctx context.Context, user *User) (*User, error)
}
