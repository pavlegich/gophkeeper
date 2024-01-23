// Package main contains actions for building and running the server.
package main

import (
	"github.com/pavlegich/gophkeeper/internal/infra/logger"
	"github.com/pavlegich/gophkeeper/internal/server/app"
	"go.uber.org/zap"
)

func main() {
	idleConnsClosed := make(chan struct{})

	if err := app.Run(idleConnsClosed); err != nil {
		logger.Log.Error("main: run app failed",
			zap.Error(err))
	}

	<-idleConnsClosed
	logger.Log.Info("quit")
}
