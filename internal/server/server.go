// Package server contains Server object and its methods.
package server

import (
	"context"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/pavlegich/gophkeeper/internal/infra/config"
	"github.com/pavlegich/gophkeeper/internal/server/controllers/middlewares"
)

// Server contains server attributes.
type Server struct {
	server *http.Server
	config *config.ServerConfig
}

// NewServer initializes controller and router, returns new server object.
func NewServer(ctx context.Context, mh *chi.Mux, cfg *config.ServerConfig) (*Server, error) {
	r := chi.NewRouter()
	r.Use(middlewares.Recovery)
	r.Mount("/", mh)

	srv := &http.Server{
		Addr:    cfg.Address,
		Handler: r,
	}

	return &Server{
		server: srv,
		config: cfg,
	}, nil
}

// GetAddress returns server's address.
func (s *Server) GetAddress(ctx context.Context) string {
	return s.config.Address
}

// Serve start listening the network by the server.
func (s *Server) Serve(ctx context.Context) error {
	return s.server.ListenAndServe()
}

// Shutdown gracefully shuts down the server without interrupting any active connections.
func (s *Server) Shutdown(ctx context.Context) error {
	return s.server.Shutdown(ctx)
}
