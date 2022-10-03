package util

import "go.uber.org/zap"

func SetupLogger() {
	logger, _ := zap.NewProduction()
	zap.ReplaceGlobals(logger)
}

func SetupDevLogger() {
	logger, _ := zap.NewDevelopment()
	zap.ReplaceGlobals(logger)
}
