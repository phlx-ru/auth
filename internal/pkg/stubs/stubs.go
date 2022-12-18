package stubs

import (
	"github.com/go-kratos/kratos/v2/log"
	"gopkg.in/alexcesaro/statsd.v2"

	"auth/internal/pkg/logger"
	"auth/internal/pkg/metrics"
	"auth/internal/pkg/watcher"
)

//go:generate moq -out stubs_moq.go -stub . Metrics Logger LoggerBase

type Metrics interface {
	metrics.Metrics
}

type Logger interface {
	logger.Logger
}

type LoggerBase interface {
	log.Logger
}

func NewMetricsMuted() Metrics {
	client, _ := statsd.New(statsd.Mute(true))
	return client
}

type LoggerBaseMuted struct{}

func (l *LoggerBaseMuted) Log(_ log.Level, _ ...any) error {
	return nil
}

func NewLoggerBaseMuted() log.Logger {
	return &LoggerBaseMuted{}
}

func NewLoggerMuted() logger.Logger {
	return logger.NewHelper(NewLoggerBaseMuted(), `module`, `internal/pkg/stubs`)
}

func NewWatcherMuted(metricPrefix string) *watcher.Watcher {
	return watcher.New(metricPrefix, NewLoggerMuted(), NewMetricsMuted())
}
