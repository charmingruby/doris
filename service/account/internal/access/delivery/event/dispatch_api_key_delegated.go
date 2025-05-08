package event

import (
	"context"

	"github.com/charmingruby/doris/lib/core/custom_err"
	"github.com/charmingruby/doris/lib/delivery/proto/gen/account"
	"github.com/charmingruby/doris/service/account/internal/access/core/event"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (h *Handler) DispatchAPIKeyDelegated(ctx context.Context, event event.APIKeyDelegatedMessage) error {
	apiKeyDelegation := account.ApiKeyDelegated{
		Id:      event.ID,
		NewTier: event.NewTier,
		OldTier: event.OldTier,
		SentAt:  timestamppb.New(event.SentAt),
	}

	msgBytes, err := proto.Marshal(&apiKeyDelegation)
	if err != nil {
		return custom_err.NewErrSerializationFailed(err)
	}

	topic := h.topics[apiKeyDelegatedIdentifier]

	if err := h.pub.Publish(ctx, topic, msgBytes); err != nil {
		return custom_err.NewErrMessagingPublishFailed(topic, msgBytes, err)
	}

	h.logger.Debug("event sent", "topic", topic)

	return nil
}
