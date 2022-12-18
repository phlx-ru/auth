package data

import (
	"context"
	databaseSql "database/sql"
	"errors"
	"strings"
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
	metricPrefixCode = `data.code`
)

type CodeRepo struct {
	data    Database
	metric  metrics.Metrics
	logger  *log.Helper
	watcher *watcher.Watcher
}

func NewCodeRepo(data Database, logs log.Logger, metric metrics.Metrics) *CodeRepo {
	loggerHelper := logger.NewHelper(logs, `ts`, log.DefaultTimestamp, `scope`, metricPrefixCode)
	return &CodeRepo{
		data:    data,
		metric:  metric,
		logger:  loggerHelper,
		watcher: watcher.New(metricPrefixCode, loggerHelper, metric),
	}
}

func (c *CodeRepo) Create(ctx context.Context, code *ent.Code) (*ent.Code, error) {
	var err error
	if code == nil {
		err = errors.New("code is empty")
		return nil, err
	}
	defer c.watcher.OnPreparedMethod(`Create`).WithFields(map[string]any{
		"userId":    code.UserID,
		"expiredAt": code.ExpiredAt,
	}).Results(func() (context.Context, error) {
		return ctx, err
	})

	code, err = c.client(ctx).Create().
		SetUserID(code.UserID).
		SetContent(code.Content).
		SetExpiredAt(code.ExpiredAt).
		Save(ctx)
	return code, err
}

func (c *CodeRepo) FindForUser(ctx context.Context, userID int) (*ent.Code, error) {
	var err error
	defer c.watcher.OnPreparedMethod(`FindForUser`).WithFields(map[string]any{
		"userId": userID,
	}).Results(func() (context.Context, error) {
		return ctx, err
	})

	var code *ent.Code
	actualTime := time.Now()
	code, err = c.client(ctx).Query().
		Where(codeFilterByUserID(userID)).
		Where(codeFilterNotExpired(actualTime)).
		Order(codeOrderByCreatedAt(orderDesc)).
		First(ctx)
	return code, err
}

func (c *CodeRepo) Transaction(
	ctx context.Context,
	txOptions *databaseSql.TxOptions,
	processes ...func(repoCtx context.Context) error,
) error {
	var err error
	defer c.watcher.OnPreparedMethod(`Transaction`).WithFields(map[string]any{
		"txOptions":       txOptions,
		"processesLength": len(processes),
	}).Results(func() (context.Context, error) {
		return ctx, err
	})

	err = transaction(c.data, c.logger)(ctx, txOptions, processes...)
	return err
}

func (c *CodeRepo) client(ctx context.Context) *ent.CodeClient {
	return client(c.data)(ctx).Code
}

func codeFilterByUserID(userID int) predicate.Code {
	return func(selector *sql.Selector) {
		selector.Where(sql.P().EQ(`user_id`, userID))
	}
}

func codeFilterNotExpired(forTime time.Time) predicate.Code {
	return func(selector *sql.Selector) {
		selector.Where(sql.P().GT(`expired_at`, forTime))
	}
}

func codeOrderByCreatedAt(direction string) ent.OrderFunc {
	if strings.ToLower(direction) == orderDesc {
		return ent.Desc(`created_at`)
	}
	return ent.Asc(`created_at`)
}
