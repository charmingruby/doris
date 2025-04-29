package memory

import (
	"context"
	"encoding/json"

	"github.com/charmingruby/doris/service/identity/internal/access/core/event"
)

type EventHandler struct {
	Pub Publisher
}

func NewEventHandler(pub Publisher) *EventHandler {
	return &EventHandler{
		Pub: pub,
	}
}

func (h *EventHandler) SendOTPNotification(ctx context.Context, event *event.SendOTPNotificationMessage) error {
	if !h.Pub.IsHealthy {
		return ErrUnhealthyDatasource
	}

	msg, err := json.Marshal(event)
	if err != nil {
		return err
	}

	return h.Pub.Publish(ctx, "send_otp_notification", msg)
}
