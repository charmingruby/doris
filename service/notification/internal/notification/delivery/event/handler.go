package event

import (
	"context"
	"time"

	"github.com/charmingruby/doris/lib/delivery/messaging"
	"github.com/charmingruby/doris/lib/instrumentation"
	"github.com/charmingruby/doris/service/notification/internal/notification/core/service"
)

const (
	receiveNotificationIdentifier = iota
)

type Handler struct {
	sub    messaging.Subscriber
	logger *instrumentation.Logger
	topics map[int]string
	svc    *service.Service
}

type TopicInput struct {
	ReceiveNotificationTopic string
}

func NewHandler(logger *instrumentation.Logger, sub messaging.Subscriber, svc *service.Service, in TopicInput) *Handler {
	topics := make(map[int]string, 1)

	topics[receiveNotificationIdentifier] = in.ReceiveNotificationTopic

	return &Handler{
		logger: logger,
		sub:    sub,
		topics: topics,
		svc:    svc,
	}
}

func (h *Handler) Subscribe() error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := h.receiveNotificationSendIntent(ctx); err != nil {
		return err
	}

	return nil
}
