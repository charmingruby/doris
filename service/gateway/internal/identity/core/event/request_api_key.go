package event

import (
	"context"
	"time"

	"github.com/charmingruby/doris/lib/core/custom_err"
	"github.com/charmingruby/doris/lib/proto/gen/notification"
	"github.com/charmingruby/doris/service/gateway/internal/identity/core/model"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (h *Handler) PublishRequestAPIKey(ctx context.Context, ak model.APIKey) error {
	msg := notification.Envelope{
		Id:     ak.ID,
		To:     ak.Email,
		SentAt: timestamppb.New(time.Now()),
		Type:   notification.EnvelopeType_API_KEY_REQUEST,
		Content: &notification.Envelope_ApiKeyRequest{
			ApiKeyRequest: &notification.APIKeyRequestContent{
				VerificationCode: ak.ConfirmationCode,
			},
		},
	}

	msgBytes, err := proto.Marshal(&msg)
	if err != nil {
		return custom_err.NewErrSerializationFailed(err)
	}

	topic := h.topics[requestAPIKeyIdentifier]

	if err := h.pub.Publish(ctx, topic, msgBytes); err != nil {
		return custom_err.NewErrMessagingPublishFailed(topic, msgBytes, err)
	}

	return nil
}
