package service

import (
	"context"

	authV1 "auth/api/auth/v1"
	"auth/internal/biz"
	"auth/internal/pkg/logger"
	"auth/internal/pkg/metrics"
	"auth/internal/pkg/watcher"

	"github.com/go-kratos/kratos/v2/log"
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
