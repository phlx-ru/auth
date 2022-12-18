package biz

import (
	"context"
	"errors"
	"fmt"

	"auth/ent"
	"auth/internal/pkg/logger"
	"auth/internal/pkg/metrics"
	"auth/internal/pkg/watcher"

	"github.com/AlekSi/pointer"
	"github.com/go-kratos/kratos/v2/log"
)

const (
	metricPrefixUser = `biz.user`
)

type UserUsecase struct {
	repo    userRepo
	metric  metrics.Metrics
	logger  logger.Logger
	watcher *watcher.Watcher
}

func NewUserUsecase(repo userRepo, metric metrics.Metrics, logs log.Logger) *UserUsecase {
	loggerHelper := logger.NewHelper(logs, `ts`, log.DefaultTimestamp, `scope`, metricPrefixUser)
	return &UserUsecase{
		repo:    repo,
		metric:  metric,
		logger:  loggerHelper,
		watcher: watcher.New(metricPrefixUser, loggerHelper, metric),
	}
}

func (u *UserUsecase) Add(ctx context.Context, dto *UserAddDTO) (*ent.User, error) {
	var err error
	if dto == nil {
		err = errors.New("userAddDTO is empty")
		return nil, err
	}
	defer u.watcher.OnPreparedMethod(`Add`).WithFields(map[string]any{
		"displayName": dto.DisplayName,
	}).Results(func() (context.Context, error) {
		return ctx, err
	})

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
	var err error
	if dto == nil {
		err = errors.New("userEditDTO is empty")
		return nil, err
	}
	defer u.watcher.OnPreparedMethod(`Edit`).WithFields(map[string]any{
		"id": dto.ID,
	}).Results(func() (context.Context, error) {
		return ctx, err
	})

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
	var err error
	if userID <= 0 {
		return nil, fmt.Errorf("user id must be positive integer, get %d", userID)
	}
	defer u.watcher.OnPreparedMethod(`Activate`).WithFields(map[string]any{
		"userID": userID,
	}).Results(func() (context.Context, error) {
		return ctx, err
	})

	var user *ent.User
	user, err = u.repo.Activate(ctx, userID)
	return user, err
}

func (u *UserUsecase) Deactivate(ctx context.Context, userID int) (*ent.User, error) {
	var err error
	if userID <= 0 {
		return nil, fmt.Errorf("user id must be positive integer, get %d", userID)
	}
	defer u.watcher.OnPreparedMethod(`Deactivate`).WithFields(map[string]any{
		"userID": userID,
	}).Results(func() (context.Context, error) {
		return ctx, err
	})

	var user *ent.User
	user, err = u.repo.Deactivate(ctx, userID)
	return user, err
}
