package stubs

import (
	"github.com/go-kratos/kratos/v2/log"
	"gopkg.in/alexcesaro/statsd.v2"

	"auth/internal/pkg/logger"
	"auth/internal/pkg/metrics"
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

type LoggerMuted struct{}

func (l *LoggerMuted) Log(_ log.Level, _ ...any) error {
	return nil
}

func NewLoggerMuted() log.Logger {
	return &LoggerMuted{}
}
