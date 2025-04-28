package event

import "context"

type Handler interface {
	SendOTPNotification(ctx context.Context, event *SendOTPNotificationMessage) error
}
