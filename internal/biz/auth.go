package biz

import (
	"context"
	"errors"
	"fmt"
	"time"

	"auth/ent"
	"auth/ent/schema"
	"auth/internal/clients"
	"auth/internal/pkg/logger"
	"auth/internal/pkg/metrics"
	"auth/internal/pkg/secrets"
	"auth/internal/pkg/template"
	"auth/internal/pkg/texts"
	"github.com/AlekSi/pointer"
	"github.com/go-kratos/kratos/v2/log"
)

const (
	sessionExpiredIntervalLong  = time.Hour * 24 * 7 // one week
	sessionExpiredIntervalShort = time.Hour

	codeLength          = 4
	passwordResetLength = 6
	randomStringLength  = 16

	loginAgainCount             = 10
	loginAgainInterval          = time.Minute
	loginByCodeAgainCount       = 3
	loginByCodeAgainInterval    = time.Minute
	resetPasswordAgainCount     = 1
	resetPasswordAgainInterval  = time.Minute
	newPasswordAgainCount       = 1
	newPasswordAgainInterval    = time.Minute
	changePasswordAgainCount    = 3
	changePasswordAgainInterval = time.Minute
	generateCodeAgainCount      = 1
	generateCodeAgainInterval   = time.Minute

	metricAuthCheckSuccess = `biz.auth.check.success`
	metricAuthCheckFailure = `biz.auth.check.failure`
	metricAuthCheckTimings = `biz.auth.check.timings`

	metricAuthLoginSuccess = `biz.auth.login.success`
	metricAuthLoginFailure = `biz.auth.login.failure`
	metricAuthLoginTimings = `biz.auth.login.timings`

	metricAuthLoginByCodeSuccess = `biz.auth.loginByCode.success`
	metricAuthLoginByCodeFailure = `biz.auth.loginByCode.failure`
	metricAuthLoginByCodeTimings = `biz.auth.loginByCode.timings`

	metricAuthResetPasswordSuccess = `biz.auth.resetPassword.success`
	metricAuthResetPasswordFailure = `biz.auth.resetPassword.failure`
	metricAuthResetPasswordTimings = `biz.auth.resetPassword.timings`

	metricAuthNewPasswordSuccess = `biz.auth.newPassword.success`
	metricAuthNewPasswordFailure = `biz.auth.newPassword.failure`
	metricAuthNewPasswordTimings = `biz.auth.newPassword.timings`

	metricAuthChangePasswordSuccess = `biz.auth.changePassword.success`
	metricAuthChangePasswordFailure = `biz.auth.changePassword.failure`
	metricAuthChangePasswordTimings = `biz.auth.changePassword.timings`

	metricAuthGenerateCodeSuccess = `biz.auth.generateCode.success`
	metricAuthGenerateCodeFailure = `biz.auth.generateCode.failure`
	metricAuthGenerateCodeTimings = `biz.auth.generateCode.timings`

	metricAuthHistorySuccess = `biz.auth.history.success`
	metricAuthHistoryFailure = `biz.auth.history.failure`
	metricAuthHistoryTimings = `biz.auth.history.timings`
)

var (
	ErrWrongPassword          = errors.New(texts.WrongPassword)
	ErrWrongCode              = errors.New(texts.WrongCode)
	ErrWrongResetHash         = errors.New(texts.WrongResetHash)
	ErrWrongOldPassword       = errors.New(texts.WrongOldPassword)
	ErrLoginTooOften          = errors.New(texts.LoginTooOften)
	ErrLoginByCodeTooOften    = errors.New(texts.LoginByCodeTooOften)
	ErrResetPasswordTooOften  = errors.New(texts.ResetPasswordTooOften)
	ErrNewPasswordTooOften    = errors.New(texts.NewPasswordTooOften)
	ErrChangePasswordTooOften = errors.New(texts.ChangePasswordTooOften)
	ErrGenerateCodeTooOften   = errors.New(texts.GenerateCodeTooOften)
)

type SessionRepo interface {
	Save(context.Context, *ent.Session) (*ent.Session, error)
	FindByToken(ctx context.Context, token string) (*ent.Session, error)
}

type CodeRepo interface {
	Save(context.Context, *ent.Code) (*ent.Code, error)
	FindForUser(ctx context.Context, userID int) (*ent.Code, error)
}

type HistoryRepo interface {
	Save(context.Context, *ent.History) (*ent.History, error)
	FindUserEvents(ctx context.Context, userID, limit, offset int) ([]*ent.History, error)
	FindLastUserEvents(
		ctx context.Context,
		userID int,
		types []string,
		interval time.Duration,
	) ([]*ent.History, error)
}

type AuthUsecase struct {
	userRepo     UserRepo
	sessionRepo  SessionRepo
	codeRepo     CodeRepo
	historyRepo  HistoryRepo
	notification clients.Notifications
	metric       metrics.Metrics
	logs         logger.Logger
}

func NewAuthUsecase(
	userRepo UserRepo,
	sessionRepo SessionRepo,
	codeRepo CodeRepo,
	historyRepo HistoryRepo,
	notification clients.Notifications,
	metric metrics.Metrics,
	logs log.Logger,
) *AuthUsecase {
	return &AuthUsecase{
		userRepo:     userRepo,
		sessionRepo:  sessionRepo,
		codeRepo:     codeRepo,
		historyRepo:  historyRepo,
		notification: notification,
		metric:       metric,
		logs:         logger.NewHelper(logs, `ts`, log.DefaultTimestamp, `scope`, `biz/auth`),
	}
}

func (a *AuthUsecase) Check(ctx context.Context, token string) (*CheckDTO, error) {
	defer a.metric.NewTiming().Send(metricAuthCheckTimings)
	if token == "" {
		return nil, errors.New(`token is empty`)
	}
	var err error
	defer func() {
		if err != nil {
			a.logs.WithContext(ctx).Errorf(`failed to check auth: %v`, err)
			a.metric.Increment(metricAuthCheckFailure)
		} else {
			a.metric.Increment(metricAuthCheckSuccess)
		}
	}()
	session, err := a.sessionRepo.FindByToken(ctx, token)
	if err != nil {
		return nil, err
	}
	user, err := a.userRepo.FindByID(ctx, session.UserID)
	if err != nil {
		return nil, err
	}
	return &CheckDTO{
		UserType:         user.Type,
		DisplayName:      user.DisplayName,
		Email:            user.Email,
		Phone:            user.Phone,
		SessionUntil:     session.ExpiredAt,
		SessionIP:        session.IP,
		SessionUserAgent: session.UserAgent,
		SessionDeviceID:  session.DeviceID,
	}, nil
}

func (a *AuthUsecase) Login(ctx context.Context, dto *LoginDTO) (string, *time.Time, error) {
	defer a.metric.NewTiming().Send(metricAuthLoginTimings)
	if dto == nil {
		return "", nil, errors.New("loginDTO is empty")
	}
	var err error
	defer func() {
		if err != nil && err != ErrWrongPassword && err != ErrLoginTooOften {
			a.logs.WithContext(ctx).Errorf(`failed to login: %v`, err)
			a.metric.Increment(metricAuthLoginFailure)
		} else {
			a.metric.Increment(metricAuthLoginSuccess)
		}
	}()
	foundedUser, err := a.userByUsername(ctx, dto.Username)
	if err != nil {
		return "", nil, err
	}
	if foundedUser.PasswordHash == nil {
		return "", nil, errors.New(`user password was not set, use login by code`)
	}
	events, err := a.historyRepo.FindLastUserEvents(
		ctx,
		foundedUser.ID,
		[]string{schema.EventLoginFailed},
		loginAgainInterval,
	)
	if err != nil {
		return "", nil, err
	}
	if len(events) > loginAgainCount {
		return "", nil, ErrLoginTooOften
	}
	match := secrets.MustCompareSourceAndHash(dto.Password, *foundedUser.PasswordHash)
	event := historyModel(foundedUser.ID, schema.EventLoginOk, dto.Stats)
	defer func() {
		_, err := a.historyRepo.Save(ctx, event)
		if err != nil {
			a.logs.WithContext(ctx).Errorf(`failed to save history: %v`, err)
		}
	}()
	if !match {
		event.Event = schema.EventLoginFailed
		return "", nil, ErrWrongPassword
	}
	session := sessionModelByLogin(foundedUser.ID, *foundedUser.PasswordHash, dto)
	_, err = a.sessionRepo.Save(ctx, session)
	if err != nil {
		return "", nil, err
	}
	return session.Token, &session.ExpiredAt, nil
}

func (a *AuthUsecase) LoginByCode(ctx context.Context, dto *LoginByCodeDTO) (string, *time.Time, error) {
	defer a.metric.NewTiming().Send(metricAuthLoginByCodeTimings)
	if dto == nil {
		return "", nil, errors.New("loginByCodeDTO is empty")
	}
	var err error
	defer func() {
		if err != nil && err != ErrWrongCode && err != ErrLoginByCodeTooOften {
			a.logs.WithContext(ctx).Errorf(`failed to login by code: %v`, err)
			a.metric.Increment(metricAuthLoginByCodeFailure)
		} else {
			a.metric.Increment(metricAuthLoginByCodeSuccess)
		}
	}()
	foundedUser, err := a.userByUsername(ctx, dto.Username)
	if err != nil {
		return "", nil, err
	}
	events, err := a.historyRepo.FindLastUserEvents(
		ctx,
		foundedUser.ID,
		[]string{schema.EventLoginByCodeFailed},
		loginByCodeAgainInterval,
	)
	if err != nil {
		return "", nil, err
	}
	if len(events) > loginByCodeAgainCount {
		return "", nil, ErrLoginByCodeTooOften
	}
	actualCode, err := a.codeRepo.FindForUser(ctx, foundedUser.ID)
	if ent.IsNotFound(err) {
		return "", nil, fmt.Errorf(`user by username %s does not have actual code`, dto.Username)
	}
	if err != nil {
		return "", nil, err
	}
	match := actualCode.Content == dto.Code
	event := historyModel(foundedUser.ID, schema.EventLoginByCodeOk, dto.Stats)
	defer func() {
		if _, err := a.historyRepo.Save(ctx, event); err != nil {
			a.logs.WithContext(ctx).Errorf(`failed to save history: %v`, err)
		}
	}()
	if !match {
		event.Event = schema.EventLoginByCodeFailed
		return "", nil, ErrWrongCode
	}
	session := sessionModelByLoginByCode(foundedUser.ID, actualCode.Content, dto)
	_, err = a.sessionRepo.Save(ctx, session)
	if err != nil {
		return "", nil, err
	}
	return session.Token, &session.ExpiredAt, nil
}

func (a *AuthUsecase) ResetPassword(ctx context.Context, dto *ResetPasswordDTO) error {
	defer a.metric.NewTiming().Send(metricAuthResetPasswordTimings)
	if dto == nil {
		return errors.New("resetPasswordDTO is empty")
	}
	var err error
	defer func() {
		if err != nil && err != ErrResetPasswordTooOften {
			a.logs.WithContext(ctx).Errorf(`failed to reset password: %v`, err)
			a.metric.Increment(metricAuthResetPasswordFailure)
		} else {
			a.metric.Increment(metricAuthResetPasswordSuccess)
		}
	}()
	foundedUser, err := a.userByUsername(ctx, dto.Username)
	if err != nil {
		return err
	}
	events, err := a.historyRepo.FindLastUserEvents(
		ctx,
		foundedUser.ID,
		[]string{schema.EventResetPasswordRequest},
		resetPasswordAgainInterval,
	)
	if err != nil {
		return err
	}
	if len(events) >= resetPasswordAgainCount {
		return ErrResetPasswordTooOften
	}

	event := historyModel(foundedUser.ID, schema.EventResetPasswordRequest, dto.Stats)
	defer func() {
		_, err := a.historyRepo.Save(ctx, event)
		if err != nil {
			a.logs.WithContext(ctx).Errorf(`failed to save history: %v`, err)
		}
	}()
	hash := makeHash(passwordResetLength, dto.Username, makeRandomString(randomStringLength))
	foundedUser.PasswordReset = &hash
	_, err = a.userRepo.Update(ctx, foundedUser)
	if err != nil {
		return err
	}

	interpolation := map[string]any{
		"hash": hash,
	}
	if foundedUser.TelegramChatID != nil {
		text := template.MustInterpolate(texts.ResetPasswordTelegramBody, interpolation)
		_, err = a.notification.EnqueueTelegramWithMarkdown(ctx, *foundedUser.TelegramChatID, text)
	} else {
		subject := texts.ResetPasswordEmailSubject
		body := template.MustInterpolate(texts.ResetPasswordEmailBody, interpolation)
		_, err = a.notification.EnqueueMailWithHTML(ctx, *foundedUser.Email, subject, body)
	}
	return err
}

func (a *AuthUsecase) NewPassword(ctx context.Context, dto *NewPasswordDTO) error {
	defer a.metric.NewTiming().Send(metricAuthNewPasswordTimings)
	if dto == nil {
		return errors.New("newPasswordDTO is empty")
	}
	var err error
	defer func() {
		if err != nil && err != ErrNewPasswordTooOften && err != ErrWrongResetHash {
			a.logs.WithContext(ctx).Errorf(`failed to set new password: %v`, err)
			a.metric.Increment(metricAuthNewPasswordFailure)
		} else {
			a.metric.Increment(metricAuthNewPasswordSuccess)
		}
	}()
	foundedUser, err := a.userByUsername(ctx, dto.Username)
	if err != nil {
		return err
	}
	events, err := a.historyRepo.FindLastUserEvents(
		ctx,
		foundedUser.ID,
		[]string{schema.EventNewPasswordFailed},
		newPasswordAgainInterval,
	)
	if err != nil {
		return err
	}
	if len(events) >= newPasswordAgainCount {
		return ErrNewPasswordTooOften
	}
	if foundedUser.PasswordReset == nil {
		return errors.New(`user password reset hash is empty`)
	}
	match := dto.PasswordResetHash != *foundedUser.PasswordReset
	event := historyModel(foundedUser.ID, schema.EventNewPasswordOk, dto.Stats)
	defer func() {
		if _, err := a.historyRepo.Save(ctx, event); err != nil {
			a.logs.WithContext(ctx).Errorf(`failed to save history: %v`, err)
		}
	}()
	if !match {
		event.Event = schema.EventNewPasswordFailed
		return ErrWrongResetHash
	}
	foundedUser.PasswordHash = pointer.ToString(secrets.MustMakeHash(dto.Password))
	foundedUser.PasswordReset = nil
	_, err = a.userRepo.Update(ctx, foundedUser)
	return err
}

func (a *AuthUsecase) ChangePassword(ctx context.Context, dto *ChangePasswordDTO) error {
	defer a.metric.NewTiming().Send(metricAuthChangePasswordTimings)
	if dto == nil {
		return errors.New("changePasswordDTO is empty")
	}
	var err error
	defer func() {
		if err != nil && err != ErrChangePasswordTooOften && err != ErrWrongOldPassword {
			a.logs.WithContext(ctx).Errorf(`failed to change password: %v`, err)
			a.metric.Increment(metricAuthChangePasswordFailure)
		} else {
			a.metric.Increment(metricAuthChangePasswordSuccess)
		}
	}()
	foundedUser, err := a.userByUsername(ctx, dto.Username)
	if err != nil {
		return err
	}
	events, err := a.historyRepo.FindLastUserEvents(
		ctx,
		foundedUser.ID,
		[]string{schema.EventChangePasswordFailed},
		changePasswordAgainInterval,
	)
	if err != nil {
		return err
	}
	if len(events) >= changePasswordAgainCount {
		return ErrChangePasswordTooOften
	}
	match := secrets.MustCompareSourceAndHash(dto.OldPassword, *foundedUser.PasswordHash)
	event := historyModel(foundedUser.ID, schema.EventChangePasswordOk, dto.Stats)
	defer func() {
		if _, err := a.historyRepo.Save(ctx, event); err != nil {
			a.logs.WithContext(ctx).Errorf(`failed to save history: %v`, err)
		}
	}()
	if !match {
		event.Event = schema.EventChangePasswordFailed
		return ErrWrongOldPassword
	}
	foundedUser.PasswordHash = pointer.ToString(secrets.MustMakeHash(dto.NewPassword))
	_, err = a.userRepo.Update(ctx, foundedUser)
	return err
}

func (a *AuthUsecase) GenerateCode(ctx context.Context, dto *GenerateCodeDTO) error {
	defer a.metric.NewTiming().Send(metricAuthGenerateCodeTimings)
	if dto == nil {
		return errors.New("generateCodeDTO is empty")
	}
	var err error
	defer func() {
		if err != nil && err != ErrGenerateCodeTooOften {
			a.logs.WithContext(ctx).Errorf(`failed to generate code: %v`, err)
			a.metric.Increment(metricAuthGenerateCodeFailure)
		} else {
			a.metric.Increment(metricAuthGenerateCodeSuccess)
		}
	}()
	foundedUser, err := a.userByUsername(ctx, dto.Username)
	if err != nil {
		return err
	}
	events, err := a.historyRepo.FindLastUserEvents(
		ctx,
		foundedUser.ID,
		[]string{schema.EventGenerateCodeRequest},
		generateCodeAgainInterval,
	)
	if err != nil {
		return err
	}
	if len(events) >= generateCodeAgainCount {
		return ErrGenerateCodeTooOften
	}

	event := historyModel(foundedUser.ID, schema.EventGenerateCodeRequest, dto.Stats)
	defer func() {
		_, err := a.historyRepo.Save(ctx, event)
		if err != nil {
			a.logs.WithContext(ctx).Errorf(`failed to save history: %v`, err)
		}
	}()

	code := mustMakeCode(codeLength)
	_, err = a.codeRepo.Save(ctx, codeModel(foundedUser.ID, code))
	if err != nil {
		return err
	}

	interpolation := map[string]any{
		"code": code,
	}
	if foundedUser.TelegramChatID != nil {
		text := template.MustInterpolate(texts.AuthCodeTelegramBody, interpolation)
		_, err = a.notification.EnqueueTelegramWithMarkdown(ctx, *foundedUser.TelegramChatID, text)
	} else {
		subject := texts.AuthCodeEmailSubject
		body := template.MustInterpolate(texts.AuthCodeEmailBody, interpolation)
		_, err = a.notification.EnqueueMailWithHTML(ctx, *foundedUser.Email, subject, body)
	}
	return err
}

func (a *AuthUsecase) History(ctx context.Context, dto *HistoryDTO) ([]*ent.History, error) {
	defer a.metric.NewTiming().Send(metricAuthHistoryTimings)
	if dto == nil {
		return nil, errors.New("historyDTO is empty")
	}
	var err error
	defer func() {
		if err != nil {
			a.logs.WithContext(ctx).Errorf(`failed to get history: %v`, err)
			a.metric.Increment(metricAuthHistoryFailure)
		} else {
			a.metric.Increment(metricAuthHistorySuccess)
		}
	}()
	return a.historyRepo.FindUserEvents(ctx, dto.UserID, dto.Limit, dto.Offset)
}
