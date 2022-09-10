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
	DisplayName    string `validate:"required"`
	Type           string `validate:"required,user_type"`
	Phone          string `validate:"required_without=email,min=10,numeric,startswith=9"`
	Email          string `validate:"required_without=phone,email"`
	TelegramChatID string `validate:"required_without=email,numeric"`
	PasswordHash   string `validate:"required,min=8,max=256"`
	DeactivatedAt  *time.Time
}

func MakeUserAddDTOFromAddRequest(a *v1.AddRequest) (*UserAddDTO, error) {
	if a == nil {
		return nil, errors.New(`addRequest is empty`)
	}
	dto := &UserAddDTO{
		DisplayName:   a.DisplayName,
		Type:          a.Type,
		Phone:         sanitize.Phone(a.Phone),
		Email:         a.Email,
		PasswordHash:  secrets.MustMakeHash(a.Password),
		DeactivatedAt: nil,
	}
	if !a.Activated {
		dto.DeactivatedAt = pointer.ToTime(time.Now())
	}
	if err := validate.Default(dto); err != nil {
		return nil, err
	}
	return dto, nil
}
