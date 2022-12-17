package biz

import (
	"fmt"
	"strings"
	"testing"

	"github.com/AlekSi/pointer"
	"github.com/stretchr/testify/require"

	v1 "auth/api/auth/v1"
)

const (
	defaultPasswordHash = "$2a$10$and53symbolsABCDEFGHIJKLMNOPQRSTUVWXYZABCDEFGHIJKLMNO"
)

//goland:noinspection GoErrorStringFormat
func TestUserUsecase_MakeUserAddDTO(t *testing.T) {
	testCases := []struct {
		name           string
		request        *v1.AddRequest
		expectedResult *UserAddDTO
		expectedError  error
	}{
		{
			name:          "empty",
			request:       &v1.AddRequest{},
			expectedError: fmt.Errorf("Field validation for 'DisplayName' failed on the 'required' tag"),
		},
		{
			name: "display_name_too_small",
			request: &v1.AddRequest{
				DisplayName: "YO",
			},
			expectedError: fmt.Errorf("Field validation for 'DisplayName' failed on the 'min' tag"),
		},
		{
			name: "display_name_too_big",
			request: &v1.AddRequest{
				DisplayName: strings.Repeat("Ovuvuevuevue Enyetuenwuevue Ugbemugbem Osas", 10),
			},
			expectedError: fmt.Errorf("Field validation for 'DisplayName' failed on the 'max' tag"),
		},
		{
			name: "incorrect_type",
			request: &v1.AddRequest{
				DisplayName: "John Doe",
				Type:        "anonymous",
			},
			expectedError: fmt.Errorf("Field validation for 'Type' failed on the 'user_type' tag"),
		},
		{
			name: "incorrect_email",
			request: &v1.AddRequest{
				DisplayName: "John Doe",
				Type:        "admin",
				Email:       pointer.ToString("john.doe@mail.com@unexpected"),
			},
			expectedError: fmt.Errorf("Field validation for 'Email' failed on the 'email' tag"),
		},
		{
			name: "password_too_small",
			request: &v1.AddRequest{
				DisplayName: "John Doe",
				Type:        "admin",
				Email:       pointer.ToString("john.doe@mail.com"),
				Password:    pointer.ToString("qwertyi"),
			},
			expectedError: fmt.Errorf("Field validation for 'Password' failed on the 'min' tag"),
		},
		{
			name: "password_too_big",
			request: &v1.AddRequest{
				DisplayName: "John Doe",
				Type:        "admin",
				Email:       pointer.ToString("john.doe@mail.com"),
				Password:    pointer.ToString(strings.Repeat("a", 256)),
			},
			expectedError: fmt.Errorf("Field validation for 'Password' failed on the 'max' tag"),
		},
		{
			name: "phone_required_without_email",
			request: &v1.AddRequest{
				DisplayName: "John Doe",
				Type:        "admin",
				Password:    pointer.ToString("qwertyuio"),
			},
			expectedError: fmt.Errorf("Field validation for 'Phone' failed on the 'required_without' tag"),
		},
		{
			name: "phone_too_small",
			request: &v1.AddRequest{
				DisplayName: "John Doe",
				Type:        "admin",
				Email:       pointer.ToString("john.doe@mail.com"),
				Password:    pointer.ToString("qwertyuio"),
				Phone:       pointer.ToString("123456789"),
			},
			expectedError: fmt.Errorf("Field validation for 'Phone' failed on the 'min' tag"),
		},
		{
			name: "phone_too_big",
			request: &v1.AddRequest{
				DisplayName: "John Doe",
				Type:        "admin",
				Email:       pointer.ToString("john.doe@mail.com"),
				Password:    pointer.ToString("qwertyuio"),
				Phone:       pointer.ToString("12345678901"),
			},
			expectedError: fmt.Errorf("Field validation for 'Phone' failed on the 'max' tag"),
		},
		{
			name: "phone_not_starts_with_9",
			request: &v1.AddRequest{
				DisplayName: "John Doe",
				Type:        "admin",
				Email:       pointer.ToString("john.doe@mail.com"),
				Password:    pointer.ToString("qwertyuio"),
				Phone:       pointer.ToString("88005553555"),
			},
			expectedError: fmt.Errorf("Field validation for 'Phone' failed on the 'startswith' tag"),
		},
		{
			name: "telegram_chat_not_numeric",
			request: &v1.AddRequest{
				DisplayName:    "John Doe",
				Type:           "admin",
				Password:       pointer.ToString("qwertyuio"),
				Phone:          pointer.ToString("9009009090"),
				TelegramChatId: pointer.ToString("abracadabra"),
			},
			expectedError: fmt.Errorf("Field validation for 'TelegramChatID' failed on the 'numeric' tag"),
		},
		{
			name: "valid_full_user",
			request: &v1.AddRequest{
				DisplayName:    "John Doe",
				Type:           "admin",
				Phone:          pointer.ToString("+79009009090"),
				Email:          pointer.ToString("john.doe@mail.com"),
				TelegramChatId: pointer.ToString("9009090"),
				Password:       pointer.ToString("qwertyuio"),
			},
			expectedResult: &UserAddDTO{
				DisplayName:    "John Doe",
				Type:           "admin",
				Phone:          "9009009090",
				Email:          "john.doe@mail.com",
				TelegramChatID: "9009090",
				Password:       "qwertyuio",
				PasswordHash:   defaultPasswordHash,
				DeactivatedAt:  nil,
			},
		},
		{
			name: "valid_user_with_email_only",
			request: &v1.AddRequest{
				DisplayName: "John Doe",
				Type:        "admin",
				Email:       pointer.ToString("john.doe@mail.com"),
				Password:    pointer.ToString("qwertyuio"),
			},
			expectedResult: &UserAddDTO{
				DisplayName:  "John Doe",
				Type:         "admin",
				Email:        "john.doe@mail.com",
				Password:     "qwertyuio",
				PasswordHash: defaultPasswordHash,
			},
		},
		{
			name: "valid_user_with_phone_only",
			request: &v1.AddRequest{
				DisplayName: "John Doe",
				Type:        "admin",
				Phone:       pointer.ToString("+79009009090"),
				Password:    pointer.ToString("qwertyuio"),
			},
			expectedResult: &UserAddDTO{
				DisplayName:  "John Doe",
				Type:         "admin",
				Phone:        "9009009090",
				Password:     "qwertyuio",
				PasswordHash: defaultPasswordHash,
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			userUsecase := NewUserUsecase(nil, nil, nil)
			actualResult, actualError := userUsecase.MakeUserAddDTO(testCase.request)
			if testCase.expectedError != nil {
				require.Error(t, actualError)
				require.Nil(t, actualResult)
				require.Truef(t, strings.Contains(actualError.Error(), testCase.expectedError.Error()),
					fmt.Sprintf("want:\n%v\ngot:\n%v", testCase.expectedError, actualError))
			} else {
				require.NoError(t, actualError)
				if actualResult.PasswordHash != "" {
					actualResult.PasswordHash = defaultPasswordHash
				}
				require.Equal(t, testCase.expectedResult, actualResult)
			}
		})
	}
}
