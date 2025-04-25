package event

import "context"

type Handler interface {
	SendAPIKeyActivation(ctx context.Context, event *APIKeyActivation) error
}
