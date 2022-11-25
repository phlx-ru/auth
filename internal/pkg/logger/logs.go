package logger

import (
	"context"
	"fmt"

	"github.com/getsentry/sentry-go"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware/tracing"
)

const (
	messageKey = "msg"
)

type Logger interface {
	WithContext(ctx context.Context) *log.Helper
	Log(level log.Level, keyvals ...interface{})
	Debug(a ...interface{})
	Debugf(format string, a ...interface{})
	Debugw(keyvals ...interface{})
	Info(a ...interface{})
	Infof(format string, a ...interface{})
	Infow(keyvals ...interface{})
	Warn(a ...interface{})
	Warnf(format string, a ...interface{})
	Warnw(keyvals ...interface{})
	Error(a ...interface{})
	Errorf(format string, a ...interface{})
	Errorw(keyvals ...interface{})
	Fatal(a ...interface{})
	Fatalf(format string, a ...interface{})
	Fatalw(keyvals ...interface{})
}

type Log struct {
	logger log.Logger
}

func (l *Log) Log(level log.Level, keyvals ...interface{}) error {
	err := l.logger.Log(level, keyvals...)
	if level == log.LevelWarn || level == log.LevelError {
		tags := ExtractMapFromKeyvals(keyvals...)
		msg, ok := tags[messageKey]
		if ok {
			delete(tags, messageKey)
			sentry.WithScope(func(scope *sentry.Scope) {
				scope.SetExtras(tags)
				if level == log.LevelWarn {
					sentry.CaptureMessage(fmt.Sprintf(`warning: %v`, msg))
				}
				if level == log.LevelError {
					sentry.CaptureException(fmt.Errorf("error: %v", msg))
				}
			})
		}
	}
	return err
}

func ExtractMapFromKeyvals(keyvals ...interface{}) map[string]any {
	res := map[string]any{}
	if len(keyvals) == 0 {
		return res
	}
	if len(keyvals) == 1 {
		res[messageKey] = keyvals[0]
		return res
	}
	length := len(keyvals)
	for i := 0; i < length; i += 2 {
		key := fmt.Sprintf("%v", keyvals[i])
		var value any
		if i != length-1 {
			value = keyvals[i+1]
		}
		res[key] = value
	}
	return res
}

func NewWithSentry() log.Logger {
	return &Log{
		logger: log.DefaultLogger,
	}
}

func New(id, name, version, level string) *log.Filter {
	loggerWithSentry := NewWithSentry()
	// loggerWithSentry := log.DefaultLogger
	loggerInstance := log.With(
		loggerWithSentry,
		"ts", log.DefaultTimestamp,
		"caller", log.DefaultCaller,
		"service.id", id,
		"service.name", name,
		"service.version", version,
		"trace.id", tracing.TraceID(),
		"span.id", tracing.SpanID(),
	)
	logLevel := log.ParseLevel(level)
	logger := log.NewFilter(loggerInstance, log.FilterLevel(logLevel))
	log.SetLogger(logger) // TODO Is it ok?
	return logger
}

func NewHelper(logger log.Logger, kv ...interface{}) *log.Helper {
	return log.NewHelper(log.With(logger, kv...))
}
