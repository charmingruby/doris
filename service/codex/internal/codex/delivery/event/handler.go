package event

import (
	"context"

	"github.com/charmingruby/doris/lib/delivery/messaging"
	"github.com/charmingruby/doris/lib/instrumentation"
	"github.com/charmingruby/doris/service/codex/internal/codex/core/usecase"
)

const (
	codexDocumentUploadedIdentifier = iota
)

type Handler struct {
	logger *instrumentation.Logger
	pub    messaging.Publisher
	sub    messaging.Subscriber
	topics map[int]string
	uc     *usecase.UseCase
}

type TopicInput struct {
	CodexDocumentUploaded string
}

func NewHandler(logger *instrumentation.Logger, pub messaging.Publisher, sub messaging.Subscriber, in TopicInput) *Handler {
	topics := make(map[int]string, 1)

	topics[codexDocumentUploadedIdentifier] = in.CodexDocumentUploaded

	return &Handler{
		logger: logger,
		pub:    pub,
		sub:    sub,
		topics: topics,
	}
}

func (h *Handler) Subscribe(uc *usecase.UseCase) error {
	ctx := context.Background()

	h.uc = uc

	if err := h.onCodexDocumentUploaded(ctx); err != nil {
		return err
	}

	return nil
}
