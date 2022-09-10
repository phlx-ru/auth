package data

import (
	"context"
	databaseSql "database/sql"
	"errors"
	"strings"
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
	metricCodeSaveTimings        = `data.code.save.timings`
	metricCodeFindForUserTimings = `data.code.findForUser.timings`
	metricCodeTransactionTimings = `data.code.transaction.timings`
)

type codeRepo struct {
	data   Database
	metric metrics.Metrics
	logs   *log.Helper
}

func NewCodeRepo(data Database, logs log.Logger, metric metrics.Metrics) biz.CodeRepo {
	return &codeRepo{
		data:   data,
		metric: metric,
		logs:   logger.NewHelper(logs, `ts`, log.DefaultTimestamp, `scope`, `data/code`),
	}
}

func (c *codeRepo) Save(ctx context.Context, code *ent.Code) (*ent.Code, error) {
	defer c.metric.NewTiming().Send(metricCodeSaveTimings)
	if code == nil {
		return nil, errors.New("code is empty")
	}

	return c.client(ctx).Create().
		SetUserID(code.UserID).
		SetContent(code.Content).
		SetExpiredAt(code.ExpiredAt).
		Save(ctx)
}

func (c *codeRepo) FindForUser(ctx context.Context, userID int) (*ent.Code, error) {
	defer c.metric.NewTiming().Send(metricCodeFindForUserTimings)
	if userID <= 0 {
		return nil, errors.New("userID must be greater than 0")
	}
	actualTime := time.Now()
	return c.client(ctx).Query().
		Where(codeFilterByUserID(userID)).
		Where(codeFilterNotExpired(actualTime)).
		Order(codeOrderByCreatedAt(`desc`)).
		First(ctx)
}

func (c *codeRepo) Transaction(
	ctx context.Context,
	txOptions *databaseSql.TxOptions,
	processes ...func(repoCtx context.Context) error,
) error {
	defer c.metric.NewTiming().Send(metricCodeTransactionTimings)
	return transaction(c.data, c.logs)(ctx, txOptions, processes...)
}

func (c *codeRepo) client(ctx context.Context) *ent.CodeClient {
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
	if strings.ToLower(direction) == `desc` {
		return ent.Desc(`created_at`)
	}
	return ent.Asc(`created_at`)
}
