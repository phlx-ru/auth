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
	metricSessionSaveTimings          = `data.session.save.timings`
	metricSessionUpdateTimings        = `data.session.update.timings`
	metricSessionFindByUserIDTimings  = `data.session.findByUserId.timings`
	metricSessionFindByTokenIDTimings = `data.session.findByToken.timings`
)

type sessionRepo struct {
	data   Database
	metric metrics.Metrics
	logs   *log.Helper
}

func NewSessionRepo(data Database, logs log.Logger, metric metrics.Metrics) biz.SessionRepo {
	return &sessionRepo{
		data:   data,
		metric: metric,
		logs:   logger.NewHelper(logs, `ts`, log.DefaultTimestamp, `scope`, `data/user`),
	}
}

func (s *sessionRepo) Create(ctx context.Context, session *ent.Session) (*ent.Session, error) {
	defer s.metric.NewTiming().Send(metricSessionSaveTimings)
	if session == nil {
		return nil, errors.New(`session is empty`)
	}

	return s.client(ctx).Create().
		SetUserID(session.UserID).
		SetToken(session.Token).
		SetIP(session.IP).
		SetUserAgent(session.UserAgent).
		SetNillableDeviceID(session.DeviceID).
		SetExpiredAt(session.ExpiredAt).
		SetIsActive(session.IsActive).
		Save(ctx)
}

func (s *sessionRepo) Update(ctx context.Context, session *ent.Session) (*ent.Session, error) {
	defer s.metric.NewTiming().Send(metricSessionUpdateTimings)
	if session == nil {
		return nil, errors.New(`session is empty`)
	}

	return s.client(ctx).UpdateOne(session).
		SetUserID(session.UserID).
		SetToken(session.Token).
		SetIP(session.IP).
		SetUserAgent(session.UserAgent).
		SetNillableDeviceID(session.DeviceID).
		SetExpiredAt(session.ExpiredAt).
		SetIsActive(session.IsActive).
		Save(ctx)
}

func (s *sessionRepo) FindByUserID(ctx context.Context, userID int) (*ent.Session, error) {
	defer s.metric.NewTiming().Send(metricSessionFindByUserIDTimings)
	if userID <= 0 {
		return nil, errors.New(`userID is incorrect`)
	}

	actualTime := time.Now()
	return s.client(ctx).Query().
		Where(sessionFilterByUserID(userID)).
		Where(sessionFilterNotExpired(actualTime)).
		Where(sessionFilterByIsActive(true)).
		Order(sessionOrderByCreatedAt(`desc`)).
		First(ctx)
}

func (s *sessionRepo) FindByToken(ctx context.Context, token string) (*ent.Session, error) {
	defer s.metric.NewTiming().Send(metricSessionFindByTokenIDTimings)
	if token == "" {
		return nil, errors.New(`token is empty`)
	}

	actualTime := time.Now()
	return s.client(ctx).Query().
		Where(sessionFilterByToken(token)).
		Where(sessionFilterNotExpired(actualTime)).
		Where(sessionFilterByIsActive(true)).
		Order(sessionOrderByCreatedAt(`desc`)).
		First(ctx)
}

func (s *sessionRepo) Transaction(
	ctx context.Context,
	txOptions *databaseSql.TxOptions,
	processes ...func(repoCtx context.Context) error,
) error {
	defer s.metric.NewTiming().Send(metricUserTransactionTimings)
	return transaction(s.data, s.logs)(ctx, txOptions, processes...)
}

func (s *sessionRepo) client(ctx context.Context) *ent.SessionClient {
	return client(s.data)(ctx).Session
}

func sessionFilterByUserID(userID int) predicate.Session {
	return func(selector *sql.Selector) {
		selector.Where(sql.P().EQ(`user_id`, userID))
	}
}

func sessionFilterByToken(token string) predicate.Session {
	return func(selector *sql.Selector) {
		selector.Where(sql.P().EQ(`token`, token))
	}
}

func sessionFilterNotExpired(forTime time.Time) predicate.Session {
	return func(selector *sql.Selector) {
		selector.Where(sql.P().GT(`expired_at`, forTime))
	}
}

func sessionFilterByIsActive(isActive bool) predicate.Session {
	return func(selector *sql.Selector) {
		selector.Where(sql.P().EQ(`is_active`, isActive))
	}
}

func sessionOrderByCreatedAt(direction string) ent.OrderFunc {
	if strings.ToLower(direction) == `desc` {
		return ent.Desc(`created_at`)
	}
	return ent.Asc(`created_at`)
}
