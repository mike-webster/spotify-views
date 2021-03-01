package logging

import (
	"context"
	"os"
	"time"

	"github.com/mike-webster/spotify-views/keys"
	"github.com/sirupsen/logrus"
)

var _logger *logrus.Logger

func newLogger() *logrus.Logger {
	if _logger != nil {
		return _logger
	}
	_logger := logrus.New()
	if os.Getenv("GO_ENV") == "development" {
		_logger.Formatter = &logrus.TextFormatter{}
	} else {
		_logger.Formatter = &logrus.JSONFormatter{
			TimestampFormat: time.RFC3339Nano,
		}
	}
	_logger.Level = logrus.DebugLevel
	return _logger
}

// GetLogger will either return the logger from the context or
// it will return a new one with the default configuration.
func GetLogger(ctx context.Context) *logrus.Logger {
	if ctx == nil {
		return newLogger()
	}

	var logger *logrus.Logger
	l := keys.GetContextValue(ctx, keys.ContextLogger)
	logger, ok := l.(*logrus.Logger)
	if !ok {
		return newLogger()
	}

	return logger
}

// SetLogger will assign the given logger in the given context
func SetLogger(ctx context.Context, logger *logrus.Logger) {
	if ctx == nil {
		return
	}

	if logger == nil {
		ctx = context.WithValue(ctx, keys.ContextLogger, GetLogger(nil))
		return
	}

	ctx = context.WithValue(ctx, keys.ContextLogger, logger)
}
