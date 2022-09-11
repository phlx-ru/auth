package biz

import (
	"errors"
	"strconv"
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
	Phone          string `validate:"required_without=email,min=10,numeric,startswith=9"`
	Email          string `validate:"required_without=phone,email"`
	TelegramChatID string `validate:"required_without=email,numeric"`
	PasswordHash   string `validate:"required,min=8,max=255"`
	DeactivatedAt  *time.Time
}

func (u *UserUsecase) MakeUserAddDTO(a *v1.AddRequest) (*UserAddDTO, error) {
	if a == nil {
		return nil, errors.New(`addRequest is empty`)
	}
	dto := &UserAddDTO{
		DisplayName:   a.DisplayName,
		Type:          a.Type,
		DeactivatedAt: nil,
	}
	if a.Phone != nil {
		dto.Phone = sanitize.Phone(*a.Phone)
	}
	if a.Email != nil {
		dto.Email = *a.Email
	}
	if a.TelegramChatId != nil {
		dto.TelegramChatID = *a.TelegramChatId
	}
	if a.Password != nil {
		dto.PasswordHash = secrets.MustMakeHash(*a.Password)
	}
	if !a.Activated {
		dto.DeactivatedAt = pointer.ToTime(time.Now())
	}
	if err := validate.Default(dto); err != nil {
		return nil, err
	}
	return dto, nil
}

type UserEditDTO struct {
	ID             int     `validate:"required,min:1"`
	DisplayName    *string `validate:"min=3,max=255"`
	Type           *string `validate:"user_type"`
	Phone          *string `validate:"min=10,numeric,startswith=9"`
	Email          *string `validate:"email"`
	TelegramChatID *string `validate:"numeric"`
	PasswordHash   *string `validate:"min=8,max=256"`
}

func (u *UserUsecase) MakeUserEditDTO(e *v1.EditRequest) (*UserEditDTO, error) {
	if e == nil {
		return nil, errors.New(`editRequest is empty`)
	}
	id, err := strconv.ParseInt(e.Id, 10, 64)
	if err != nil {
		return nil, err
	}
	dto := &UserEditDTO{
		ID:          int(id),
		DisplayName: e.DisplayName,
		Type:        e.Type,
	}
	if e.Phone != nil && *e.Phone != "" {
		dto.Phone = pointer.ToString(sanitize.Phone(*e.Phone))
	}
	if e.Email != nil && *e.Email != "" {
		dto.Email = e.Email
	}
	if e.TelegramChatId != nil && *e.TelegramChatId != "" {
		dto.TelegramChatID = e.TelegramChatId
	}
	if e.Password != nil {
		dto.PasswordHash = pointer.ToString(secrets.MustMakeHash(*e.Password))
	}
	if err := validate.Default(dto); err != nil {
		return nil, err
	}
	if e.Phone != nil && *e.Phone == "" {
		dto.Phone = pointer.ToString("")
	}
	if e.Email != nil && *e.Email == "" {
		dto.Email = pointer.ToString("")
	}
	if e.TelegramChatId != nil && *e.TelegramChatId == "" {
		dto.TelegramChatID = pointer.ToString("")
	}
	return dto, nil
}
