package event

import (
	"context"
	"time"

	"github.com/charmingruby/doris/lib/delivery/messaging"
	"github.com/charmingruby/doris/lib/instrumentation/logger"
)

const (
	apiKeyActivationIdentifier = iota
)

type Handler struct {
	sub messaging.Subscriber
	log *logger.Logger

	// identifier -> topic
	topics map[int]string
}

type HandlerInput struct {
	APIKeyActivationTopic string
}

func NewHandler(log *logger.Logger, sub messaging.Subscriber, in HandlerInput) *Handler {
	topics := make(map[int]string, 1)

	topics[apiKeyActivationIdentifier] = in.APIKeyActivationTopic

	return &Handler{
		log:    log,
		sub:    sub,
		topics: topics,
	}
}

func (h *Handler) Subscribe() error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := h.receiveAPIKeyActivationSubscription(ctx); err != nil {
		return err
	}

	return nil
}
