package event

import "context"

type Handler interface {
	DispatchSendOTPNotification(ctx context.Context, message SendOTPNotification) error
	DispatchAPIKeyDelegated(ctx context.Context, message APIKeyDelegated) error
	DispatchAPIKeyActivated(ctx context.Context, message APIKeyActivated) error
}
