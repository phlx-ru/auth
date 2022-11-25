package service

import (
	"context"

	authV1 "auth/api/auth/v1"
	"auth/ent"
	"auth/internal/biz"
	"auth/internal/pkg/logger"
	"auth/internal/pkg/metrics"
	"auth/internal/pkg/strings"

	"github.com/AlekSi/pointer"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/transport"
	"google.golang.org/protobuf/types/known/timestamppb"
)

const (
	metricPrefixAuth = `service.auth`
)

type AuthService struct {
	authV1.UnimplementedAuthServer
	usecase *biz.AuthUsecase
	metric  metrics.Metrics
	logger  *log.Helper
}

func NewAuthService(
	usecase *biz.AuthUsecase,
	metric metrics.Metrics,
	logs log.Logger,
) *AuthService {
	return &AuthService{
		usecase: usecase,
		metric:  metric,
		logger:  logger.NewHelper(logs, "ts", log.DefaultTimestamp, "scope", "service/auth"),
	}
}

func (a *AuthService) postProcess(ctx context.Context, method string, err error) {
	if err != nil {
		a.logger.WithContext(ctx).Errorf(`auth service method "%s" failed: %v`, method, err)
		a.metric.Increment(strings.Metric(metricPrefixAuth, method, `failure`))
	} else {
		a.metric.Increment(strings.Metric(metricPrefixAuth, method, `success`))
	}
}

func (a *AuthService) Check(ctx context.Context, req *authV1.CheckRequest) (*authV1.CheckResponse, error) {
	method := `check`
	defer a.metric.NewTiming().Send(strings.Metric(metricPrefixAuth, method, `timings`))
	var err error
	defer func() { a.postProcess(ctx, method, err) }()

	dto, err := a.usecase.Check(ctx, req.Token)
	if ent.IsNotFound(err) {
		err = authV1.ErrorNotFound(err.Error())
	}
	if err != nil {
		return nil, err
	}
	return &authV1.CheckResponse{
		User: &authV1.CheckResponse_User{
			Id:          int64(dto.UserID),
			Type:        dto.UserType,
			DisplayName: dto.DisplayName,
			Email:       dto.Email,
			Phone:       dto.Phone,
		},
		Session: &authV1.CheckResponse_Session{
			Until:     timestamppb.New(dto.SessionUntil),
			Ip:        &dto.SessionIP,
			UserAgent: &dto.SessionUserAgent,
			DeviceId:  dto.SessionDeviceID,
		},
	}, nil
}

func (a *AuthService) Login(ctx context.Context, req *authV1.LoginRequest) (*authV1.LoginResponse, error) {
	method := `login`
	defer a.metric.NewTiming().Send(strings.Metric(metricPrefixAuth, method, `timings`))
	var err error
	defer func() { a.postProcess(ctx, method, err) }()

	if req.Stats == nil {
		req.Stats = statsFromRequestContext(ctx)
	}

	dto, err := a.usecase.MakeLoginDTOFromLoginRequest(req)
	if err != nil {
		err = authV1.ErrorValidationFailed(err.Error())
		return nil, err
	}

	token, expiredAt, err := a.usecase.Login(ctx, dto)
	switch err {
	case biz.ErrLoginTooOften:
		err = authV1.ErrorTooOften(err.Error())
	case biz.ErrWrongPassword:
		err = authV1.ErrorWrongInput(err.Error())
	default:
		if err != nil {
			err = authV1.ErrorInternalError(err.Error())
		}
	}
	if err != nil {
		return nil, err
	}

	return &authV1.LoginResponse{
		Token: token,
		Until: timestamppb.New(*expiredAt),
	}, nil
}

func (a *AuthService) LoginByCode(ctx context.Context, req *authV1.LoginByCodeRequest) (*authV1.LoginResponse, error) {
	method := `loginByCode`
	defer a.metric.NewTiming().Send(strings.Metric(metricPrefixAuth, method, `timings`))
	var err error
	defer func() { a.postProcess(ctx, method, err) }()

	if req.Stats == nil {
		req.Stats = statsFromRequestContext(ctx)
	}

	dto, err := a.usecase.MakeLoginByCodeFromLoginByCodeRequest(req)
	if err != nil {
		err = authV1.ErrorValidationFailed(err.Error())
		return nil, err
	}

	token, expiredAt, err := a.usecase.LoginByCode(ctx, dto)
	switch err {
	case biz.ErrLoginByCodeTooOften:
		err = authV1.ErrorTooOften(err.Error())
	case biz.ErrWrongCode:
		err = authV1.ErrorWrongInput(err.Error())
	default:
		if err != nil {
			err = authV1.ErrorInternalError(err.Error())
		}
	}
	if err != nil {
		return nil, err
	}

	return &authV1.LoginResponse{
		Token: token,
		Until: timestamppb.New(*expiredAt),
	}, nil
}

func (a *AuthService) ResetPassword(ctx context.Context, req *authV1.ResetPasswordRequest) (
	*authV1.AuthNothing,
	error,
) {
	method := `resetPassword`
	defer a.metric.NewTiming().Send(strings.Metric(metricPrefixAuth, method, `timings`))
	var err error
	defer func() { a.postProcess(ctx, method, err) }()

	if req.Stats == nil {
		req.Stats = statsFromRequestContext(ctx)
	}

	dto, err := a.usecase.MakeResetPasswordDTO(req)
	if err != nil {
		err = authV1.ErrorValidationFailed(err.Error())
		return nil, err
	}

	err = a.usecase.ResetPassword(ctx, dto)
	if err == biz.ErrResetPasswordTooOften {
		err = authV1.ErrorTooOften(err.Error())
	} else if err != nil {
		err = authV1.ErrorInternalError(err.Error())
	}
	if err != nil {
		return nil, err
	}

	return &authV1.AuthNothing{}, nil
}

func (a *AuthService) NewPassword(ctx context.Context, req *authV1.NewPasswordRequest) (*authV1.AuthNothing, error) {
	method := `newPassword`
	defer a.metric.NewTiming().Send(strings.Metric(metricPrefixAuth, method, `timings`))
	var err error
	defer func() { a.postProcess(ctx, method, err) }()

	if req.Stats == nil {
		req.Stats = statsFromRequestContext(ctx)
	}

	dto, err := a.usecase.MakeNewPasswordDTO(req)
	if err != nil {
		err = authV1.ErrorValidationFailed(err.Error())
		return nil, err
	}

	err = a.usecase.NewPassword(ctx, dto)
	switch err {
	case biz.ErrNewPasswordTooOften:
		err = authV1.ErrorTooOften(err.Error())
	case biz.ErrWrongResetHash:
		err = authV1.ErrorWrongInput(err.Error())
	default:
		if err != nil {
			err = authV1.ErrorInternalError(err.Error())
		}
	}
	if err != nil {
		return nil, err
	}

	return &authV1.AuthNothing{}, nil
}

func (a *AuthService) ChangePassword(ctx context.Context, req *authV1.ChangePasswordRequest) (
	*authV1.AuthNothing,
	error,
) {
	method := `changePassword`
	defer a.metric.NewTiming().Send(strings.Metric(metricPrefixAuth, method, `timings`))
	var err error
	defer func() { a.postProcess(ctx, method, err) }()

	if req.Stats == nil {
		req.Stats = statsFromRequestContext(ctx)
	}

	dto, err := a.usecase.MakeChangePasswordDTO(req)
	if err != nil {
		err = authV1.ErrorValidationFailed(err.Error())
		return nil, err
	}

	err = a.usecase.ChangePassword(ctx, dto)
	switch err {
	case biz.ErrChangePasswordTooOften:
		err = authV1.ErrorTooOften(err.Error())
	case biz.ErrWrongOldPassword:
		err = authV1.ErrorWrongInput(err.Error())
	default:
		if err != nil {
			err = authV1.ErrorInternalError(err.Error())
		}
	}
	if err != nil {
		return nil, err
	}

	return &authV1.AuthNothing{}, nil
}

func (a *AuthService) GenerateCode(ctx context.Context, req *authV1.GenerateCodeRequest) (
	*authV1.AuthNothing,
	error,
) {
	method := `generateCode`
	defer a.metric.NewTiming().Send(strings.Metric(metricPrefixAuth, method, `timings`))
	var err error
	defer func() { a.postProcess(ctx, method, err) }()

	if req.Stats == nil {
		req.Stats = statsFromRequestContext(ctx)
	}

	dto, err := a.usecase.MakeGenerateCodeDTO(req)
	if err != nil {
		err = authV1.ErrorValidationFailed(err.Error())
		return nil, err
	}

	err = a.usecase.GenerateCode(ctx, dto)
	if err == biz.ErrGenerateCodeTooOften {
		err = authV1.ErrorTooOften(err.Error())
	} else if err != nil {
		err = authV1.ErrorInternalError(err.Error())
	}
	if err != nil {
		return nil, err
	}

	return &authV1.AuthNothing{}, nil
}

func (a *AuthService) History(ctx context.Context, req *authV1.HistoryRequest) (*authV1.HistoryResponse, error) {
	method := `history`
	defer a.metric.NewTiming().Send(strings.Metric(metricPrefixAuth, method, `timings`))
	var err error
	defer func() { a.postProcess(ctx, method, err) }()

	resp := &authV1.HistoryResponse{
		Items: make([]*authV1.HistoryItem, 0),
	}

	dto, err := a.usecase.MakeHistoryDTO(req)
	if err != nil {
		err = authV1.ErrorValidationFailed(err.Error())
		return nil, err
	}

	items, err := a.usecase.History(ctx, dto)
	if err != nil {
		err = authV1.ErrorInternalError(err.Error())
		return nil, err
	}

	for _, item := range items {
		history := &authV1.HistoryItem{
			Id:        int64(item.ID),
			When:      timestamppb.New(item.CreatedAt),
			Event:     item.Event,
			Ip:        pointer.GetString(item.IP),
			UserAgent: pointer.GetString(item.UserAgent),
		}
		resp.Items = append(resp.Items, history)
	}

	return resp, nil
}

func statsFromRequestContext(ctx context.Context) *authV1.Stats {
	tr, ok := transport.FromServerContext(ctx)
	if !ok {
		return nil
	}
	h := tr.RequestHeader()
	return &authV1.Stats{
		Ip:        h.Get("X-Real-Ip"),
		UserAgent: h.Get("User-Agent"),
		DeviceId:  pointer.ToString(h.Get("DeviceId")),
	}
}
