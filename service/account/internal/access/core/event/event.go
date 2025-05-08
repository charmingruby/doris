package event

import "context"

type Handler interface {
	DispatchSendOTPNotification(ctx context.Context, event SendOTPNotificationMessage) error
	DispatchAPIKeyDelegated(ctx context.Context, event APIKeyDelegatedMessage) error
}
