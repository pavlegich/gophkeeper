// Package main contains actions for building and running the client.
package main

import (
	"github.com/pavlegich/gophkeeper/internal/client/app"
	"github.com/pavlegich/gophkeeper/internal/infra/logger"
	"go.uber.org/zap"
)

func main() {
	if err := app.Run(); err != nil {
		logger.Log.Error("main: run app failed",
			zap.Error(err))
	}
	logger.Log.Info("quit")
}
