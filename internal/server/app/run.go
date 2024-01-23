// Package app contains the main methods for running the server.
package app

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/pavlegich/gophkeeper/internal/infra/config"
	"github.com/pavlegich/gophkeeper/internal/infra/logger"
	"github.com/pavlegich/gophkeeper/internal/server"
	_ "go.uber.org/automaxprocs"
	"go.uber.org/zap"
)

// Run initialized the main app components and runs the server.
func Run(idleConnsClosed chan struct{}) error {
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

	// Server
	srv, err := server.NewServer(ctx, cfg)
	if err != nil {
		return fmt.Errorf("Run: server initialization failed %w", err)
	}

	// Server graceful shutdown
	go func() {
		<-ctx.Done()
		ctxShutdown, cancelShutdown := context.WithTimeout(ctx, 5*time.Second)
		defer cancelShutdown()

		err := srv.Shutdown(ctxShutdown)
		if err != nil {
			logger.Log.Error("server shutdown failed",
				zap.Error(err))
		}

		logger.Log.Info("shutting down gracefully...")
		close(idleConnsClosed)
	}()

	logger.Log.Info("running server", zap.String("addr", srv.GetAddress(ctx)))

	return srv.Serve(ctx)
}
