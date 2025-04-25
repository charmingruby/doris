package event

import (
	"context"
)

func (h *Handler) receiveAPIKeyActivationSubscription(ctx context.Context) error {
	return h.sub.Subscribe(ctx, h.topics[apiKeyActivationIdentifier], func(message []byte) error {
		h.log.Debug("received message", "topic", h.topics[apiKeyActivationIdentifier], "message", string(message))

		return nil
	})
}
