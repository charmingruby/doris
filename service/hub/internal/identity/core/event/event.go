package event

import "context"

type Handler interface {
	SendOTP(ctx context.Context, event *OTP) error
}
