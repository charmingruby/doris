package event

import (
	"context"

	"github.com/charmingruby/doris/lib/core/custom_err"
	"github.com/charmingruby/doris/lib/core/privilege"
	"github.com/charmingruby/doris/lib/delivery/proto/gen/account"
	"github.com/charmingruby/doris/service/account/internal/access/core/event"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (h *Handler) DispatchAPIKeyActivated(ctx context.Context, message event.APIKeyActivated) error {
	tier, err := privilege.MapTierToProto(message.Tier)
	if err != nil {
		return err
	}

	apiKeyDelegation := account.ApiKeyActivatedEvent{
		Id:     message.ID,
		Tier:   tier,
		SentAt: timestamppb.New(message.SentAt),
	}

	msgBytes, err := proto.Marshal(&apiKeyDelegation)
	if err != nil {
		return custom_err.NewErrSerializationFailed(err)
	}

	topic := h.topics[apiKeyActivatedIdentifier]

	if err := h.pub.Publish(ctx, topic, msgBytes); err != nil {
		return custom_err.NewErrMessagingPublishFailed(topic, msgBytes, err)
	}

	h.logger.Debug("event sent", "topic", topic, "message", message)

	return nil
}
