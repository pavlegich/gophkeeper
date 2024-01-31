// Package main contains actions for building and running the client.
package main

import (
	"context"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/pavlegich/gophkeeper/internal/client"
	"github.com/pavlegich/gophkeeper/internal/client/controllers"
	"github.com/pavlegich/gophkeeper/internal/client/domains/rwmanager"
	"github.com/pavlegich/gophkeeper/internal/client/utils"
	"github.com/pavlegich/gophkeeper/internal/infra/config"
	"github.com/pavlegich/gophkeeper/internal/infra/logger"
	"go.uber.org/zap"
)

// go run -ldflags "-X main.buildVersion=v1.0.1 -X main.buildDate=$(date +'%Y/%m/%d') -X main.buildCommit=1d1wdd1f" main.go
var buildVersion string = "v1.0.0"
var buildDate string = "N/A"
var buildCommit string = "N/A"

func main() {
	// Context
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt,
		syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	defer stop()

	// Manager for read and write
	// f, _ := os.Open("commands")
	// rw := rwmanager.NewRWManager(ctx, f, os.Stdout)
	rw := rwmanager.NewRWManager(ctx, os.Stdin, os.Stdout)

	// Versions
	rw.Write(ctx, "Build version: "+buildVersion)
	rw.WriteString(ctx, "Build date: "+buildDate)
	rw.WriteString(ctx, "Build commit: "+buildCommit)

	// Greeting
	rw.WriteString(ctx, utils.Greet)

	// WaitGroup
	wg := &sync.WaitGroup{}

	// Logger
	err := logger.Init(ctx, "Info")
	if err != nil {
		logger.Log.Error("main: logger initialization failed", zap.Error(err))
	}
	defer logger.Log.Sync()

	// Configuration
	cfg := config.NewClientConfig(ctx)
	err = cfg.ParseFlags(ctx)
	if err != nil {
		logger.Log.Error("main: parse flags failed", zap.Error(err))
	}

	// Client
	ctrl := controllers.NewController(ctx, rw, cfg)
	client, err := client.NewClient(ctx, ctrl, rw, cfg)
	if err != nil {
		logger.Log.Error("main: create new client failed", zap.Error(err))
	}

	wg.Add(1)
	go func() {
		client.Serve(ctx)
		stop()
		wg.Done()
	}()

	// Client graceful shutdown
	<-ctx.Done()
	if ctx.Err() != nil {
		logger.Log.Info("shutting down gracefully...",
			zap.Error(ctx.Err()))

		connsClosed := make(chan struct{})
		go func() {
			wg.Wait()
			close(connsClosed)
		}()

		select {
		case <-connsClosed:
		case <-time.After(15 * time.Second):
			panic("shutdown timeout")
		}
	}

	rw.Write(ctx, utils.Quit)
}
