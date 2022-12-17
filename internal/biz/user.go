package biz

import (
	"context"
	"fmt"

	"auth/ent"
	"auth/internal/pkg/logger"
	"auth/internal/pkg/metrics"
	"auth/internal/pkg/strings"

	"github.com/AlekSi/pointer"
	"github.com/go-kratos/kratos/v2/log"
)

const (
	metricPrefixUser = `biz.user`
)

type UserUsecase struct {
	repo   userRepo
	metric metrics.Metrics
	logger logger.Logger
}

func NewUserUsecase(repo userRepo, metric metrics.Metrics, logs log.Logger) *UserUsecase {
	return &UserUsecase{
		repo:   repo,
		metric: metric,
		logger: logger.NewHelper(logs, `ts`, log.DefaultTimestamp, `scope`, `biz/user`),
	}
}

func (u *UserUsecase) postProcess(ctx context.Context, method string, err error) {
	if err != nil {
		u.logger.WithContext(ctx).Errorf(`user service method "%s" failed: %v`, method, err)
		u.metric.Increment(strings.Metric(metricPrefixUser, method, `failure`))
	} else {
		u.metric.Increment(strings.Metric(metricPrefixUser, method, `success`))
	}
}

func (u *UserUsecase) Add(ctx context.Context, dto *UserAddDTO) (*ent.User, error) {
	method := `add`
	defer u.metric.NewTiming().Send(strings.Metric(metricPrefixUser, method, `timings`))
	var err error
	defer func() { u.postProcess(ctx, method, err) }()

	user := &ent.User{
		DisplayName:   dto.DisplayName,
		Type:          dto.Type,
		DeactivatedAt: dto.DeactivatedAt,
	}

	if dto.Email != "" {
		user.Email = pointer.ToString(dto.Email)
	}

	if dto.Phone != "" {
		user.Phone = pointer.ToString(dto.Phone)
	}

	if dto.PasswordHash != "" {
		user.PasswordHash = pointer.ToString(dto.PasswordHash)
	}

	user, err = u.repo.Create(ctx, user)
	return user, err
}

func (u *UserUsecase) Edit(ctx context.Context, dto *UserEditDTO) (*ent.User, error) {
	method := `edit`
	defer u.metric.NewTiming().Send(strings.Metric(metricPrefixUser, method, `timings`))
	var err error
	defer func() { u.postProcess(ctx, method, err) }()

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

func (u *UserUsecase) Activate(ctx context.Context, userID int) (*ent.User, error) {
	method := `activate`
	defer u.metric.NewTiming().Send(strings.Metric(metricPrefixUser, method, `timings`))
	var err error
	defer func() { u.postProcess(ctx, method, err) }()

	if userID <= 0 {
		return nil, fmt.Errorf("user id must be positive integer, get %d", userID)
	}

	var user *ent.User
	user, err = u.repo.Activate(ctx, userID)
	return user, err
}

func (u *UserUsecase) Deactivate(ctx context.Context, userID int) (*ent.User, error) {
	method := `deactivate`
	defer u.metric.NewTiming().Send(strings.Metric(metricPrefixUser, method, `timings`))
	var err error
	defer func() { u.postProcess(ctx, method, err) }()

	if userID <= 0 {
		return nil, fmt.Errorf("user id must be positive integer, get %d", userID)
	}

	var user *ent.User
	user, err = u.repo.Deactivate(ctx, userID)
	return user, err
}
