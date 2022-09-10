package data

import (
	"context"
	databaseSql "database/sql"
	"errors"

	"auth/ent"
	"auth/ent/predicate"
	"auth/internal/biz"
	"auth/internal/pkg/logger"
	"auth/internal/pkg/metrics"
	"entgo.io/ent/dialect/sql"

	"github.com/go-kratos/kratos/v2/log"
)

const (
	metricUserSaveTimings        = `data.user.save.timings`
	metricUserUpdateTimings      = `data.user.update.timings`
	metricUserFindByIDTimings    = `data.user.findById.timings`
	metricUserFindByEmailTimings = `data.user.findByEmail.timings`
	metricUserFindByPhoneTimings = `data.user.findByPhone.timings`
	metricUserTransactionTimings = `data.user.transaction.timings`
)

type userRepo struct {
	data   Database
	metric metrics.Metrics
	logs   *log.Helper
}

func NewUserRepo(data Database, logs log.Logger, metric metrics.Metrics) biz.UserRepo {
	return &userRepo{
		data:   data,
		metric: metric,
		logs:   logger.NewHelper(logs, `ts`, log.DefaultTimestamp, `scope`, `data/user`),
	}
}

func (u *userRepo) Save(ctx context.Context, user *ent.User) (*ent.User, error) {
	defer u.metric.NewTiming().Send(metricUserSaveTimings)
	if user == nil {
		return nil, errors.New("user is empty")
	}

	return u.client(ctx).Create().
		SetDisplayName(user.DisplayName).
		SetType(user.Type).
		SetNillableEmail(user.Email).
		SetNillablePhone(user.Phone).
		SetNillableTelegramChatID(user.TelegramChatID).
		SetNillablePasswordHash(user.PasswordHash).
		SetNillablePasswordReset(user.PasswordReset).
		SetNillableDeactivatedAt(user.DeactivatedAt).
		Save(ctx)
}

func (u *userRepo) Update(ctx context.Context, user *ent.User) (*ent.User, error) {
	defer u.metric.NewTiming().Send(metricUserUpdateTimings)
	if user == nil {
		return nil, errors.New("user is empty")
	}

	return u.client(ctx).UpdateOne(user).
		SetDisplayName(user.DisplayName).
		SetType(user.Type).
		SetNillableEmail(user.Email).
		SetNillablePhone(user.Phone).
		SetNillableTelegramChatID(user.TelegramChatID).
		SetNillablePasswordHash(user.PasswordHash).
		SetNillablePasswordReset(user.PasswordReset).
		SetNillableDeactivatedAt(user.DeactivatedAt).
		Save(ctx)
}

func (u *userRepo) FindByID(ctx context.Context, id int) (*ent.User, error) {
	defer u.metric.NewTiming().Send(metricUserFindByIDTimings)
	return u.client(ctx).Get(ctx, id)
}

func (u *userRepo) FindByEmail(ctx context.Context, email string) (*ent.User, error) {
	defer u.metric.NewTiming().Send(metricUserFindByEmailTimings)
	return u.client(ctx).Query().
		Where(userFilterActive()).
		Where(userFilterByEmail(email)).
		Only(ctx)
}

func (u *userRepo) FindByPhone(ctx context.Context, phone string) (*ent.User, error) {
	defer u.metric.NewTiming().Send(metricUserFindByPhoneTimings)
	return u.client(ctx).Query().
		Where(userFilterActive()).
		Where(userFilterByPhone(phone)).
		Only(ctx)
}

func (u *userRepo) Transaction(
	ctx context.Context,
	txOptions *databaseSql.TxOptions,
	processes ...func(repoCtx context.Context) error,
) error {
	defer u.metric.NewTiming().Send(metricUserTransactionTimings)
	return transaction(u.data, u.logs)(ctx, txOptions, processes...)
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
