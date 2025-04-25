package event

import (
	"context"

	"github.com/charmingruby/doris/lib/core/custom_err"
	"github.com/charmingruby/doris/lib/proto/gen/notification"
	"github.com/charmingruby/doris/service/hub/internal/identity/core/event"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (h *Handler) SendAPIKeyActivation(ctx context.Context, event *event.APIKeyActivation) error {
	envelope := notification.Envelope{
		Id:            event.ID,
		To:            event.To,
		RecipientName: event.RecipientName,
		SentAt:        timestamppb.New(event.SentAt),
		Type:          notification.EnvelopeType_API_KEY_ACTIVATION,
		Content: &notification.Envelope_ApiKeyActivation{
			ApiKeyActivation: &notification.APIKeyActivationContent{
				ActivationCode: event.ActivationCode,
			},
		},
	}

	msgBytes, err := proto.Marshal(&envelope)
	if err != nil {
		return custom_err.NewErrSerializationFailed(err)
	}

	topic := h.topics[apiKeyActivationIdentifier]

	if err := h.pub.Publish(ctx, topic, msgBytes); err != nil {
		return custom_err.NewErrMessagingPublishFailed(topic, msgBytes, err)
	}

	return nil
}
