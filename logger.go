package main

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func defaultLogger(ctx *gin.Context) *logrus.Logger {
	if ctx == nil {
		return newLogger()
	}

	var logger *logrus.Logger
	l, exists := ctx.Get("logger")
	if !exists {
		logger = newLogger()
		ctx.Set("logger", logger)
		return logger
	}

	logger, ok := l.(*logrus.Logger)
	if !ok {
		return newLogger()
	}

	return logger
}

func newLogger() *logrus.Logger {
	logger := logrus.New()
	logger.Formatter = &logrus.JSONFormatter{
		TimestampFormat: time.RFC3339Nano,
	}
	logger.Level = logrus.DebugLevel
	return logger
}