package data

import (
	"context"
	databaseSql "database/sql"
	"errors"
	"time"

	"auth/ent"
	"auth/ent/predicate"
	"auth/internal/biz"
	"auth/internal/pkg/logger"
	"auth/internal/pkg/metrics"
	"auth/internal/pkg/strings"

	"entgo.io/ent/dialect/sql"
	"github.com/go-kratos/kratos/v2/log"
)

type historyRepo struct {
	data   Database
	metric metrics.Metrics
	logger *log.Helper
}

func NewHistoryRepo(data Database, logs log.Logger, metric metrics.Metrics) biz.HistoryRepo {
	return &historyRepo{
		data:   data,
		metric: metric,
		logger: logger.NewHelper(logs, `ts`, log.DefaultTimestamp, `scope`, `data/history`),
	}
}

func (h *historyRepo) postProcess(ctx context.Context, method string, err error) {
	if err != nil {
		h.logger.WithContext(ctx).Errorf(`history data method "%s" failed: %v`, method, err)
		h.metric.Increment(strings.Metric(metricPrefix, method, `failure`))
	} else {
		h.metric.Increment(strings.Metric(metricPrefix, method, `success`))
	}
}

func (h *historyRepo) Create(ctx context.Context, history *ent.History) (*ent.History, error) {
	method := `create`
	defer h.metric.NewTiming().Send(strings.Metric(metricPrefix, method, `timings`))
	var err error
	defer func() { h.postProcess(ctx, method, err) }()

	if history == nil {
		err = errors.New("code is empty")
		return nil, err
	}

	history, err = h.client(ctx).Create().
		SetUserID(history.UserID).
		SetEvent(history.Event).
		SetNillableIP(history.IP).
		SetNillableUserAgent(history.UserAgent).
		Save(ctx)
	return history, err
}

func (h *historyRepo) FindLastUserEvents(
	ctx context.Context,
	userID int,
	types []string,
	interval time.Duration,
) ([]*ent.History, error) {
	method := `findLastUserEvents`
	defer h.metric.NewTiming().Send(strings.Metric(metricPrefix, method, `timings`))
	var err error
	defer func() { h.postProcess(ctx, method, err) }()

	if userID <= 0 {
		err = errors.New("userID must be greater than 0")
	}
	if len(types) == 0 {
		err = errors.New("types is empty")
	}
	if interval <= 0 {
		err = errors.New("interval must be greater than 0")
	}
	if err != nil {
		return nil, err
	}

	var histories []*ent.History
	actualTime := time.Now()
	histories, err = h.client(ctx).Query().
		Where(historyFilterByUserID(userID)).
		Where(historyFilterByTypes(types)).
		Where(historyFilterByLastInterval(interval, actualTime)).
		All(ctx)
	return histories, err
}

func (h *historyRepo) FindUserEvents(ctx context.Context, userID, limit, offset int) ([]*ent.History, error) {
	method := `findUserEvents`
	defer h.metric.NewTiming().Send(strings.Metric(metricPrefix, method, `timings`))
	var err error
	defer func() { h.postProcess(ctx, method, err) }()

	if userID <= 0 {
		err = errors.New("userID must be greater than 0")
	}
	if limit <= 0 {
		err = errors.New("limit must be greater than 0")
	}
	if limit > 1000 {
		err = errors.New("limit must be less or equal than 1000")
	}
	if offset < 0 {
		err = errors.New("offset must be greater or equal than 0")
	}
	if offset > 10000 {
		err = errors.New("offset must be less or equal than 10000")
	}
	if err != nil {
		return nil, err
	}

	var histories []*ent.History
	histories, err = h.client(ctx).Query().
		Where(historyFilterByUserID(userID)).
		Offset(offset).
		Limit(limit).
		All(ctx)
	return histories, err
}

func (h *historyRepo) Transaction(
	ctx context.Context,
	txOptions *databaseSql.TxOptions,
	processes ...func(repoCtx context.Context) error,
) error {
	method := `transaction`
	defer h.metric.NewTiming().Send(strings.Metric(metricPrefix, method, `timings`))
	var err error
	defer func() { h.postProcess(ctx, method, err) }()

	err = transaction(h.data, h.logger)(ctx, txOptions, processes...)
	return err
}

func (h *historyRepo) client(ctx context.Context) *ent.HistoryClient {
	return client(h.data)(ctx).History
}

func historyFilterByUserID(userID int) predicate.History {
	return func(selector *sql.Selector) {
		selector.Where(sql.P().EQ(`user_id`, userID))
	}
}

func historyFilterByTypes(types []string) predicate.History {
	return func(selector *sql.Selector) {
		selector.Where(sql.P().In(`event`, itemsToAny(types)...))
	}
}

func historyFilterByLastInterval(interval time.Duration, now time.Time) predicate.History {
	return func(selector *sql.Selector) {
		selector.Where(sql.P().GTE(`created_at`, now.Add(-interval)))
	}
}

func itemsToAny[T comparable](items []T) []any {
	res := []any{}
	for _, item := range items {
		res = append(res, item)
	}
	return res
}
