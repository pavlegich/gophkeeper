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
	"github.com/pavlegich/gophkeeper/internal/common/infra/config"
	"github.com/pavlegich/gophkeeper/internal/common/infra/logger"
	"go.uber.org/zap"
)

// go run -ldflags "-X 'main.buildVersion=v1.0.0' -X 'main.buildDate=$(date +'%d/%m/%Y')'" main.go
var buildVersion string = "N/A"
var buildDate string = "N/A"

func main() {
	// Context
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt,
		syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	defer stop()

	// Manager for read and write
	rw := rwmanager.NewRWManager(ctx, os.Stdin, os.Stdout)

	// Versions
	rw.Writeln(ctx, "Build version: "+buildVersion)
	rw.Writeln(ctx, "Build date: "+buildDate)

	// Greeting
	rw.Writeln(ctx, utils.Greet)
	// WaitGroup
	wg := &sync.WaitGroup{}

	// Logger
	err := logger.Init(ctx, "Panic")
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

	// Run client
	wg.Add(1)
	go func() {
		err := client.Serve(ctx)
		if err != nil {
			logger.Log.Error("main: client serve error", zap.Error(err))
		}
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
		case <-time.After(5 * time.Second):
			rw.Writeln(ctx, "\n"+utils.UnexpectedQuit)
			panic("shutdown timeout")
		}
	}

	rw.Writeln(ctx, utils.Quit)
}
