package event

import (
	"context"
	"time"

	"github.com/charmingruby/doris/lib/delivery/messaging"
	"github.com/charmingruby/doris/lib/instrumentation/logger"
	"github.com/charmingruby/doris/service/notification/internal/notification/core/service"
)

const (
	receiveNotificationIdentifier = iota
)

type Handler struct {
	sub    messaging.Subscriber
	log    *logger.Logger
	topics map[int]string
	svc    *service.Service
}

type TopicInput struct {
	ReceiveNotificationTopic string
}

func NewHandler(log *logger.Logger, sub messaging.Subscriber, svc *service.Service, in TopicInput) *Handler {
	topics := make(map[int]string, 1)

	topics[receiveNotificationIdentifier] = in.ReceiveNotificationTopic

	return &Handler{
		log:    log,
		sub:    sub,
		topics: topics,
		svc:    svc,
	}
}

func (h *Handler) Subscribe() error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := h.receiveNotification(ctx); err != nil {
		return err
	}

	return nil
}
