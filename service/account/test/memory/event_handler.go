package memory

import (
	"context"
	"encoding/json"

	"github.com/charmingruby/doris/service/account/internal/access/core/event"
)

type EventHandler struct {
	Pub Publisher
}

func NewEventHandler(pub Publisher) *EventHandler {
	return &EventHandler{
		Pub: pub,
	}
}

func (h *EventHandler) DispatchSendOTPNotification(ctx context.Context, event event.SendOTPNotification) error {
	if !h.Pub.IsHealthy {
		return ErrUnhealthyDatasource
	}

	msg, err := json.Marshal(event)
	if err != nil {
		return err
	}

	return h.Pub.Publish(ctx, "notifications.otp.send", msg)
}

func (h *EventHandler) DispatchAPIKeyDelegated(ctx context.Context, event event.APIKeyDelegated) error {
	if !h.Pub.IsHealthy {
		return ErrUnhealthyDatasource
	}

	msg, err := json.Marshal(event)
	if err != nil {
		return err
	}

	return h.Pub.Publish(ctx, "api-keys.delegated", msg)
}

func (h *EventHandler) DispatchAPIKeyActivated(ctx context.Context, event event.APIKeyActivated) error {
	if !h.Pub.IsHealthy {
		return ErrUnhealthyDatasource
	}

	msg, err := json.Marshal(event)
	if err != nil {
		return err
	}

	return h.Pub.Publish(ctx, "api-keys.activated", msg)
}
