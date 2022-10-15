package biz

import (
	"errors"
	"strings"
	"time"

	v1 "auth/api/auth/v1"
	"auth/internal/pkg/sanitize"
	"auth/internal/pkg/validate"
)

type Stats struct {
	IP        string  `json:"ip" validate:"omitempty,ip"`
	UserAgent string  `json:"userAgent" validate:"omitempty,max=4096"`
	DeviceID  *string `json:"deviceId" validate:"omitempty,max=4096"`
}

type LoginDTO struct {
	Username string `json:"username" validate:"required,min=3,max=255"`
	Password string `json:"password" validate:"required,min=8,max=255"`
	Remember bool   `json:"remember"`
	Stats    *Stats `json:"stats"`
}

func makeStats(s *v1.Stats) *Stats {
	if s == nil {
		return nil
	}
	return &Stats{
		IP:        s.Ip,
		UserAgent: s.UserAgent,
		DeviceID:  s.DeviceId,
	}
}

func (a *AuthUsecase) MakeLoginDTOFromLoginRequest(l *v1.LoginRequest) (*LoginDTO, error) {
	if l == nil {
		return nil, errors.New(`loginRequest is empty`)
	}
	dto := &LoginDTO{
		Username: sanitizeUsername(l.Username),
		Password: l.Password,
		Remember: l.Remember,
		Stats:    makeStats(l.Stats),
	}
	if err := validate.Default(dto); err != nil {
		return nil, err
	}
	return dto, nil
}

// CheckDTO is output DTO for Check() method
type CheckDTO struct {
	UserID           int
	UserType         string
	DisplayName      string
	Email            *string
	Phone            *string
	SessionUntil     time.Time
	SessionIP        string
	SessionUserAgent string
	SessionDeviceID  *string
}

type LoginByCodeDTO struct {
	Username string `json:"username" validate:"required,min=3,max=255"`
	Code     string `json:"code" validate:"required,min=3,max=64"`
	Remember bool   `json:"remember"`
	Stats    *Stats `json:"stats"`
}

func (a *AuthUsecase) MakeLoginByCodeFromLoginByCodeRequest(l *v1.LoginByCodeRequest) (*LoginByCodeDTO, error) {
	if l == nil {
		return nil, errors.New(`loginByCodeRequest is empty`)
	}
	dto := &LoginByCodeDTO{
		Username: sanitizeUsername(l.Username),
		Code:     strings.TrimSpace(l.Code),
		Remember: l.Remember,
		Stats:    makeStats(l.Stats),
	}
	if err := validate.Default(dto); err != nil {
		return nil, err
	}
	return dto, nil
}

type ResetPasswordDTO struct {
	Username string `json:"username" validate:"required,min=3,max=255"`
	Stats    *Stats `json:"stats"`
}

func (a *AuthUsecase) MakeResetPasswordDTO(l *v1.ResetPasswordRequest) (*ResetPasswordDTO, error) {
	if l == nil {
		return nil, errors.New(`resetPasswordRequest is empty`)
	}
	dto := &ResetPasswordDTO{
		Username: sanitizeUsername(l.Username),
		Stats:    makeStats(l.Stats),
	}
	if err := validate.Default(dto); err != nil {
		return nil, err
	}
	return dto, nil
}

type NewPasswordDTO struct {
	Username          string `json:"username" validate:"required,min=3,max=255"`
	PasswordResetHash string `json:"passwordResetHash" validate:"required,min=4,max=255"`
	Password          string `json:"password" validate:"required,min=8,max=255"`
	Stats             *Stats `json:"stats"`
}

func sanitizeUsername(username string) string {
	u := strings.TrimSpace(username)
	if !strings.Contains(u, `@`) {
		// case for phone as username
		u = sanitize.Phone(u)
	}
	return u
}

func (a *AuthUsecase) MakeNewPasswordDTO(l *v1.NewPasswordRequest) (*NewPasswordDTO, error) {
	if l == nil {
		return nil, errors.New(`newPasswordRequest is empty`)
	}
	dto := &NewPasswordDTO{
		Username:          sanitizeUsername(l.Username),
		PasswordResetHash: strings.TrimSpace(l.PasswordResetHash),
		Password:          l.Password,
		Stats:             makeStats(l.Stats),
	}
	if err := validate.Default(dto); err != nil {
		return nil, err
	}
	return dto, nil
}

type ChangePasswordDTO struct {
	Username    string `json:"username" validate:"required,min=3,max=255"`
	OldPassword string `json:"oldPassword" validate:"required,min=8,max=255"`
	NewPassword string `json:"newPassword" validate:"required,min=8,max=255"`
	Stats       *Stats `json:"stats"`
}

func (a *AuthUsecase) MakeChangePasswordDTO(l *v1.ChangePasswordRequest) (*ChangePasswordDTO, error) {
	if l == nil {
		return nil, errors.New(`changePasswordRequest is empty`)
	}
	dto := &ChangePasswordDTO{
		Username:    sanitizeUsername(l.Username),
		NewPassword: l.NewPassword,
		OldPassword: l.OldPassword,
		Stats:       makeStats(l.Stats),
	}
	if err := validate.Default(dto); err != nil {
		return nil, err
	}
	return dto, nil
}

type GenerateCodeDTO struct {
	Username string `json:"username" validate:"required,min=3,max=255"`
	Stats    *Stats `json:"stats"`
}

func (a *AuthUsecase) MakeGenerateCodeDTO(l *v1.GenerateCodeRequest) (*GenerateCodeDTO, error) {
	if l == nil {
		return nil, errors.New(`generateCodeRequest is empty`)
	}
	dto := &GenerateCodeDTO{
		Username: sanitizeUsername(l.Username),
		Stats:    makeStats(l.Stats),
	}
	if err := validate.Default(dto); err != nil {
		return nil, err
	}
	return dto, nil
}

type HistoryDTO struct {
	UserID int `json:"userId" validate:"required,min=1"`
	Limit  int `json:"limit" validate:"required,min=0,max=1000"`
	Offset int `json:"offset" validate:"min=0,max=10000"`
}

func (a *AuthUsecase) MakeHistoryDTO(l *v1.HistoryRequest) (*HistoryDTO, error) {
	if l == nil {
		return nil, errors.New(`historyRequest is empty`)
	}
	dto := &HistoryDTO{
		UserID: int(l.UserId),
		Limit:  int(l.Limit),
		Offset: int(l.Offset),
	}
	if err := validate.Default(dto); err != nil {
		return nil, err
	}
	return dto, nil
}
