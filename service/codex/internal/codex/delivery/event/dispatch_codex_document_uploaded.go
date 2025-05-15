package event

import (
	"context"

	"github.com/charmingruby/doris/lib/core/custom_err"
	"github.com/charmingruby/doris/lib/delivery/proto/gen/codex"
	"github.com/charmingruby/doris/service/codex/internal/codex/core/event"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (h *Handler) DispatchCodexDocumentUploaded(ctx context.Context, message event.CodexDocumentUploaded) error {
	codexDocumentUploaded := codex.CodexDocumentUploadedEvent{
		Id:            message.ID,
		CodexId:       message.CodexID,
		CorrelationId: message.CorrelationID,
		ImageUrl:      message.ImageURL,
		SentAt:        timestamppb.New(message.SentAt),
	}

	msgBytes, err := proto.Marshal(&codexDocumentUploaded)
	if err != nil {
		return custom_err.NewErrSerializationFailed(err)
	}

	topic := h.topics[codexDocumentUploadedIdentifier]

	if err := h.pub.Publish(ctx, topic, msgBytes); err != nil {
		return custom_err.NewErrMessagingPublishFailed(topic, msgBytes, err)
	}

	h.logger.Debug("event sent", "topic", topic, "message", message)

	return nil
}
