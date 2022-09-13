package service

import (
	"context"
	"net/http"

	authV1 "auth/api/auth/v1"
	"auth/ent"
	"auth/internal/biz"
	"auth/internal/pkg/logger"
	"auth/internal/pkg/metrics"

	"github.com/AlekSi/pointer"
	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/transport"
	"google.golang.org/protobuf/types/known/timestamppb"
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

func (s *AuthService) Check(ctx context.Context, req *authV1.CheckRequest) (*authV1.CheckResponse, error) {
	dto, err := s.usecase.Check(ctx, req.Token)
	if ent.IsNotFound(err) {
		return nil, errors.New(
			http.StatusNotFound,
			`NOT_FOUND`,
			`session not found or expired`,
		) // TODO change to generated error
	}
	if err != nil {
		return nil, err
	}
	return &authV1.CheckResponse{
		User: &authV1.CheckResponse_User{
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

func (s *AuthService) Login(ctx context.Context, req *authV1.LoginRequest) (*authV1.LoginResponse, error) {
	if req.Stats == nil {
		req.Stats = statsFromRequestContext(ctx)
	}

	dto, err := s.usecase.MakeLoginDTOFromLoginRequest(req)
	if err != nil {
		return nil, err
	}

	token, expiredAt, err := s.usecase.Login(ctx, dto)
	if err != nil {
		return nil, err
	}

	return &authV1.LoginResponse{
		Token: token,
		Until: timestamppb.New(*expiredAt),
	}, nil
}

func (s *AuthService) LoginByCode(ctx context.Context, req *authV1.LoginByCodeRequest) (*authV1.LoginResponse, error) {
	if req.Stats == nil {
		req.Stats = statsFromRequestContext(ctx)
	}

	dto, err := s.usecase.MakeLoginByCodeFromLoginByCodeRequest(req)
	if err != nil {
		return nil, err
	}

	token, expiredAt, err := s.usecase.LoginByCode(ctx, dto)
	if err != nil {
		return nil, err
	}

	return &authV1.LoginResponse{
		Token: token,
		Until: timestamppb.New(*expiredAt),
	}, nil
}

func (s *AuthService) ResetPassword(ctx context.Context, req *authV1.ResetPasswordRequest) (
	*authV1.AuthNothing,
	error,
) {
	if req.Stats == nil {
		req.Stats = statsFromRequestContext(ctx)
	}

	dto, err := s.usecase.MakeResetPasswordDTO(req)
	if err != nil {
		return nil, err
	}

	err = s.usecase.ResetPassword(ctx, dto)
	if err != nil {
		return nil, err
	}

	return &authV1.AuthNothing{}, nil
}

func (s *AuthService) NewPassword(ctx context.Context, req *authV1.NewPasswordRequest) (*authV1.AuthNothing, error) {
	if req.Stats == nil {
		req.Stats = statsFromRequestContext(ctx)
	}

	dto, err := s.usecase.MakeNewPasswordDTO(req)
	if err != nil {
		return nil, err
	}

	err = s.usecase.NewPassword(ctx, dto)
	if err != nil {
		return nil, err
	}

	return &authV1.AuthNothing{}, nil
}

func (s *AuthService) ChangePassword(ctx context.Context, req *authV1.ChangePasswordRequest) (
	*authV1.AuthNothing,
	error,
) {
	if req.Stats == nil {
		req.Stats = statsFromRequestContext(ctx)
	}

	dto, err := s.usecase.MakeChangePasswordDTO(req)
	if err != nil {
		return nil, err
	}

	err = s.usecase.ChangePassword(ctx, dto)
	if err != nil {
		return nil, err
	}

	return &authV1.AuthNothing{}, nil
}

func (s *AuthService) GenerateCode(ctx context.Context, req *authV1.GenerateCodeRequest) (
	*authV1.AuthNothing,
	error,
) {
	if req.Stats == nil {
		req.Stats = statsFromRequestContext(ctx)
	}

	dto, err := s.usecase.MakeGenerateCodeDTO(req)
	if err != nil {
		return nil, err
	}

	err = s.usecase.GenerateCode(ctx, dto)
	if err == biz.ErrGenerateCodeTooOften { // TODO MAKE ALL BIZ LOGIC ERROR PROCESSING WITH REASONS
		return nil, errors.New(http.StatusBadRequest, `generate_code_too_often`, err.Error())
	}
	if err != nil {
		return nil, err
	}

	return &authV1.AuthNothing{}, nil
}

func (s *AuthService) History(ctx context.Context, req *authV1.HistoryRequest) (*authV1.HistoryResponse, error) {
	resp := &authV1.HistoryResponse{
		Items: make([]*authV1.HistoryItem, 0),
	}

	dto, err := s.usecase.MakeHistoryDTO(req)
	if err != nil {
		return nil, err
	}

	items, err := s.usecase.History(ctx, dto)
	if err != nil {
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
		Ip:        h.Get("X-Real-Ip"), // TODO Maybe X-Forwarded-For
		UserAgent: h.Get("User-Agent"),
		DeviceId:  pointer.ToString(h.Get("DeviceId")),
	}
}
