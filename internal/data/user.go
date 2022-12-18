package data

import (
	"context"
	databaseSql "database/sql"
	"errors"
	"fmt"
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
	metricPrefixUser = `data.user`
)

var (
	UserDefaultListOrderFields    = []string{`updated_at`}
	UserDefaultListOrderDirection = orderDesc
	UserAllowedListOrderFields    = []string{`id`, `type`, `created_at`, `updated_at`}
)

type UserRepo struct {
	data    Database
	metric  metrics.Metrics
	logger  *log.Helper
	watcher *watcher.Watcher
}

func NewUserRepo(data Database, logs log.Logger, metric metrics.Metrics) *UserRepo {
	loggerHelper := logger.NewHelper(logs, `ts`, log.DefaultTimestamp, `scope`, metricPrefixUser)
	return &UserRepo{
		data:    data,
		metric:  metric,
		logger:  loggerHelper,
		watcher: watcher.New(metricPrefixUser, loggerHelper, metric),
	}
}

func (u *UserRepo) List(ctx context.Context, limit, offset int64, orderFields []string, orderDirection string) (
	[]*ent.User,
	error,
) {
	var err error
	defer u.watcher.OnPreparedMethod(`List`).WithFields(map[string]any{
		"limit":          limit,
		"offset":         offset,
		"orderFields":    orderFields,
		"orderDirection": orderDirection,
	}).Results(func() (context.Context, error) {
		return ctx, err
	})

	if limit <= 0 {
		err = errors.New("limit must be greater than 0")
		return nil, err
	}

	maxLimit := int64(1000)
	if limit > maxLimit {
		err = fmt.Errorf("limit must be less than or equal %d", maxLimit)
		return nil, err
	}

	if offset < 0 {
		err = errors.New("offset must be greater than or equal 0")
		return nil, err
	}

	maxOffset := int64(1 << 16)
	if offset > maxOffset {
		err = fmt.Errorf("offset must be less than or equal %d", maxOffset)
		return nil, err
	}

	if orderDirection != orderAsc && orderDirection != orderDesc {
		err = fmt.Errorf("order direction must by '%s' or '%s'", orderAsc, orderDesc)
		return nil, err
	}

	for _, orderField := range orderFields {
		allowed := false
		for _, allowedField := range UserAllowedListOrderFields {
			if orderField == allowedField {
				allowed = true
				continue
			}
		}
		if !allowed {
			err = fmt.Errorf("order field '%s' is not allowed", orderField)
			return nil, err
		}
	}

	list, err := u.client(ctx).Query().
		Where(userFilterActive()).
		Limit(int(limit)).
		Offset(int(offset)).
		Order(userOrderByFields(orderFields, orderDirection)).
		All(ctx)

	return list, err
}

func (u *UserRepo) Create(ctx context.Context, user *ent.User) (*ent.User, error) {
	var err error
	if user == nil {
		err = errors.New("user is empty")
		return nil, err
	}
	defer u.watcher.OnPreparedMethod(`Create`).WithFields(map[string]any{
		"displayName": user.DisplayName,
	}).Results(func() (context.Context, error) {
		return ctx, err
	})

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
func (u *UserRepo) Update(ctx context.Context, user *ent.User) (*ent.User, error) {
	var err error
	if user == nil {
		err = errors.New("user is empty")
		return nil, err
	}
	defer u.watcher.OnPreparedMethod(`Update`).WithFields(map[string]any{
		"id": user.ID,
	}).Results(func() (context.Context, error) {
		return ctx, err
	})

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

func (u *UserRepo) Activate(ctx context.Context, userID int) (*ent.User, error) {
	var err error
	defer u.watcher.OnPreparedMethod(`Activate`).WithFields(map[string]any{
		"userId": userID,
	}).Results(func() (context.Context, error) {
		return ctx, err
	})

	var user *ent.User
	user, err = u.client(ctx).UpdateOneID(userID).ClearDeactivatedAt().Save(ctx)
	return user, err
}

func (u *UserRepo) Deactivate(ctx context.Context, userID int) (*ent.User, error) {
	var err error
	defer u.watcher.OnPreparedMethod(`Deactivate`).WithFields(map[string]any{
		"userId": userID,
	}).Results(func() (context.Context, error) {
		return ctx, err
	})

	var user *ent.User
	user, err = u.client(ctx).UpdateOneID(userID).SetDeactivatedAt(time.Now()).Save(ctx)
	return user, err
}

func (u *UserRepo) FindByID(ctx context.Context, userID int) (*ent.User, error) {
	var err error
	defer u.watcher.OnPreparedMethod(`FindByID`).WithFields(map[string]any{
		"userId": userID,
	}).Results(func() (context.Context, error) {
		return ctx, err
	})

	var user *ent.User
	user, err = u.client(ctx).Get(ctx, userID)
	return user, err
}

func (u *UserRepo) FindByEmail(ctx context.Context, email string) (*ent.User, error) {
	var err error
	defer u.watcher.OnPreparedMethod(`FindByEmail`).WithFields(map[string]any{
		"email": email,
	}).Results(func() (context.Context, error) {
		return ctx, err
	})

	var user *ent.User
	user, err = u.client(ctx).Query().
		Where(userFilterActive()).
		Where(userFilterByEmail(email)).
		Only(ctx)
	return user, err
}

func (u *UserRepo) FindByPhone(ctx context.Context, phone string) (*ent.User, error) {
	var err error
	defer u.watcher.OnPreparedMethod(`FindByPhone`).WithFields(map[string]any{
		"phone": phone,
	}).Results(func() (context.Context, error) {
		return ctx, err
	})

	var user *ent.User
	user, err = u.client(ctx).Query().
		Where(userFilterActive()).
		Where(userFilterByPhone(phone)).
		Only(ctx)
	return user, err
}

func (u *UserRepo) Transaction(
	ctx context.Context,
	txOptions *databaseSql.TxOptions,
	processes ...func(repoCtx context.Context) error,
) error {
	var err error
	defer u.watcher.OnPreparedMethod(`Transaction`).WithFields(map[string]any{
		"txOptions":       txOptions,
		"processesLength": len(processes),
	}).Results(func() (context.Context, error) {
		return ctx, err
	})

	err = transaction(u.data, u.logger)(ctx, txOptions, processes...)
	return err
}

func (u *UserRepo) client(ctx context.Context) *ent.UserClient {
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

func userOrderByFields(fields []string, direction string) ent.OrderFunc {
	if strings.ToLower(direction) == orderDesc {
		return ent.Desc(fields...)
	}
	return ent.Asc(fields...)
}
