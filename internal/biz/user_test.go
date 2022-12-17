package biz

import (
	"context"
	"testing"
	"time"

	"github.com/AlekSi/pointer"
	"github.com/stretchr/testify/require"

	v1 "auth/api/auth/v1"
	"auth/ent"
	"auth/internal/pkg/stubs"
)

func TestUserUsecase_Add(t *testing.T) {
	mockedUserRepo := &userRepoMock{
		ActivateFunc: nil,
		CreateFunc: func(contextMoqParam context.Context, user *ent.User) (*ent.User, error) {
			return &ent.User{
				ID:                     1,
				DisplayName:            user.DisplayName,
				Type:                   user.Type,
				Email:                  user.Email,
				Phone:                  nil,
				TelegramChatID:         nil,
				PasswordHash:           user.PasswordHash,
				PasswordReset:          nil,
				PasswordResetExpiredAt: nil,
				CreatedAt:              time.Time{},
				UpdatedAt:              time.Time{},
				DeactivatedAt:          nil,
			}, nil
		},
		DeactivateFunc:  nil,
		FindByEmailFunc: nil,
		FindByIDFunc:    nil,
		FindByPhoneFunc: nil,
		UpdateFunc:      nil,
	}

	mockedMetrics := stubs.NewMetricsMuted()

	mockedLog := stubs.NewLoggerMuted()

	userUsecase := NewUserUsecase(mockedUserRepo, mockedMetrics, mockedLog)

	ctx := context.Background()

	userAddDTO, err := userUsecase.MakeUserAddDTO(&v1.AddRequest{
		DisplayName:    "Johnny Silverhand",
		Type:           "admin",
		Phone:          nil,
		Email:          pointer.ToString("silver.rocker@arasaka.lies"),
		TelegramChatId: nil,
		Password:       pointer.ToString("smasherMustDie"),
		Deactivated:    false,
	})
	require.NoError(t, err)

	actualUser, err := userUsecase.Add(ctx, userAddDTO)
	require.NoError(t, err)
	require.Equal(t, userAddDTO.DisplayName, actualUser.DisplayName)
	require.Equal(t, userAddDTO.Type, actualUser.Type)
}
