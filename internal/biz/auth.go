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
	"auth/internal/pkg/strings"
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

	codeExpiredInterval   = 1 * time.Minute
	resetPasswordInterval = 24 * time.Hour // 1 Day

	loginAgainCount             = 10
	loginAgainInterval          = time.Minute
	loginByCodeAgainCount       = 3
	loginByCodeAgainInterval    = time.Minute
	resetPasswordAgainCount     = 1
	resetPasswordAgainInterval  = 10 * time.Minute
	newPasswordAgainCount       = 1
	newPasswordAgainInterval    = time.Minute
	changePasswordAgainCount    = 3
	changePasswordAgainInterval = time.Minute
	generateCodeAgainCount      = 1
	generateCodeAgainInterval   = 2 * time.Minute

	metricPrefixAuth = `biz.auth`
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

type AuthUsecase struct {
	userRepo     UserRepo
	sessionRepo  SessionRepo
	codeRepo     CodeRepo
	historyRepo  HistoryRepo
	notification clients.Notifications
	metric       metrics.Metrics
	logger       logger.Logger
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
		logger:       logger.NewHelper(logs, `ts`, log.DefaultTimestamp, `scope`, `biz/auth`),
	}
}

func (a *AuthUsecase) postProcess(ctx context.Context, method string, err error) {
	ignoreErrors := []error{
		ErrWrongPassword,
		ErrWrongCode,
		ErrWrongResetHash,
		ErrWrongOldPassword,
		ErrLoginTooOften,
		ErrLoginByCodeTooOften,
		ErrResetPasswordTooOften,
		ErrNewPasswordTooOften,
		ErrChangePasswordTooOften,
		ErrGenerateCodeTooOften,
	}
	if err != nil {
		for _, ignoreError := range ignoreErrors {
			if errors.Is(err, ignoreError) {
				err = nil
				break
			}
		}
	}
	if err != nil {
		a.logger.WithContext(ctx).Errorf(`biz auth method "%s" failed: %v`, method, err)
		a.metric.Increment(strings.Metric(metricPrefixAuth, method, `failure`))
	} else {
		a.metric.Increment(strings.Metric(metricPrefixAuth, method, `success`))
	}
}

func (a *AuthUsecase) Check(ctx context.Context, token string) (*CheckDTO, error) {
	method := `check`
	defer a.metric.NewTiming().Send(strings.Metric(metricPrefixAuth, method, `timings`))
	var err error
	defer func() { a.postProcess(ctx, method, err) }()

	if token == "" {
		err = errors.New(`token is empty`)
		return nil, err
	}

	session, err := a.sessionRepo.FindByToken(ctx, token)
	if err != nil {
		return nil, err
	}
	user, err := a.userRepo.FindByID(ctx, session.UserID)
	if err != nil {
		return nil, err
	}
	return &CheckDTO{
		UserID:           user.ID,
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
	method := `login`
	defer a.metric.NewTiming().Send(strings.Metric(metricPrefixAuth, method, `timings`))
	var err error
	defer func() { a.postProcess(ctx, method, err) }()

	if dto == nil {
		err = errors.New("loginDTO is empty")
		return "", nil, err
	}

	foundedUser, err := a.userByUsername(ctx, dto.Username)
	if err != nil {
		return "", nil, err
	}
	if foundedUser.PasswordHash == nil {
		err = errors.New(`user password was not set, use login by code`)
		return "", nil, err
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
		err = ErrLoginTooOften
		return "", nil, err
	}
	match := secrets.MustCompareSourceAndHash(dto.Password, *foundedUser.PasswordHash)
	event := historyModel(foundedUser.ID, schema.EventLoginOk, dto.Stats)
	defer func() {
		_, err = a.historyRepo.Create(ctx, event)
	}()
	if !match {
		event.Event = schema.EventLoginFailed
		err = ErrWrongPassword
		return "", nil, err
	}
	session := sessionModelByLogin(foundedUser.ID, *foundedUser.PasswordHash, dto)
	_, err = a.sessionRepo.Create(ctx, session)
	if err != nil {
		return "", nil, err
	}
	return session.Token, &session.ExpiredAt, nil
}

func (a *AuthUsecase) LoginByCode(ctx context.Context, dto *LoginByCodeDTO) (string, *time.Time, error) {
	method := `loginByCode`
	defer a.metric.NewTiming().Send(strings.Metric(metricPrefixAuth, method, `timings`))
	var err error
	defer func() { a.postProcess(ctx, method, err) }()

	if dto == nil {
		err = errors.New("loginByCodeDTO is empty")
		return "", nil, err
	}

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
		err = ErrLoginByCodeTooOften
		return "", nil, err
	}
	var actualCode *ent.Code
	actualCode, err = a.codeRepo.FindForUser(ctx, foundedUser.ID)
	if ent.IsNotFound(err) {
		err = fmt.Errorf(`user by username %s does not have actual code`, dto.Username)
		return "", nil, err
	}
	if err != nil {
		return "", nil, err
	}
	match := actualCode.Content == dto.Code
	event := historyModel(foundedUser.ID, schema.EventLoginByCodeOk, dto.Stats)
	defer func() {
		_, err = a.historyRepo.Create(ctx, event)
	}()
	if !match {
		event.Event = schema.EventLoginByCodeFailed
		return "", nil, ErrWrongCode
	}
	session := sessionModelByLoginByCode(foundedUser.ID, actualCode.Content, dto)
	_, err = a.sessionRepo.Create(ctx, session)
	if err != nil {
		return "", nil, err
	}
	return session.Token, &session.ExpiredAt, nil
}

func (a *AuthUsecase) ResetPassword(ctx context.Context, dto *ResetPasswordDTO) error {
	method := `resetPassword`
	defer a.metric.NewTiming().Send(strings.Metric(metricPrefixAuth, method, `timings`))
	var err error
	defer func() { a.postProcess(ctx, method, err) }()

	if dto == nil {
		err = errors.New("resetPasswordDTO is empty")
		return err
	}

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
		err = ErrResetPasswordTooOften
		return err
	}

	event := historyModel(foundedUser.ID, schema.EventResetPasswordRequest, dto.Stats)
	defer func() {
		_, err = a.historyRepo.Create(ctx, event)
	}()
	hash := makeHash(passwordResetLength, dto.Username, makeRandomString(randomStringLength))
	foundedUser.PasswordReset = &hash
	foundedUser.PasswordResetExpiredAt = pointer.ToTime(time.Now().Add(resetPasswordInterval))
	_, err = a.userRepo.Update(ctx, foundedUser)
	if err != nil {
		return err
	}

	err = a.sendNotifyWithPasswordReset(ctx, foundedUser, dto.Username, hash)
	return err
}

func (a *AuthUsecase) NewPassword(ctx context.Context, dto *NewPasswordDTO) error {
	method := `newPassword`
	defer a.metric.NewTiming().Send(strings.Metric(metricPrefixAuth, method, `timings`))
	var err error
	defer func() { a.postProcess(ctx, method, err) }()

	if dto == nil {
		err = errors.New("newPasswordDTO is empty")
		return err
	}

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
		err = ErrNewPasswordTooOften
		return err
	}
	if foundedUser.PasswordReset == nil {
		err = errors.New(`user password reset hash is empty`)
		return err
	}
	if foundedUser.PasswordResetExpiredAt == nil {
		err = errors.New(`user password reset expiration time must be set`)
		return err
	}
	if foundedUser.PasswordResetExpiredAt.Before(time.Now()) {
		foundedUser.PasswordReset = nil
		foundedUser.PasswordResetExpiredAt = nil
		if _, err = a.userRepo.Update(ctx, foundedUser); err != nil {
			a.logger.WithContext(ctx).Warnf(`failed to save founded user: %v`, err)
		}
		err = errors.New(`password reset hash is expired, try reset password again`)
		return err
	}
	match := dto.PasswordResetHash == *foundedUser.PasswordReset
	event := historyModel(foundedUser.ID, schema.EventNewPasswordOk, dto.Stats)
	defer func() {
		_, err = a.historyRepo.Create(ctx, event)
	}()
	if !match {
		event.Event = schema.EventNewPasswordFailed
		err = ErrWrongResetHash
		return err
	}
	foundedUser.PasswordHash = pointer.ToString(secrets.MustMakeHash(dto.Password))
	foundedUser.PasswordReset = nil
	_, err = a.userRepo.Update(ctx, foundedUser)
	return err
}

func (a *AuthUsecase) ChangePassword(ctx context.Context, dto *ChangePasswordDTO) error {
	method := `changePassword`
	defer a.metric.NewTiming().Send(strings.Metric(metricPrefixAuth, method, `timings`))
	var err error
	defer func() { a.postProcess(ctx, method, err) }()

	if dto == nil {
		err = errors.New("changePasswordDTO is empty")
		return err
	}

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
		err = ErrChangePasswordTooOften
		return err
	}
	match := secrets.MustCompareSourceAndHash(dto.OldPassword, *foundedUser.PasswordHash)
	event := historyModel(foundedUser.ID, schema.EventChangePasswordOk, dto.Stats)
	defer func() {
		_, err = a.historyRepo.Create(ctx, event)
	}()
	if !match {
		event.Event = schema.EventChangePasswordFailed
		err = ErrWrongOldPassword
		return err
	}
	foundedUser.PasswordHash = pointer.ToString(secrets.MustMakeHash(dto.NewPassword))
	_, err = a.userRepo.Update(ctx, foundedUser)
	return err
}

func (a *AuthUsecase) GenerateCode(ctx context.Context, dto *GenerateCodeDTO) error {
	method := `generateCode`
	defer a.metric.NewTiming().Send(strings.Metric(metricPrefixAuth, method, `timings`))
	var err error
	defer func() { a.postProcess(ctx, method, err) }()

	if dto == nil {
		err = errors.New("generateCodeDTO is empty")
		return err
	}

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
		err = ErrGenerateCodeTooOften
		return err
	}

	event := historyModel(foundedUser.ID, schema.EventGenerateCodeRequest, dto.Stats)
	defer func() {
		_, err = a.historyRepo.Create(ctx, event)
	}()

	code := mustMakeCode(codeLength)
	_, err = a.codeRepo.Create(ctx, codeModel(foundedUser.ID, code, codeExpiredInterval))
	if err != nil {
		return err
	}

	err = a.sendNotifyWithAuthCode(ctx, foundedUser, dto.Username, code)
	return err
}

func (a *AuthUsecase) History(ctx context.Context, dto *HistoryDTO) ([]*ent.History, error) {
	method := `history`
	defer a.metric.NewTiming().Send(strings.Metric(metricPrefixAuth, method, `timings`))
	var err error
	defer func() { a.postProcess(ctx, method, err) }()

	if dto == nil {
		err = errors.New("historyDTO is empty")
		return nil, err
	}

	var histories []*ent.History
	histories, err = a.historyRepo.FindUserEvents(ctx, dto.UserID, dto.Limit, dto.Offset)
	return histories, err
}
