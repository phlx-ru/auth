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

func (u *userRepo) Create(ctx context.Context, user *ent.User) (*ent.User, error) {
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
		SetNillablePasswordResetExpiredAt(user.PasswordResetExpiredAt).
		SetNillableDeactivatedAt(user.DeactivatedAt).
		Save(ctx)
}

// Update all fields of user record. CAUTION: if field in 'user' not set — it will be cleared
func (u *userRepo) Update(ctx context.Context, user *ent.User) (*ent.User, error) {
	defer u.metric.NewTiming().Send(metricUserUpdateTimings)
	if user == nil {
		return nil, errors.New("user is empty")
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

	return updated.Save(ctx)
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
