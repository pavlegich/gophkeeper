// Package http contains object of data handler,
// functions for activating the data handler in controller
// and data handlers.
package http

import (
	"database/sql"

	"github.com/go-chi/chi"
	"github.com/pavlegich/gophkeeper/internal/infra/config"
	"github.com/pavlegich/gophkeeper/internal/server/domains/data"
	repo "github.com/pavlegich/gophkeeper/internal/server/domains/data/repository"
)

// DataHandler contains objects for work
// with user handlers.
type DataHandler struct {
	Config  *config.ServerConfig
	Service data.Service
}

// Activate activates handler for user.
func Activate(r *chi.Mux, cfg *config.ServerConfig, db *sql.DB) {
	s := data.NewDataService(repo.NewDataRepository(db))
	newHandler(r, cfg, s)
}

// newHandler initializes handler for user.
func newHandler(r *chi.Mux, cfg *config.ServerConfig, s data.Service) {
	// h := &DataHandler{
	// 	Config:  cfg,
	// 	Service: s,
	// }
}
