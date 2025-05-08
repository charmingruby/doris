package event

import "context"

type Handler interface {
	DispatchSendOTPNotification(ctx context.Context, event SendOTPNotification) error
	DispatchAPIKeyDelegated(ctx context.Context, event APIKeyDelegated) error
}
