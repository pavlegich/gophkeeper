// Package http contains object of user handler,
// functions for activating the user handler in controller
// and user handlers.
package http

import (
	"database/sql"

	"github.com/go-chi/chi"
	"github.com/pavlegich/gophkeeper/internal/infra/config"
	"github.com/pavlegich/gophkeeper/internal/server/domains/user"
)

// UserHandler contatins objects for work
// with user handlers.
type UserHandler struct {
	config  *config.ServerConfig
	service user.Service
}

// Activate activates handler for user object.
func Activate(r *chi.Mux, cfg *config.ServerConfig, db *sql.DB) {
	// s := user.NewUserService(repo.NewUserRepo(db))
}

// newHandler initializes handler fpr
func newHandler(r *chi.Mux, cfg *config.ServerConfig, s user.Service) {
	// h := &UserHandler{
	// config:  cfg,
	// service: s,
	// }
	// r.Post("/api/user/register", h.Handle)
}
