package event

import (
	"context"

	"github.com/charmingruby/doris/lib/core/custom_err"
	"github.com/charmingruby/doris/lib/proto/gen/notification"
	"github.com/charmingruby/doris/service/gateway/internal/identity/core/event"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (h *Handler) PublishAPIKeyRequest(ctx context.Context, event *event.APIKeyRequest) error {
	envelope := notification.Envelope{
		Id:     event.ID,
		To:     event.To,
		SentAt: timestamppb.New(event.SentAt),
		Type:   notification.EnvelopeType_API_KEY_REQUEST,
		Content: &notification.Envelope_ApiKeyRequest{
			ApiKeyRequest: &notification.APIKeyRequestContent{
				VerificationCode: event.ConfirmationCode,
			},
		},
	}

	msgBytes, err := proto.Marshal(&envelope)
	if err != nil {
		return custom_err.NewErrSerializationFailed(err)
	}

	topic := h.topics[apiKeyRequestIdentifier]

	if err := h.pub.Publish(ctx, topic, msgBytes); err != nil {
		return custom_err.NewErrMessagingPublishFailed(topic, msgBytes, err)
	}

	return nil
}
