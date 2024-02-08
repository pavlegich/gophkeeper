// Package handlers contains server controller object and
// methods for building the server route.
package handlers

import (
	"context"
	"database/sql"

	"github.com/go-chi/chi/v5"
	"github.com/pavlegich/gophkeeper/internal/common/infra/config"
	"github.com/pavlegich/gophkeeper/internal/server/controllers/middlewares"
	data "github.com/pavlegich/gophkeeper/internal/server/domains/data/controllers/http"
	users "github.com/pavlegich/gophkeeper/internal/server/domains/user/controllers/http"
)

// Controller contains database and configuration
// for building the server router.
type Controller struct {
	db  *sql.DB
	cfg *config.ServerConfig
}

// NewController creates and returns new server controller.
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
	r.Use(middlewares.Recovery)
	r.Use(middlewares.WithAuth(c.cfg.Token))
	r.Use(middlewares.WithCompress)

	users.Activate(ctx, r, c.cfg, c.db)
	data.Activate(ctx, r, c.cfg, c.db)

	return r, nil
}
