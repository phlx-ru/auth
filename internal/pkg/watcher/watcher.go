package watcher

import (
	"context"
	"errors"
	"fmt"
	"reflect"
	"runtime/debug"
	"strings"

	"github.com/go-kratos/kratos/v2/log"

	pkgStrings "auth/internal/pkg/strings"
)

type Watcher struct {
	logger              logger
	metrics             metrics
	method              string
	metricPrefix        string
	timing              Timing
	fields              map[string]any
	asWarning           bool
	ignoredErrors       []error
	ignoredErrorsChecks []func(error) bool
}

// New builds watcher, that can use for a struct variable
// Usage:
//
//	type SleepingStruct struct {
//		metrics metrics
//		logger  logger
//		watcher *Watcher
//	}
//
//	func NewSleepingStruct(metrics metrics, logger logger) *SleepingStruct {
//		return &SleepingStruct{
//			metrics: metrics,
//			logger:  logger,
//			watcher: New(`metric.prefix`, logger, metrics),
//		}
//	}
//
//	func (s *SleepingStruct) SleepForAWhile(ctx context.Context, duration time.Duration, errorMessage string) error {
//		var err error
//		defer s.watcher.OnPreparedMethod(`SleepForAWhile`).WithTimings().Results(func() (context.Context, error) {
//			return ctx, err
//		})
//		if errorMessage != "" {
//			err = fmt.Errorf(errorMessage)
//		}
//		time.Sleep(duration)
//		return err
//	}
func New(metricPrefix string, logger logger, metrics metrics) *Watcher {
	return &Watcher{
		logger:       logger,
		metrics:      metrics,
		metricPrefix: metricPrefix,
		fields:       map[string]any{},
	}
}

// Make makes fluent interface base with only metric prefix.
// Usage:
//
//	func (s *SleepingStruct) SleepForAWhile(ctx context.Context, duration time.Duration, errorMessage string) error {
//		var err error
//		defer Make(`services.prefix`).
//			OnPreparedMethod(`SleepForAWhile`).
//			WithLogger(s.logger).
//			WithMetrics(s.metrics).
//			WithTimings().
//			Results(func() (context.Context, error) { return ctx, err })
//		if errorMessage != "" {
//			err = fmt.Errorf(errorMessage)
//		}
//		time.Sleep(duration)
//		return err
//	}
func Make(metricPrefix string) *Watcher {
	return &Watcher{
		metricPrefix: metricPrefix,
		fields:       map[string]any{},
	}
}

func (w *Watcher) WithMetrics(metrics metrics) *Watcher {
	n := *w
	n.metrics = metrics
	return &n
}

func (w *Watcher) WithLogger(logger logger) *Watcher {
	n := *w
	n.logger = logger
	return &n
}

// WithTimings add timings. WARNING: if call it before WithMetrics() then empty timing will be used
func (w *Watcher) WithTimings() *Watcher {
	n := *w
	if !isNil(n.metrics) {
		n.timing = n.metrics.NewTiming()
	} else {
		n.timing = NewEmptyTiming()
	}
	return &n
}

func (w *Watcher) WithFields(fields map[string]any) *Watcher {
	n := *w
	n.fields = fields
	return &n
}

func (w *Watcher) WithIgnoredErrors(ignoredErrors []error) *Watcher {
	n := *w
	n.ignoredErrors = ignoredErrors
	return &n
}

func (w *Watcher) WithIgnoredErrorsChecks(ignoredErrorsChecks []func(error) bool) *Watcher {
	n := *w
	n.ignoredErrorsChecks = ignoredErrorsChecks
	return &n
}

func (w *Watcher) AsWarning() *Watcher {
	n := *w
	n.asWarning = true
	return &n
}

func (w *Watcher) AsError() *Watcher {
	n := *w
	n.asWarning = false
	return &n
}

func prepareMethodForMetric(method string) string {
	if method == "" {
		return ""
	}
	return strings.ToLower(method[:1]) + method[1:]
}

func (w *Watcher) OnPreparedMethod(method string) *Watcher {
	n := *w
	n.method = prepareMethodForMetric(method)
	return &n
}

func (w *Watcher) OnMethod(method string) *Watcher {
	n := *w
	n.method = method
	return &n
}

type ContextAndErrorCatcher func() (context.Context, error)

func (w *Watcher) Results(catcher ContextAndErrorCatcher) {
	ctx, err := catcher()
	result := `success`
	isIgnored := isIgnoredError(err, w.ignoredErrors, w.ignoredErrorsChecks)
	if err != nil && !isIgnored {
		result = `failure`
	}
	if !isNil(w.logger) {
		if w.method == "" {
			w.logger.WithContext(ctx).Errorf("empty 'method' on watcher for metric prefix '%s'", w.metricPrefix)
			w.method = "unknown"
		}
		w.fields[log.DefaultMessageKey] = fmt.Sprintf(`%s has %s on %s`, w.metricPrefix, result, w.method)
		if err != nil && !isIgnored {
			w.fields[`error`] = err
			w.fields[`stack`] = string(debug.Stack())
		}
		kvs := make([]any, 0, len(w.fields)*2)
		for field, value := range w.fields {
			kvs = append(kvs, field, value)
		}
		if err != nil && !isIgnored {
			if w.asWarning {
				w.logger.WithContext(ctx).Warnw(kvs...)
			} else {
				w.logger.WithContext(ctx).Errorw(kvs...)
			}
		} else {
			w.logger.WithContext(ctx).Infow(kvs...)
		}
	}
	metricStarts := w.method
	if w.metricPrefix != "" {
		metricStarts = pkgStrings.Metric(w.metricPrefix, w.method)
	}
	if !isNil(w.timing) {
		w.timing.Send(pkgStrings.Metric(metricStarts, `timings`, result))
	}
	if !isNil(w.metrics) {
		w.metrics.Increment(pkgStrings.Metric(metricStarts, result))
	}
}

func isNil(value any) bool {
	return value == nil || reflect.ValueOf(value).IsZero()
}

func isIgnoredError(err error, ignoredErrors []error, ignoredErrorsChecks []func(error) bool) bool {
	if err == nil {
		return true
	}
	for _, ignoredError := range ignoredErrors {
		if errors.Is(err, ignoredError) {
			return true
		}
	}
	for _, ignoredErrorCheck := range ignoredErrorsChecks {
		if ignoredErrorCheck(err) {
			return true
		}
	}
	return false
}
