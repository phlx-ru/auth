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

	metricUserEditSuccess = `biz.user.edit.success`
	metricUserEditFailure = `biz.user.edit.failure`
	metricUserEditTimings = `biz.user.edit.timings`
)

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

	user, err := u.repo.Create(ctx, user)

	if err != nil {
		u.metric.Increment(metricUserAddFailure)
		u.logs.WithContext(ctx).Error(err)
		return nil, err
	}

	u.metric.Increment(metricUserAddSuccess)
	return user, err
}

func (u *UserUsecase) Edit(ctx context.Context, dto *UserEditDTO) (*ent.User, error) {
	defer u.metric.NewTiming().Send(metricUserEditTimings)

	var err error
	defer func() {
		if err != nil {
			u.logs.WithContext(ctx).Errorf(`failed to edit user: %v`, err)
			u.metric.Increment(metricUserEditFailure)
		} else {
			u.metric.Increment(metricUserEditSuccess)
		}
	}()

	user, err := u.repo.FindByID(ctx, dto.ID)
	if err != nil {
		return nil, err
	}

	if dto.Type != nil {
		user.Type = *dto.Type
	}

	if dto.Email != nil {
		if *dto.Email == "" {
			user.Email = nil
		} else {
			user.Email = dto.Email
		}
	}

	if dto.Phone != nil {
		if *dto.Phone == "" {
			user.Phone = nil
		} else {
			user.Phone = dto.Phone
		}
	}

	if dto.TelegramChatID != nil {
		if *dto.TelegramChatID == "" {
			user.TelegramChatID = nil
		} else {
			user.TelegramChatID = dto.TelegramChatID
		}
	}

	if dto.DisplayName != nil {
		user.DisplayName = *dto.DisplayName
	}

	if dto.PasswordHash != nil {
		if *dto.PasswordHash == "" {
			user.PasswordHash = nil
		} else {
			user.PasswordHash = dto.PasswordHash
		}
	}

	user, err = u.repo.Update(ctx, user)
	return user, err
}
