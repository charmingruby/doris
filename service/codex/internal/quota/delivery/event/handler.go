package event

import (
	"context"

	"github.com/charmingruby/doris/lib/delivery/messaging"
	"github.com/charmingruby/doris/lib/instrumentation"
	"github.com/charmingruby/doris/service/codex/internal/quota/core/usecase"
)

const (
	apiKeyDelegatedIdentifier = iota
	apiKeyActivatedIdentifier
)

type Handler struct {
	logger *instrumentation.Logger
	sub    messaging.Subscriber
	topics map[int]string
	uc     *usecase.UseCase
}

type TopicInput struct {
	APIKeyDelegated string
	APIKeyActivated string
}

func NewHandler(logger *instrumentation.Logger, sub messaging.Subscriber, uc *usecase.UseCase, in TopicInput) *Handler {
	topics := make(map[int]string, 2)

	topics[apiKeyDelegatedIdentifier] = in.APIKeyDelegated
	topics[apiKeyActivatedIdentifier] = in.APIKeyActivated

	return &Handler{
		logger: logger,
		sub:    sub,
		topics: topics,
		uc:     uc,
	}
}

func (h *Handler) Subscribe() error {
	ctx := context.Background()

	if err := h.onAPIKeyActivated(ctx); err != nil {
		return err
	}

	if err := h.onAPIKeyDelegated(ctx); err != nil {
		return err
	}

	return nil
}
