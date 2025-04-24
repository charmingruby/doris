package memory

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/charmingruby/doris/service/gateway/internal/identity/core/event"
)

type EventHandler struct {
	Pub Publisher
}

func NewEventHandler(pub Publisher) *EventHandler {
	return &EventHandler{
		Pub: pub,
	}
}

func (h *EventHandler) SendAPIKeyActivationCode(ctx context.Context, event *event.APIKeyActivation) error {
	if !h.Pub.IsHealthy {
		return errors.New("publisher is not healthy")
	}

	msg, err := json.Marshal(event)
	if err != nil {
		return err
	}

	return h.Pub.Publish(ctx, "api_key_activation", msg)
}
