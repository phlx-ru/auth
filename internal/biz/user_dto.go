package biz

import (
	"errors"
	"time"

	v1 "auth/api/auth/v1"
	"auth/internal/pkg/sanitize"
	"auth/internal/pkg/secrets"
	"auth/internal/pkg/validate"

	"github.com/AlekSi/pointer"
)

type UserAddDTO struct {
	DisplayName    string `validate:"required,min=3,max=255"`
	Type           string `validate:"required,user_type"`
	Phone          string `validate:"required_without=Email,omitempty,min=10,max=10,numeric,startswith=9"`
	Email          string `validate:"required_without=Phone,omitempty,email"`
	TelegramChatID string `validate:"omitempty,numeric"`
	Password       string `validate:"required,min=8,max=255"`
	PasswordHash   string
	DeactivatedAt  *time.Time
}

func (u *UserUsecase) MakeUserAddDTO(request *v1.AddRequest) (*UserAddDTO, error) {
	if request == nil {
		return nil, errors.New(`addRequest is empty`)
	}
	dto := &UserAddDTO{
		DisplayName:   request.DisplayName,
		Type:          request.Type,
		DeactivatedAt: nil,
	}
	if request.Phone != nil {
		dto.Phone = sanitize.Phone(*request.Phone)
	}
	if request.Email != nil {
		dto.Email = *request.Email
	}
	if request.TelegramChatId != nil {
		dto.TelegramChatID = *request.TelegramChatId
	}
	if request.Password != nil {
		dto.Password = *request.Password
		dto.PasswordHash = secrets.MustMakeHash(*request.Password)
	}
	if request.Deactivated {
		dto.DeactivatedAt = pointer.ToTime(time.Now())
	}
	if err := validate.Default(dto); err != nil {
		return nil, err
	}
	return dto, nil
}

type UserEditDTO struct {
	ID             int     `validate:"required,min=1"`
	DisplayName    *string `validate:"omitempty,min=3,max=255"`
	Type           *string `validate:"omitempty,user_type"`
	Phone          *string `validate:"omitempty,min=10,numeric,startswith=9"`
	Email          *string `validate:"omitempty,email"`
	TelegramChatID *string `validate:"omitempty,numeric"`
	PasswordHash   *string `validate:"omitempty,min=8,max=256"`
}

func (u *UserUsecase) MakeUserEditDTO(request *v1.EditRequest) (*UserEditDTO, error) {
	if request == nil {
		return nil, errors.New(`editRequest is empty`)
	}
	dto := &UserEditDTO{
		ID:          int(request.Id),
		DisplayName: request.DisplayName,
		Type:        request.Type,
	}
	if request.Phone != nil && *request.Phone != "" {
		dto.Phone = pointer.ToString(sanitize.Phone(*request.Phone))
	}
	if request.Email != nil && *request.Email != "" {
		dto.Email = request.Email
	}
	if request.TelegramChatId != nil && *request.TelegramChatId != "" {
		dto.TelegramChatID = request.TelegramChatId
	}
	if request.Password != nil {
		dto.PasswordHash = pointer.ToString(secrets.MustMakeHash(*request.Password))
	}
	if err := validate.Default(dto); err != nil {
		return nil, err
	}
	if request.Phone != nil && *request.Phone == "" {
		dto.Phone = pointer.ToString("")
	}
	if request.Email != nil && *request.Email == "" {
		dto.Email = pointer.ToString("")
	}
	if request.TelegramChatId != nil && *request.TelegramChatId == "" {
		dto.TelegramChatID = pointer.ToString("")
	}
	return dto, nil
}

type UserActivateDTO struct {
	ID int `validate:"required,min=1"`
}

func (u *UserUsecase) MakeUserActivateDTO(request *v1.ActivateRequest) (*UserActivateDTO, error) {
	if request == nil {
		return nil, errors.New(`activateRequest is empty`)
	}
	dto := &UserActivateDTO{
		ID: int(request.Id),
	}
	if err := validate.Default(dto); err != nil {
		return nil, err
	}
	return dto, nil
}

type UserDeactivateDTO struct {
	ID int `validate:"required,min=1"`
}

func (u *UserUsecase) MakeUserDeactivateDTO(request *v1.DeactivateRequest) (*UserDeactivateDTO, error) {
	if request == nil {
		return nil, errors.New(`deactivateRequest is empty`)
	}
	dto := &UserDeactivateDTO{
		ID: int(request.Id),
	}
	if err := validate.Default(dto); err != nil {
		return nil, err
	}
	return dto, nil
}
