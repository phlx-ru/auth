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

const (
	metricPrefix = `data.user`
)

type userRepo struct {
	data   Database
	metric metrics.Metrics
	logger *log.Helper
}

func NewUserRepo(data Database, logs log.Logger, metric metrics.Metrics) biz.UserRepo {
	return &userRepo{
		data:   data,
		metric: metric,
		logger: logger.NewHelper(logs, `ts`, log.DefaultTimestamp, `scope`, `data/user`),
	}
}

func (u *userRepo) postProcess(ctx context.Context, method string, err error) {
	if err != nil {
		u.logger.WithContext(ctx).Errorf(`user data method "%s" failed: %v`, method, err)
		u.metric.Increment(strings.Metric(metricPrefix, method, `failure`))
	} else {
		u.metric.Increment(strings.Metric(metricPrefix, method, `success`))
	}
}

func (u *userRepo) Create(ctx context.Context, user *ent.User) (*ent.User, error) {
	method := `create`
	defer u.metric.NewTiming().Send(strings.Metric(metricPrefix, method, `timings`))
	var err error
	defer func() { u.postProcess(ctx, method, err) }()

	if user == nil {
		err = errors.New("user is empty")
		return nil, err
	}

	user, err = u.client(ctx).Create().
		SetDisplayName(user.DisplayName).
		SetType(user.Type).
		SetNillableEmail(user.Email).
		SetNillablePhone(user.Phone).
		SetNillableTelegramChatID(user.TelegramChatID).
		SetNillablePasswordHash(user.PasswordHash).
		SetNillablePasswordReset(user.PasswordReset).
		SetNillablePasswordResetExpiredAt(user.PasswordResetExpiredAt).
		SetNillableDeactivatedAt(user.DeactivatedAt).
		Save(ctx)

	return user, err
}

// Update all fields of user record. CAUTION: if field in 'user' not set — it will be cleared
func (u *userRepo) Update(ctx context.Context, user *ent.User) (*ent.User, error) {
	method := `update`
	defer u.metric.NewTiming().Send(strings.Metric(metricPrefix, method, `timings`))
	var err error
	defer func() { u.postProcess(ctx, method, err) }()

	if user == nil {
		err = errors.New("user is empty")
		return nil, err
	}

	updated := u.client(ctx).UpdateOne(user).
		SetDisplayName(user.DisplayName).
		SetType(user.Type)

	if user.Email != nil {
		updated.SetEmail(*user.Email)
	} else {
		updated.ClearEmail()
	}

	if user.Phone != nil {
		updated.SetPhone(*user.Phone)
	} else {
		updated.ClearEmail()
	}

	if user.TelegramChatID != nil {
		updated.SetTelegramChatID(*user.TelegramChatID)
	} else {
		updated.ClearTelegramChatID()
	}

	if user.DeactivatedAt != nil {
		updated.SetDeactivatedAt(*user.DeactivatedAt)
	} else {
		updated.ClearDeactivatedAt()
	}

	if user.PasswordHash != nil {
		updated.SetPasswordHash(*user.PasswordHash)
	} else {
		updated.ClearPasswordHash()
	}

	if user.PasswordReset != nil {
		updated.SetPasswordReset(*user.PasswordReset)
	} else {
		updated.ClearPasswordReset()
	}

	if user.PasswordResetExpiredAt != nil {
		updated.SetPasswordResetExpiredAt(*user.PasswordResetExpiredAt)
	} else {
		updated.ClearPasswordResetExpiredAt()
	}

	user, err = updated.Save(ctx)
	return user, err
}

func (u *userRepo) Activate(ctx context.Context, userID int) (*ent.User, error) {
	method := `activate`
	defer u.metric.NewTiming().Send(strings.Metric(metricPrefix, method, `timings`))
	var err error
	defer func() { u.postProcess(ctx, method, err) }()

	var user *ent.User
	user, err = u.client(ctx).UpdateOneID(userID).ClearDeactivatedAt().Save(ctx)
	return user, err
}

func (u *userRepo) Deactivate(ctx context.Context, userID int) (*ent.User, error) {
	method := `deactivate`
	defer u.metric.NewTiming().Send(strings.Metric(metricPrefix, method, `timings`))
	var err error
	defer func() { u.postProcess(ctx, method, err) }()

	var user *ent.User
	user, err = u.client(ctx).UpdateOneID(userID).SetDeactivatedAt(time.Now()).Save(ctx)
	return user, err
}

func (u *userRepo) FindByID(ctx context.Context, id int) (*ent.User, error) {
	method := `findByID`
	defer u.metric.NewTiming().Send(strings.Metric(metricPrefix, method, `timings`))
	var err error
	defer func() { u.postProcess(ctx, method, err) }()

	var user *ent.User
	user, err = u.client(ctx).Get(ctx, id)
	return user, err
}

func (u *userRepo) FindByEmail(ctx context.Context, email string) (*ent.User, error) {
	method := `findByEmail`
	defer u.metric.NewTiming().Send(strings.Metric(metricPrefix, method, `timings`))
	var err error
	defer func() { u.postProcess(ctx, method, err) }()

	var user *ent.User
	user, err = u.client(ctx).Query().
		Where(userFilterActive()).
		Where(userFilterByEmail(email)).
		Only(ctx)
	return user, err
}

func (u *userRepo) FindByPhone(ctx context.Context, phone string) (*ent.User, error) {
	method := `findByPhone`
	defer u.metric.NewTiming().Send(strings.Metric(metricPrefix, method, `timings`))
	var err error
	defer func() { u.postProcess(ctx, method, err) }()

	var user *ent.User
	user, err = u.client(ctx).Query().
		Where(userFilterActive()).
		Where(userFilterByPhone(phone)).
		Only(ctx)
	return user, err
}

func (u *userRepo) Transaction(
	ctx context.Context,
	txOptions *databaseSql.TxOptions,
	processes ...func(repoCtx context.Context) error,
) error {
	method := `transaction`
	defer u.metric.NewTiming().Send(strings.Metric(metricPrefix, method, `timings`))
	var err error
	defer func() { u.postProcess(ctx, method, err) }()

	err = transaction(u.data, u.logger)(ctx, txOptions, processes...)
	return err
}

func (u *userRepo) client(ctx context.Context) *ent.UserClient {
	return client(u.data)(ctx).User
}

func userFilterActive() predicate.User {
	return func(selector *sql.Selector) {
		selector.Where(sql.P().IsNull(`deactivated_at`))
	}
}

func userFilterByEmail(email string) predicate.User {
	return func(selector *sql.Selector) {
		selector.Where(sql.P().EQ(`email`, email))
	}
}

func userFilterByPhone(phone string) predicate.User {
	return func(selector *sql.Selector) {
		selector.Where(sql.P().EQ(`phone`, phone))
	}
}
