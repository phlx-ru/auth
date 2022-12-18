package data

import (
	"context"
	databaseSql "database/sql"
	"errors"
	"time"

	"auth/ent"
	"auth/ent/predicate"
	"auth/internal/pkg/logger"
	"auth/internal/pkg/metrics"
	"auth/internal/pkg/watcher"

	"entgo.io/ent/dialect/sql"
	"github.com/go-kratos/kratos/v2/log"
)

const (
	metricPrefixHistory = `data.history`
)

type HistoryRepo struct {
	data    Database
	metric  metrics.Metrics
	logger  *log.Helper
	watcher *watcher.Watcher
}

func NewHistoryRepo(data Database, logs log.Logger, metric metrics.Metrics) *HistoryRepo {
	loggerHelper := logger.NewHelper(logs, `ts`, log.DefaultTimestamp, `scope`, metricPrefixHistory)
	return &HistoryRepo{
		data:    data,
		metric:  metric,
		logger:  loggerHelper,
		watcher: watcher.New(metricPrefixHistory, loggerHelper, metric),
	}
}

func (h *HistoryRepo) Create(ctx context.Context, history *ent.History) (*ent.History, error) {
	var err error
	if history == nil {
		err = errors.New("code is empty")
		return nil, err
	}
	defer h.watcher.OnPreparedMethod(`Create`).WithFields(map[string]any{
		"userId": history.UserID,
		"event":  history.Event,
	}).Results(func() (context.Context, error) {
		return ctx, err
	})

	history, err = h.client(ctx).Create().
		SetUserID(history.UserID).
		SetEvent(history.Event).
		SetNillableIP(history.IP).
		SetNillableUserAgent(history.UserAgent).
		Save(ctx)
	return history, err
}

func (h *HistoryRepo) FindLastUserEvents(
	ctx context.Context,
	userID int,
	types []string,
	interval time.Duration,
) ([]*ent.History, error) {
	var err error
	defer h.watcher.OnPreparedMethod(`FindLastUserEvents`).WithFields(map[string]any{
		"userId":   userID,
		"types":    types,
		"interval": interval,
	}).Results(func() (context.Context, error) {
		return ctx, err
	})

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

func (h *HistoryRepo) FindUserEvents(ctx context.Context, userID, limit, offset int) ([]*ent.History, error) {
	var err error
	defer h.watcher.OnPreparedMethod(`FindUserEvents`).WithFields(map[string]any{
		"userId": userID,
		"limit":  limit,
		"offset": offset,
	}).Results(func() (context.Context, error) {
		return ctx, err
	})

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

func (h *HistoryRepo) Transaction(
	ctx context.Context,
	txOptions *databaseSql.TxOptions,
	processes ...func(repoCtx context.Context) error,
) error {
	var err error
	defer h.watcher.OnPreparedMethod(`Transaction`).WithFields(map[string]any{
		"txOptions":       txOptions,
		"processesLength": len(processes),
	}).Results(func() (context.Context, error) {
		return ctx, err
	})

	err = transaction(h.data, h.logger)(ctx, txOptions, processes...)
	return err
}

func (h *HistoryRepo) client(ctx context.Context) *ent.HistoryClient {
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
