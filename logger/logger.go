package logger

import (
	"bounzr/iam/config"
	"go.uber.org/zap"
)

/**
var Logger *zap.Logger

func Init(){
	var cfg zap.Config
	cfg = config.IAM.Logger
	Logger, err := cfg.Build()
	if err != nil {
		panic(err)
	}
	defer Logger.Sync()

	Logger.Info("logger construction succeeded")
}
**/

func GetLogger() *zap.Logger {
	var cfg zap.Config
	cfg = config.IAM.Logger
	logger, err := cfg.Build()
	if err != nil {
		panic(err)
	}
	defer logger.Sync()
	return logger
}
