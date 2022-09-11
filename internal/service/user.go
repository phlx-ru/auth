package service

import (
	"context"
	"fmt"

	v1 "auth/api/auth/v1"
	"auth/internal/biz"
	"auth/internal/pkg/logger"
	"auth/internal/pkg/metrics"

	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/log"
)

const (
	ErrorReasonCommon = `COMMON` // TODO Make custom errors

	metricUserAddTimings = `service.user.add.timings`
	metricUserAddSuccess = `service.user.add.success`
	metricUserAddFailure = `service.user.add.failure`
)

type UserService struct {
	v1.UnimplementedUserServer

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

func (s *UserService) Add(ctx context.Context, req *v1.AddRequest) (*v1.AddResponse, error) {
	defer s.metric.NewTiming().Send(metricUserAddTimings)
	var err error
	defer func() {
		if err != nil {
			s.logger.WithContext(ctx).Error(fmt.Sprintf("failed to add user: %v", err))
			s.metric.Increment(metricUserAddFailure)
		} else {
			s.metric.Increment(metricUserAddSuccess)
		}
	}()
	addDTO, err := s.usecase.MakeUserAddDTO(req)
	if err != nil {
		return nil, err
	}
	user, err := s.usecase.Add(ctx, addDTO)
	if err != nil {
		return nil, errors.InternalServer(ErrorReasonCommon, err.Error())
	}
	return &v1.AddResponse{
		Id: fmt.Sprintf(`%d`, user.ID),
	}, nil
}

func (s *UserService) Edit(ctx context.Context, req *v1.EditRequest) (*v1.UserNothing, error) {
	return &v1.UserNothing{}, nil
}
func (s *UserService) Activate(ctx context.Context, req *v1.ActivateRequest) (*v1.UserNothing, error) {
	return &v1.UserNothing{}, nil
}
func (s *UserService) Deactivate(ctx context.Context, req *v1.DeactivateRequest) (*v1.UserNothing, error) {
	return &v1.UserNothing{}, nil
}
