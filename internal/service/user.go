package service

import (
	"context"
	"strings"

	authV1 "auth/api/auth/v1"
	"auth/ent"
	"auth/internal/biz"
	"auth/internal/pkg/logger"
	"auth/internal/pkg/metrics"
	"auth/internal/pkg/watcher"

	"github.com/go-kratos/kratos/v2/log"
	"google.golang.org/protobuf/types/known/timestamppb"
)

const (
	metricPrefixUser = `service.user`
)

type UserService struct {
	authV1.UnimplementedUserServer
	usecase *biz.UserUsecase
	metric  metrics.Metrics
	logger  *log.Helper
	watcher *watcher.Watcher
}

func NewUserService(usecase *biz.UserUsecase, metric metrics.Metrics, logs log.Logger) *UserService {
	loggerHelper := logger.NewHelper(logs, "ts", log.DefaultTimestamp, "scope", metricPrefixUser)
	return &UserService{
		usecase: usecase,
		metric:  metric,
		logger:  loggerHelper,
		watcher: watcher.New(metricPrefixUser, loggerHelper, metric).WithIgnoredErrorsChecks([]func(error) bool{
			authV1.IsValidationFailed,
		}),
	}
}

func (u *UserService) List(ctx context.Context, req *authV1.ListRequest) (*authV1.ListResponse, error) {
	var err error
	defer u.watcher.OnPreparedMethod(`List`).WithFields(map[string]any{
		"limit":  req.Limit,
		"offset": req.Offset,
	}).Results(func() (context.Context, error) {
		return ctx, err
	})

	list, err := u.usecase.List(ctx, req.Limit, req.Offset)
	if err != nil {
		if strings.Contains(err.Error(), `must be`) {
			return nil, authV1.ErrorValidationFailed(err.Error())
		}
		return nil, authV1.ErrorInternalError(err.Error())
	}

	users := make([]*authV1.ListUser, 0, len(list))
	for _, user := range list {
		users = append(users, transformEntUserToListUser(user))
	}
	return &authV1.ListResponse{
		Users: users,
	}, nil
}

func transformEntUserToListUser(user *ent.User) *authV1.ListUser {
	if user == nil {
		return nil
	}
	listUser := &authV1.ListUser{
		Id:             int64(user.ID),
		DisplayName:    user.DisplayName,
		Type:           user.Type,
		Email:          user.Email,
		Phone:          user.Phone,
		TelegramChatId: user.TelegramChatID,
		CreatedAt:      timestamppb.New(user.CreatedAt),
		UpdatedAt:      timestamppb.New(user.UpdatedAt),
		DeactivatedAt:  nil,
	}
	if user.DeactivatedAt != nil {
		listUser.DeactivatedAt = timestamppb.New(*user.DeactivatedAt)
	}
	return listUser
}

func (u *UserService) Add(ctx context.Context, req *authV1.AddRequest) (*authV1.AddResponse, error) {
	var err error
	defer u.watcher.OnPreparedMethod(`Add`).WithFields(map[string]any{
		"displayName": req.DisplayName,
	}).Results(func() (context.Context, error) {
		return ctx, err
	})

	addDTO, err := u.usecase.MakeUserAddDTO(req)
	if err != nil {
		return nil, authV1.ErrorValidationFailed(err.Error())
	}
	user, err := u.usecase.Add(ctx, addDTO)
	if err != nil {
		return nil, authV1.ErrorInternalError(err.Error())
	}
	return &authV1.AddResponse{
		Id: int64(user.ID),
	}, nil
}

func (u *UserService) Edit(ctx context.Context, req *authV1.EditRequest) (*authV1.UserNothing, error) {
	var err error
	defer u.watcher.OnPreparedMethod(`Edit`).WithFields(map[string]any{
		"userId": req.Id,
	}).Results(func() (context.Context, error) {
		return ctx, err
	})

	editDTO, err := u.usecase.MakeUserEditDTO(req)
	if err != nil {
		return nil, authV1.ErrorValidationFailed(err.Error())
	}
	_, err = u.usecase.Edit(ctx, editDTO)
	if err != nil {
		return nil, authV1.ErrorInternalError(err.Error())
	}
	return &authV1.UserNothing{}, nil
}

func (u *UserService) Activate(ctx context.Context, req *authV1.ActivateRequest) (*authV1.UserNothing, error) {
	var err error
	defer u.watcher.OnPreparedMethod(`Activate`).WithFields(map[string]any{
		"userId": req.Id,
	}).Results(func() (context.Context, error) {
		return ctx, err
	})

	activateDTO, err := u.usecase.MakeUserActivateDTO(req)
	if err != nil {
		return nil, authV1.ErrorValidationFailed(err.Error())
	}
	_, err = u.usecase.Activate(ctx, activateDTO.ID)
	if err != nil {
		return nil, authV1.ErrorInternalError(err.Error())
	}
	return &authV1.UserNothing{}, nil
}

func (u *UserService) Deactivate(ctx context.Context, req *authV1.DeactivateRequest) (*authV1.UserNothing, error) {
	var err error
	defer u.watcher.OnPreparedMethod(`Deactivate`).WithFields(map[string]any{
		"userId": req.Id,
	}).Results(func() (context.Context, error) {
		return ctx, err
	})

	deactivateDTO, err := u.usecase.MakeUserDeactivateDTO(req)
	if err != nil {
		return nil, authV1.ErrorValidationFailed(err.Error())
	}
	_, err = u.usecase.Deactivate(ctx, deactivateDTO.ID)
	if err != nil {
		return nil, authV1.ErrorInternalError(err.Error())
	}
	return &authV1.UserNothing{}, nil
}
