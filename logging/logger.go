package logging

import (
	"context"
	"time"

	"github.com/sirupsen/logrus"
)

var _logger *logrus.Logger

type ContextKey string

var ContextLogger = ContextKey("logger")

func newLogger() *logrus.Logger {
	if _logger != nil {
		return _logger
	}
	_logger := logrus.New()
	_logger.Formatter = &logrus.JSONFormatter{
		TimestampFormat: time.RFC3339Nano,
	}
	_logger.Level = logrus.DebugLevel
	return _logger
}

// GetLogger will either return the logger from the context or
// it will return a new one with the default configuration.
func GetLogger(ctx *context.Context) *logrus.Logger {
	if ctx == nil {
		return newLogger()
	}

	var logger *logrus.Logger
	localCtx := *ctx
	l := localCtx.Value(ContextLogger)
	logger, ok := l.(*logrus.Logger)
	if !ok {
		l = localCtx.Value(string(ContextLogger))
		logger, ok := l.(*logrus.Logger)
		if !ok {
			return newLogger()
		}

		return logger
	}

	return logger
}
