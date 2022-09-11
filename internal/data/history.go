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

	"entgo.io/ent/dialect/sql"
	"github.com/go-kratos/kratos/v2/log"
)

const (
	metricHistorySaveTimings               = `data.history.save.timings`
	metricHistoryFindLastUserEventsTimings = `data.history.findLastUserEvents.timings`
	metricHistoryFindUserEventsTimings     = `data.history.findUserEvents.timings`
	metricHistoryTransactionTimings        = `data.history.transaction.timings`
)

type historyRepo struct {
	data   Database
	metric metrics.Metrics
	logs   *log.Helper
}

func NewHistoryRepo(data Database, logs log.Logger, metric metrics.Metrics) biz.HistoryRepo {
	return &historyRepo{
		data:   data,
		metric: metric,
		logs:   logger.NewHelper(logs, `ts`, log.DefaultTimestamp, `scope`, `data/history`),
	}
}

func (h *historyRepo) Create(ctx context.Context, history *ent.History) (*ent.History, error) {
	defer h.metric.NewTiming().Send(metricHistorySaveTimings)
	if history == nil {
		return nil, errors.New("code is empty")
	}

	return h.client(ctx).Create().
		SetUserID(history.UserID).
		SetEvent(history.Event).
		SetNillableIP(history.IP).
		SetNillableUserAgent(history.UserAgent).
		Save(ctx)
}

func (h *historyRepo) FindLastUserEvents(
	ctx context.Context,
	userID int,
	types []string,
	interval time.Duration,
) ([]*ent.History, error) {
	defer h.metric.NewTiming().Send(metricHistoryFindLastUserEventsTimings)
	if userID <= 0 {
		return nil, errors.New("userID must be greater than 0")
	}
	if len(types) == 0 {
		return nil, errors.New("types is empty")
	}
	if interval <= 0 {
		return nil, errors.New("interval must be greater than 0")
	}

	actualTime := time.Now()
	return h.client(ctx).Query().
		Where(historyFilterByUserID(userID)).
		Where(historyFilterByTypes(types)).
		Where(historyFilterByLastInterval(interval, actualTime)).
		All(ctx)
}

func (h *historyRepo) FindUserEvents(ctx context.Context, userID, limit, offset int) ([]*ent.History, error) {
	defer h.metric.NewTiming().Send(metricHistoryFindUserEventsTimings)
	if userID <= 0 {
		return nil, errors.New("userID must be greater than 0")
	}
	if limit <= 0 {
		return nil, errors.New("limit must be greater than 0")
	}
	if limit > 1000 {
		return nil, errors.New("limit must be less or equal than 1000")
	}
	if offset < 0 {
		return nil, errors.New("offset must be greater or equal than 0")
	}
	if offset > 10000 {
		return nil, errors.New("offset must be less or equal than 10000")
	}
	return h.client(ctx).Query().
		Where(historyFilterByUserID(userID)).
		Offset(offset).
		Limit(limit).
		All(ctx)
}

func (h *historyRepo) Transaction(
	ctx context.Context,
	txOptions *databaseSql.TxOptions,
	processes ...func(repoCtx context.Context) error,
) error {
	defer h.metric.NewTiming().Send(metricHistoryTransactionTimings)
	return transaction(h.data, h.logs)(ctx, txOptions, processes...)
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
