package event

import (
	"context"

	"github.com/charmingruby/doris/lib/delivery/messaging"
	"github.com/charmingruby/doris/lib/instrumentation"
	"github.com/charmingruby/doris/service/notification/internal/notification/core/service"
)

const (
	sendOTPNotificationIdentifier = iota
)

type Handler struct {
	sub    messaging.Subscriber
	logger *instrumentation.Logger
	topics map[int]string
	svc    *service.Service
}

type TopicInput struct {
	SendOTPNotification string
}

func NewHandler(logger *instrumentation.Logger, sub messaging.Subscriber, svc *service.Service, in TopicInput) *Handler {
	topics := make(map[int]string, 1)

	topics[sendOTPNotificationIdentifier] = in.SendOTPNotification

	return &Handler{
		logger: logger,
		sub:    sub,
		topics: topics,
		svc:    svc,
	}
}

func (h *Handler) Subscribe() error {
	ctx := context.Background()

	if err := h.onSendOTPNotification(ctx); err != nil {
		return err
	}

	return nil
}
