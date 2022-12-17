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
	pkgStrings "auth/internal/pkg/strings"

	"entgo.io/ent/dialect/sql"
	"github.com/go-kratos/kratos/v2/log"
)

type SessionRepo struct {
	data   Database
	metric metrics.Metrics
	logger *log.Helper
}

func NewSessionRepo(data Database, logs log.Logger, metric metrics.Metrics) *SessionRepo {
	return &SessionRepo{
		data:   data,
		metric: metric,
		logger: logger.NewHelper(logs, `ts`, log.DefaultTimestamp, `scope`, `data/user`),
	}
}

func (s *SessionRepo) postProcess(ctx context.Context, method string, err error) {
	if err != nil {
		s.logger.WithContext(ctx).Errorf(`session data method "%s" failed: %v`, method, err)
		s.metric.Increment(pkgStrings.Metric(metricPrefix, method, `failure`))
	} else {
		s.metric.Increment(pkgStrings.Metric(metricPrefix, method, `success`))
	}
}

func (s *SessionRepo) Create(ctx context.Context, session *ent.Session) (*ent.Session, error) {
	method := `create` // nolint: goconst
	defer s.metric.NewTiming().Send(pkgStrings.Metric(metricPrefix, method, `timings`))
	var err error
	defer func() { s.postProcess(ctx, method, err) }()

	if session == nil {
		err = errors.New(`session is empty`)
		return nil, err
	}

	session, err = s.client(ctx).Create().
		SetUserID(session.UserID).
		SetToken(session.Token).
		SetIP(session.IP).
		SetUserAgent(session.UserAgent).
		SetNillableDeviceID(session.DeviceID).
		SetExpiredAt(session.ExpiredAt).
		SetIsActive(session.IsActive).
		Save(ctx)
	return session, err
}

func (s *SessionRepo) Update(ctx context.Context, session *ent.Session) (*ent.Session, error) {
	method := `update` // nolint: goconst
	defer s.metric.NewTiming().Send(pkgStrings.Metric(metricPrefix, method, `timings`))
	var err error
	defer func() { s.postProcess(ctx, method, err) }()

	if session == nil {
		err = errors.New(`session is empty`)
		return nil, err
	}

	updated := s.client(ctx).UpdateOne(session).
		SetUserID(session.UserID).
		SetToken(session.Token).
		SetIP(session.IP).
		SetUserAgent(session.UserAgent).
		SetExpiredAt(session.ExpiredAt).
		SetIsActive(session.IsActive)

	if session.DeviceID != nil {
		updated.SetDeviceID(*session.DeviceID)
	} else {
		updated.ClearDeviceID()
	}

	session, err = updated.Save(ctx)
	return session, err
}

func (s *SessionRepo) FindByUserID(ctx context.Context, userID int) (*ent.Session, error) {
	method := `findByUserID`
	defer s.metric.NewTiming().Send(pkgStrings.Metric(metricPrefix, method, `timings`))
	var err error
	defer func() { s.postProcess(ctx, method, err) }()

	var session *ent.Session
	actualTime := time.Now()
	session, err = s.client(ctx).Query().
		Where(sessionFilterByUserID(userID)).
		Where(sessionFilterNotExpired(actualTime)).
		Where(sessionFilterByIsActive(true)).
		Order(sessionOrderByCreatedAt(`desc`)).
		First(ctx)
	return session, err
}

func (s *SessionRepo) FindByToken(ctx context.Context, token string) (*ent.Session, error) {
	method := `findByToken`
	defer s.metric.NewTiming().Send(pkgStrings.Metric(metricPrefix, method, `timings`))
	var err error
	defer func() { s.postProcess(ctx, method, err) }()

	if token == "" {
		err = errors.New(`token is empty`)
		return nil, err
	}

	var session *ent.Session
	actualTime := time.Now()
	session, err = s.client(ctx).Query().
		Where(sessionFilterByToken(token)).
		Where(sessionFilterNotExpired(actualTime)).
		Where(sessionFilterByIsActive(true)).
		Order(sessionOrderByCreatedAt(`desc`)).
		First(ctx)
	return session, err
}

func (s *SessionRepo) Transaction(
	ctx context.Context,
	txOptions *databaseSql.TxOptions,
	processes ...func(repoCtx context.Context) error,
) error {
	method := `transaction`
	defer s.metric.NewTiming().Send(pkgStrings.Metric(metricPrefix, method, `timings`))
	var err error
	defer func() { s.postProcess(ctx, method, err) }()

	err = transaction(s.data, s.logger)(ctx, txOptions, processes...)
	return err
}

func (s *SessionRepo) client(ctx context.Context) *ent.SessionClient {
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
