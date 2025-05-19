package event

import (
	"context"

	"github.com/charmingruby/doris/lib/core/custom_err"
	"github.com/charmingruby/doris/lib/delivery/proto/gen/codex"
	"github.com/charmingruby/doris/service/codex/internal/codex/core/usecase"
	"google.golang.org/protobuf/proto"
)

func (h *Handler) onCodexDocumentUploaded(ctx context.Context) error {
	topic := h.topics[codexDocumentUploadedIdentifier]

	return h.sub.Subscribe(ctx, topic, func(message []byte) error {
		var c codex.CodexDocumentUploadedEvent

		if err := proto.Unmarshal(message, &c); err != nil {
			h.logger.Error("failed to unmarshal message", "error", err)

			return custom_err.NewErrSerializationFailed(err)
		}

		h.logger.Debug("event received", "topic", topic, "message", &c)

		if err := h.uc.GenerateDocumentEmbedding(ctx, usecase.GenerateDocumentEmbeddingInput{
			DocumentID:    c.GetId(),
			CodexID:       c.GetCodexId(),
			CorrelationID: c.GetCorrelationId(),
			ImageURL:      c.GetImageUrl(),
		}); err != nil {
			h.logger.Error("failed to generate document embedding", "error", err)

			return err
		}

		return nil
	})
}
