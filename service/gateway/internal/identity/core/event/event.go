package event

import "context"

type Handler interface {
	SendAPIKeyActivationCode(ctx context.Context, event *APIKeyActivation) error
}
