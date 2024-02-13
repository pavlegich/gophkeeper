// Package main contains actions for building and running the server.
package main

import (
	"github.com/pavlegich/gophkeeper/internal/common/infra/logger"
	"github.com/pavlegich/gophkeeper/internal/server/app"
	"go.uber.org/zap"
)

func main() {
	if err := app.Run(); err != nil {
		logger.Log.Error("main: run app failed",
			zap.Error(err))
	}
	logger.Log.Info("quit")
}
