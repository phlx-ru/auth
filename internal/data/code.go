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
	pkgStrings "auth/internal/pkg/strings"

	"entgo.io/ent/dialect/sql"
	"github.com/go-kratos/kratos/v2/log"
)

type codeRepo struct {
	data   Database
	metric metrics.Metrics
	logger *log.Helper
}

func NewCodeRepo(data Database, logs log.Logger, metric metrics.Metrics) biz.CodeRepo {
	return &codeRepo{
		data:   data,
		metric: metric,
		logger: logger.NewHelper(logs, `ts`, log.DefaultTimestamp, `scope`, `data/code`),
	}
}

func (c *codeRepo) postProcess(ctx context.Context, method string, err error) {
	if err != nil {
		c.logger.WithContext(ctx).Errorf(`history data method "%s" failed: %v`, method, err)
		c.metric.Increment(pkgStrings.Metric(metricPrefix, method, `failure`))
	} else {
		c.metric.Increment(pkgStrings.Metric(metricPrefix, method, `success`))
	}
}

func (c *codeRepo) Create(ctx context.Context, code *ent.Code) (*ent.Code, error) {
	method := `create` // nolint: goconst
	defer c.metric.NewTiming().Send(pkgStrings.Metric(metricPrefix, method, `timings`))
	var err error
	defer func() { c.postProcess(ctx, method, err) }()

	if code == nil {
		err = errors.New("code is empty")
		return nil, err
	}

	code, err = c.client(ctx).Create().
		SetUserID(code.UserID).
		SetContent(code.Content).
		SetExpiredAt(code.ExpiredAt).
		Save(ctx)
	return code, err
}

func (c *codeRepo) FindForUser(ctx context.Context, userID int) (*ent.Code, error) {
	method := `findForUser`
	defer c.metric.NewTiming().Send(pkgStrings.Metric(metricPrefix, method, `timings`))
	var err error
	defer func() { c.postProcess(ctx, method, err) }()

	var code *ent.Code
	actualTime := time.Now()
	code, err = c.client(ctx).Query().
		Where(codeFilterByUserID(userID)).
		Where(codeFilterNotExpired(actualTime)).
		Order(codeOrderByCreatedAt(orderDesc)).
		First(ctx)
	return code, err
}

func (c *codeRepo) Transaction(
	ctx context.Context,
	txOptions *databaseSql.TxOptions,
	processes ...func(repoCtx context.Context) error,
) error {
	method := `transaction` // nolint: goconst
	defer c.metric.NewTiming().Send(pkgStrings.Metric(metricPrefix, method, `timings`))
	var err error
	defer func() { c.postProcess(ctx, method, err) }()

	err = transaction(c.data, c.logger)(ctx, txOptions, processes...)
	return err
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
	if strings.ToLower(direction) == orderDesc {
		return ent.Desc(`created_at`)
	}
	return ent.Asc(`created_at`)
}
