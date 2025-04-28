package memory

import (
	"context"
	"encoding/json"

	"github.com/charmingruby/doris/service/hub/internal/identity/core/event"
)

type EventHandler struct {
	Pub Publisher
}

func NewEventHandler(pub Publisher) *EventHandler {
	return &EventHandler{
		Pub: pub,
	}
}

func (h *EventHandler) SendOTP(ctx context.Context, event *event.OTP) error {
	if !h.Pub.IsHealthy {
		return ErrUnhealthyDatasource
	}

	msg, err := json.Marshal(event)
	if err != nil {
		return err
	}

	return h.Pub.Publish(ctx, "api_key_activation", msg)
}
