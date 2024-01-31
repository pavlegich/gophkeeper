// Package user contains objects and methods
// for user registration, login and authorization.
package user

import "context"

// User contains information about user.
type User struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

// Service describes methods related with user.
type Service interface {
	Register(ctx context.Context) error
	Login(ctx context.Context) error
}
