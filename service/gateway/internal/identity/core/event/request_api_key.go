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

type APIKeyRequestEvent struct {
	ID               string
	To               string
	SentAt           time.Time
	VerificationCode string
}

func (h *Handler) PublishAPIKeyRequest(ctx context.Context, event *APIKeyRequestEvent) error {
	envelope := notification.Envelope{
		Id:     event.ID,
		To:     event.To,
		SentAt: timestamppb.New(event.SentAt),
		Type:   notification.EnvelopeType_API_KEY_REQUEST,
		Content: &notification.Envelope_ApiKeyRequest{
			ApiKeyRequest: &notification.APIKeyRequestContent{
				VerificationCode: event.VerificationCode,
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

type APIKeyRequestMessageMapper struct{}

func (m *APIKeyRequestMessageMapper) MapToExternal(in *model.APIKey) *notification.Envelope {
	return &notification.Envelope{
		Id:     in.ID,
		To:     in.Email,
		SentAt: timestamppb.New(time.Now()),
		Type:   notification.EnvelopeType_API_KEY_REQUEST,
		Content: &notification.Envelope_ApiKeyRequest{
			ApiKeyRequest: &notification.APIKeyRequestContent{
				VerificationCode: in.ConfirmationCode,
			},
		},
	}
}

func (m *APIKeyRequestMessageMapper) MapToInternal(in *notification.Envelope) *APIKeyRequestEvent {
	return &APIKeyRequestEvent{
		ID:               in.Id,
		To:               in.To,
		SentAt:           in.SentAt.AsTime(),
		VerificationCode: in.GetApiKeyRequest().VerificationCode,
	}
}

func (m *APIKeyRequestMessageMapper) MapFromBytes(in []byte) (*APIKeyRequestEvent, error) {
	envelope := &notification.Envelope{}

	err := proto.Unmarshal(in, envelope)
	if err != nil {
		return nil, err
	}

	return m.MapToInternal(envelope), nil
}
