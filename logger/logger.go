package logger

import "go.uber.org/zap"

var logger *zap.SugaredLogger

func init() {
	production, _ := zap.NewProduction()
	logger = production.Sugar()
	defer logger.Sync()
}

func Infof(template string, args ...interface{}) {
	logger.Infof(template, args)
}

func Warnf(template string, args ...interface{}) {
	logger.Warnf(template, args)
}

func Errorf(template string, args ...interface{}) {
	logger.Errorf(template, args)
}
