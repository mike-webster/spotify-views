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

// GetLogger will return the context logger for the request
func GetLogger(ctx context.Context) *logrus.Entry {
	iEnt := keys.GetContextValue(ctx, keys.ContextLogger)
	if iEnt != nil {
		entry, ok := iEnt.(*logrus.Entry)
		if ok {
			return entry
		}
	}

	return newLogger().WithField("entry_create", time.Now().UTC())
}

// SetRequestLogger will store the request log entry in the context
func SetRequestLogger(ctx context.Context, entry *logrus.Entry) {
	if ctx == nil || entry == nil {
		return
	}

	ctx = context.WithValue(ctx, keys.ContextLogger, entry)
}
