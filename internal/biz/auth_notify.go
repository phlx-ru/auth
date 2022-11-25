package biz

import (
	"context"
	"errors"
	"fmt"

	"auth/ent"
	"auth/internal/pkg/template"
	"auth/internal/pkg/texts"
)

func (a *AuthUsecase) sendNotifyWithPasswordReset(ctx context.Context, user *ent.User, username, hash string) error {
	interpolation := map[string]any{
		"hash": hash,
	}
	telegram := func() error {
		text := template.MustInterpolate(texts.ResetPasswordTelegramBody, interpolation)
		_, err := a.notification.EnqueueTelegramWithMarkdown(ctx, *user.TelegramChatID, text)
		return err
	}
	phone := func() error {
		text := template.MustInterpolate(texts.ResetPasswordSMSBody, interpolation)
		_, err := a.notification.EnqueueSMS(ctx, *user.Phone, text)
		return err
	}
	email := func() error {
		subject := texts.ResetPasswordEmailSubject
		body := template.MustInterpolate(texts.ResetPasswordEmailBody, interpolation)
		_, err := a.notification.EnqueueMailWithHTML(ctx, *user.Email, subject, body)
		return err
	}
	if user.Phone != nil && username == *user.Phone {
		return phone()
	}
	if user.Email != nil && username == *user.Email {
		return email()
	}
	if user.TelegramChatID != nil {
		return telegram()
	}
	a.logger.WithContext(ctx).Errorf(`notify password reset failed: user %d by username %s`, user.ID, username)
	return errors.New(`user has not set telegram, phone or email, notify with password reset failed`)
}

func (a *AuthUsecase) sendNotifyWithAuthCode(ctx context.Context, user *ent.User, username, code string) error {
	var err error
	m := int(codeExpiredInterval.Minutes())
	interpolation := map[string]any{
		"code":    code,
		"minutes": fmt.Sprintf("%d %s", m, texts.Plural(m, `минуту`, `минуты`, `минут`)),
	}
	telegram := func() error {
		text := template.MustInterpolate(texts.AuthCodeTelegramBody, interpolation)
		_, err = a.notification.EnqueueTelegramWithMarkdown(ctx, *user.TelegramChatID, text)
		return err
	}
	phone := func() error {
		text := template.MustInterpolate(texts.AuthCodeSMSBody, interpolation)
		_, err = a.notification.EnqueueSMS(ctx, *user.Phone, text)
		return err
	}
	email := func() error {
		subject := texts.AuthCodeEmailSubject
		body := template.MustInterpolate(texts.AuthCodeEmailBody, interpolation)
		_, err = a.notification.EnqueueMailWithHTML(ctx, *user.Email, subject, body)
		return err
	}
	if user.Phone != nil && username == *user.Phone {
		return phone()
	}
	if user.Email != nil && username == *user.Email {
		return email()
	}
	if user.TelegramChatID != nil {
		return telegram()
	}
	a.logger.WithContext(ctx).Errorf(`notify auth code failed: user %d by username %s`, user.ID, username)
	return errors.New(`user has not set telegram, phone or email, notify with code failed`)
}
