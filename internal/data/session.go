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
	metricPrefixSession = `data.session`
)

type SessionRepo struct {
	data    Database
	metric  metrics.Metrics
	logger  *log.Helper
	watcher *watcher.Watcher
}

func NewSessionRepo(data Database, logs log.Logger, metric metrics.Metrics) *SessionRepo {
	loggerHelper := logger.NewHelper(logs, `ts`, log.DefaultTimestamp, `scope`, metricPrefixSession)
	return &SessionRepo{
		data:    data,
		metric:  metric,
		logger:  loggerHelper,
		watcher: watcher.New(metricPrefixSession, loggerHelper, metric),
	}
}

func (s *SessionRepo) Create(ctx context.Context, session *ent.Session) (*ent.Session, error) {
	var err error
	if session == nil {
		err = errors.New(`session is empty`)
		return nil, err
	}
	defer s.watcher.OnPreparedMethod(`Create`).WithFields(map[string]any{
		"userId":    session.UserID,
		"expiredAt": session.ExpiredAt,
	}).Results(func() (context.Context, error) {
		return ctx, err
	})

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
	var err error
	if session == nil {
		err = errors.New(`session is empty`)
		return nil, err
	}
	defer s.watcher.OnPreparedMethod(`Update`).WithFields(map[string]any{
		"userId":    session.UserID,
		"expiredAt": session.ExpiredAt,
	}).Results(func() (context.Context, error) {
		return ctx, err
	})

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
	var err error
	defer s.watcher.OnPreparedMethod(`FindByUserID`).WithFields(map[string]any{
		"userId": userID,
	}).Results(func() (context.Context, error) {
		return ctx, err
	})

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
	var err error
	defer s.watcher.OnPreparedMethod(`FindByToken`).Results(func() (context.Context, error) {
		return ctx, err
	})

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
	var err error
	defer s.watcher.OnPreparedMethod(`Transaction`).WithFields(map[string]any{
		"txOptions":       txOptions,
		"processesLength": len(processes),
	}).Results(func() (context.Context, error) {
		return ctx, err
	})

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
