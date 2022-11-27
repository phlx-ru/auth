package service

import (
	"context"

	authV1 "auth/api/auth/v1"
	"auth/internal/biz"
	"auth/internal/pkg/logger"
	"auth/internal/pkg/metrics"
	"auth/internal/pkg/strings"

	"github.com/go-kratos/kratos/v2/log"
)

const (
	metricPrefixUser = `service.user`
)

type UserService struct {
	authV1.UnimplementedUserServer

	usecase *biz.UserUsecase

	metric metrics.Metrics

	logger *log.Helper
}

func NewUserService(usecase *biz.UserUsecase, metric metrics.Metrics, logs log.Logger) *UserService {
	return &UserService{
		usecase: usecase,
		metric:  metric,
		logger:  logger.NewHelper(logs, "ts", log.DefaultTimestamp, "scope", "service/user"),
	}
}

func (u *UserService) postProcess(ctx context.Context, method string, err error) {
	if err != nil {
		u.logger.WithContext(ctx).Errorf(`user service method "%s" failed: %v`, method, err)
		u.metric.Increment(strings.Metric(metricPrefixUser, method, `failure`))
	} else {
		u.metric.Increment(strings.Metric(metricPrefixUser, method, `success`))
	}
}

func (u *UserService) Add(ctx context.Context, req *authV1.AddRequest) (*authV1.AddResponse, error) {
	method := `add`
	defer u.metric.NewTiming().Send(strings.Metric(metricPrefixUser, method, `timings`))
	var err error
	defer func() { u.postProcess(ctx, method, err) }()

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
	method := `edit`
	defer u.metric.NewTiming().Send(strings.Metric(metricPrefixUser, method, `timings`))
	var err error
	defer func() { u.postProcess(ctx, method, err) }()

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
	method := `activate`
	defer u.metric.NewTiming().Send(strings.Metric(metricPrefixUser, method, `timings`))
	var err error
	defer func() { u.postProcess(ctx, method, err) }()

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
	method := `deactivate`
	defer u.metric.NewTiming().Send(strings.Metric(metricPrefixUser, method, `timings`))
	var err error
	defer func() { u.postProcess(ctx, method, err) }()

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
