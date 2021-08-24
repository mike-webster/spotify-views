package logging

import (
	"context"
	"os"
	"time"

	"github.com/mike-webster/spotify-views/keys"
	"github.com/sirupsen/logrus"
)

type LoggerFields struct {
	UserAgent   string
	Referer     string
	QueryString string
	Path        string
	Method      string
	ClientIP    string
	RequestID   string
	UserID      string
}

var _logger *logrus.Logger

func newLogger() *logrus.Logger {
	if _logger != nil {
		return _logger
	}
	_logger := logrus.New()
	if os.Getenv("GO_ENV") != "production" {
		_logger.Formatter = &logrus.TextFormatter{}
	} else {
		_logger.Formatter = &logrus.JSONFormatter{
			TimestampFormat: time.RFC3339Nano,
		}
	}
	_logger.Level = logrus.DebugLevel

	return _logger
}

func parseRequestValues(logger *logrus.Logger, lf *LoggerFields) *logrus.Entry {
	entry := logrus.NewEntry(logger)
	entry = addNonEmptyField(lf.UserAgent, "user_agent", entry)
	entry = addNonEmptyField(lf.Referer, "referer", entry)
	entry = addNonEmptyField(lf.QueryString, "query_string", entry)
	entry = addNonEmptyField(lf.Path, "path", entry)
	entry = addNonEmptyField(lf.Method, "method", entry)
	entry = addNonEmptyField(lf.ClientIP, "client_ip", entry)
	entry = addNonEmptyField(lf.RequestID, "request_id", entry)
	entry = addNonEmptyField(lf.UserID, "user_id", entry)
	return entry
}

func addNonEmptyField(val string, name string, entry *logrus.Entry) *logrus.Entry {
	if len(val) > 0 {
		return entry.WithField(name, val)
	}

	return entry
}

// GetLogger will return the context logger for the request
func GetLogger(ctx context.Context) *logrus.Entry {
	if ctx != nil {
		iEnt := keys.GetContextValue(ctx, keys.ContextLogger)
		if iEnt != nil {
			entry, ok := iEnt.(*logrus.Entry)
			if ok {
				return entry
			}
		}
	}

	logger := newLogger()

	ilf := keys.GetContextValue(ctx, keys.ContextLoggerFields)
	if ilf != nil {
		lf, ok := ilf.(*LoggerFields)
		if ok {
			return parseRequestValues(logger, lf)
		}
	}

	return logger.WithField("cached_logger", "false")
}

// SetRequestLogger will store the request log entry in the context
func SetRequestLogger(ctx context.Context, entry *logrus.Entry) {
	if ctx == nil || entry == nil {
		return
	}

	ctx = context.WithValue(ctx, keys.ContextLogger, entry)
}
