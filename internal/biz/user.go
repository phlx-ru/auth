package biz

import (
	"context"

	"auth/ent"
	"auth/internal/pkg/logger"
	"auth/internal/pkg/metrics"

	"github.com/AlekSi/pointer"
	"github.com/go-kratos/kratos/v2/log"
)

const (
	metricUserAddSuccess = `biz.user.add.success`
	metricUserAddFailure = `biz.user.add.failure`
	metricUserAddTimings = `biz.user.add.timings`
)

type UserRepo interface {
	Save(context.Context, *ent.User) (*ent.User, error)
	Update(context.Context, *ent.User) (*ent.User, error)
	FindByID(ctx context.Context, id int) (*ent.User, error)
	FindByEmail(ctx context.Context, email string) (*ent.User, error)
	FindByPhone(ctx context.Context, phone string) (*ent.User, error)
}

type UserUsecase struct {
	repo   UserRepo
	metric metrics.Metrics
	logs   logger.Logger
}

func NewUserUsecase(repo UserRepo, metric metrics.Metrics, logs log.Logger) *UserUsecase {
	return &UserUsecase{
		repo:   repo,
		metric: metric,
		logs:   logger.NewHelper(logs, `ts`, log.DefaultTimestamp, `scope`, `biz/user`),
	}
}

func (u *UserUsecase) Add(ctx context.Context, dto *UserAddDTO) (*ent.User, error) {
	defer u.metric.NewTiming().Send(metricUserAddTimings)

	user := &ent.User{
		DisplayName:   dto.DisplayName,
		Type:          dto.Type,
		Email:         pointer.ToString(dto.Email),
		Phone:         pointer.ToString(dto.Phone),
		PasswordHash:  pointer.ToString(dto.PasswordHash),
		DeactivatedAt: dto.DeactivatedAt,
	}

	user, err := u.repo.Save(ctx, user)

	if err != nil {
		u.metric.Increment(metricUserAddFailure)
		u.logs.WithContext(ctx).Error(err)
		return nil, err
	}

	u.metric.Increment(metricUserAddSuccess)
	return nil, err
}
