package event

import "context"

type Handler interface {
	DispatchSendOTPNotification(ctx context.Context, event SendOTPNotificationMessage) error
	DispatchAPIKeyDelegation(ctx context.Context, event APIKeyDelegationMessage) error
}
