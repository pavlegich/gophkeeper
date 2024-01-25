// Package app contains the main methods for running the server.
package app

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/pavlegich/gophkeeper/internal/infra/config"
	"github.com/pavlegich/gophkeeper/internal/infra/database"
	"github.com/pavlegich/gophkeeper/internal/infra/hash"
	"github.com/pavlegich/gophkeeper/internal/infra/logger"
	"github.com/pavlegich/gophkeeper/internal/server"
	"github.com/pavlegich/gophkeeper/internal/server/controllers/handlers"
	_ "go.uber.org/automaxprocs"
	"go.uber.org/zap"
)

// Run initialized the main app components and runs the server.
func Run() error {
	// Context
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	defer stop()

	// Logger
	err := logger.Init(ctx, "Info")
	if err != nil {
		return fmt.Errorf("Run: logger initialization failed %w", err)
	}
	defer logger.Log.Sync()

	// Configuration
	cfg := config.NewServerConfig(ctx)
	err = cfg.ParseFlags(ctx)
	if err != nil {
		return fmt.Errorf("Run: parse flags failed %w", err)
	}

	// Database
	db, err := database.Init(ctx, cfg.DSN)
	if err != nil {
		return fmt.Errorf("Run: database initialization failed %w", err)
	}
	defer db.Close()

	// Generate keys for token
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return fmt.Errorf("Run: generate private key failed")
	}
	cfg.Token = hash.NewToken(privateKey, &privateKey.PublicKey, cfg.TokenExp)

	// Router
	ctrl := handlers.NewController(ctx, db, cfg)
	router, err := ctrl.BuildRoute(ctx)
	if err != nil {
		return fmt.Errorf("Run: build server route failed %w", err)
	}

	// Server
	srv, err := server.NewServer(ctx, router, cfg)
	if err != nil {
		return fmt.Errorf("Run: server initialization failed %w", err)
	}

	// Server graceful shutdown
	go func() {
		<-ctx.Done()
		ctxShutdown, cancelShutdown := context.WithTimeout(ctx, 5*time.Second)
		defer cancelShutdown()

		logger.Log.Info("shutting down gracefully...")

		err := srv.Shutdown(ctxShutdown)
		if err != nil {
			logger.Log.Error("server shutdown failed",
				zap.Error(err))
		}
	}()

	logger.Log.Info("running server", zap.String("addr", srv.GetAddress(ctx)))

	return srv.Serve(ctx)
}
