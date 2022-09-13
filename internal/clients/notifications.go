package clients

//go:generate mockgen -source ./${GOFILE} -destination ./notifications_mock_test.go -package ${GOPACKAGE}

import (
	"context"
	"fmt"
	"time"

	v1 "auth/api/notification/v1"
	"auth/internal/pkg/logger"
	"auth/internal/pkg/metrics"
	"auth/internal/pkg/sanitize"

	"github.com/go-kratos/kratos/v2/log"
)

const (
	metricNotificationsEnqueueMailTimings     = `clients.notifications.enqueueMail.timings`
	metricNotificationsEnqueueMailSuccess     = `clients.notifications.enqueueMail.success`
	metricNotificationsEnqueueMailFailure     = `clients.notifications.enqueueMail.failure`
	metricNotificationsEnqueueTelegramTimings = `clients.notifications.enqueueTelegram.timings`
	metricNotificationsEnqueueTelegramSuccess = `clients.notifications.enqueueTelegram.success`
	metricNotificationsEnqueueTelegramFailure = `clients.notifications.enqueueTelegram.failure`
	metricNotificationsEnqueueSMSTimings      = `clients.notifications.enqueueTelegram.timings`
	metricNotificationsEnqueueSMSSuccess      = `clients.notifications.enqueueTelegram.success`
	metricNotificationsEnqueueSMSFailure      = `clients.notifications.enqueueTelegram.failure`
)

type NotificationsClient struct {
	client v1.NotificationClient
	metric metrics.Metrics
	logger *log.Helper
}

type Notifications interface {
	EnqueueMailWithHTML(ctx context.Context, to, subject, body string) (int64, error)
	EnqueueTelegramWithMarkdown(ctx context.Context, chatID, text string) (int64, error)
	EnqueueSMS(ctx context.Context, phone, text string) (int64, error)
}

func NewNotifications(
	ctx context.Context,
	endpoint string,
	timeout time.Duration,
	metric metrics.Metrics,
	logs log.Logger,
) (*NotificationsClient, error) {
	client, err := Default(ctx, endpoint, timeout)
	if err != nil {
		return nil, err
	}
	return &NotificationsClient{
		client: v1.NewNotificationClient(client),
		metric: metric,
		logger: logger.NewHelper(logs, "ts", log.DefaultTimestamp, "scope", "clients/notifications"),
	}, nil
}

func defaultSendRequest(setter func(*v1.SendRequest)) *v1.SendRequest {
	s := &v1.SendRequest{
		Type:      v1.Type_sms,
		Payload:   map[string]string{},
		PlannedAt: nil,
		Ttl:       300,
		SenderId:  0,
	}
	setter(s)
	return s
}

func (n *NotificationsClient) EnqueueMailWithHTML(ctx context.Context, to, subject, body string) (int64, error) {
	defer n.metric.NewTiming().Send(metricNotificationsEnqueueMailTimings)
	var err error
	defer func() {
		if err != nil {
			n.logger.WithContext(ctx).Errorf(`failed to enqueue mail notification: %v`, err)
			n.metric.Increment(metricNotificationsEnqueueMailFailure)
		} else {
			n.metric.Increment(metricNotificationsEnqueueMailSuccess)
		}
	}()

	res, err := n.client.Enqueue(
		ctx,
		defaultSendRequest(
			func(s *v1.SendRequest) {
				s.Type = v1.Type_email
				s.Payload = map[string]string{
					"to":      to,
					"subject": subject,
					"body":    body,
					"is_html": fmt.Sprintf("%t", true),
				}
			},
		),
	)
	var notificationID int64
	if res != nil {
		notificationID = res.Id
	}
	return notificationID, err
}

func (n *NotificationsClient) EnqueueTelegramWithMarkdown(ctx context.Context, chatID, text string) (int64, error) {
	defer n.metric.NewTiming().Send(metricNotificationsEnqueueTelegramTimings)
	var err error
	defer func() {
		if err != nil {
			n.logger.WithContext(ctx).Errorf(`failed to enqueue telegram notification: %v`, err)
			n.metric.Increment(metricNotificationsEnqueueTelegramFailure)
		} else {
			n.metric.Increment(metricNotificationsEnqueueTelegramSuccess)
		}
	}()
	res, err := n.client.Enqueue(
		ctx,
		defaultSendRequest(
			func(s *v1.SendRequest) {
				s.Type = v1.Type_telegram
				s.Payload = map[string]string{
					// Attributes based on https://core.telegram.org/bots/api#sendmessage
					"chat_id":         chatID,     // Unique identifier for the target chat or username of the target channel (in the format @channelusername)
					"text":            text,       // Text of the message to be sent, 1-4096 characters after entities parsing
					"parse_mode":      "markdown", // Mode for parsing entities in the message text. See formatting options (https://core.telegram.org/bots/api#formatting-options) for more details.
					"protect_content": "true",     // Protects the contents of the sent message from forwarding and saving
				}
			},
		),
	)
	var notificationID int64
	if res != nil {
		notificationID = res.Id
	}
	return notificationID, err
}

func (n *NotificationsClient) EnqueueSMS(ctx context.Context, phone, text string) (int64, error) {
	defer n.metric.NewTiming().Send(metricNotificationsEnqueueSMSTimings)
	var err error
	defer func() {
		if err != nil {
			n.logger.WithContext(ctx).Errorf(`failed to enqueue sms notification: %v`, err)
			n.metric.Increment(metricNotificationsEnqueueSMSFailure)
		} else {
			n.metric.Increment(metricNotificationsEnqueueSMSSuccess)
		}
	}()
	sanitizedPhone := sanitize.PhoneWithCountryCode(phone)
	res, err := n.client.Enqueue(
		ctx,
		defaultSendRequest(
			func(s *v1.SendRequest) {
				s.Type = v1.Type_sms
				s.Payload = map[string]string{
					"phone": sanitizedPhone, // Phone number in format 79009009090
					"text":  text,           // Text of SMS message, limit of 160 symbols for latin and 70 symbols for cyrillic
					"split": "false",        // Split to few messages if text length exceeds limit
				}
			},
		),
	)
	var notificationID int64
	if res != nil {
		notificationID = res.Id
	}
	return notificationID, err
}
