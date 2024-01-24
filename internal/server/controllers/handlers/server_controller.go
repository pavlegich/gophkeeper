// Package handlers contains server controller object and
// methods for building the server
package handlers

import (
	"context"
	"database/sql"

	"github.com/go-chi/chi/v5"
	"github.com/pavlegich/gophkeeper/internal/infra/config"
	"github.com/pavlegich/gophkeeper/internal/server/controllers/middlewares"
)

// Controller contains database and configuration
// for building the server router.
type Controller struct {
	db  *sql.DB
	cfg *config.ServerConfig
}

// NewController returns new server controller.
func NewController(ctx context.Context, db *sql.DB, cfg *config.ServerConfig) *Controller {
	return &Controller{
		db:  db,
		cfg: cfg,
	}
}

// BuildRoute creates new router and appends handlers and middlewares to it.
func (c *Controller) BuildRoute(ctx context.Context) (*chi.Mux, error) {
	r := chi.NewRouter()

	r.Use(middlewares.WithLogging)

	return r, nil
}
