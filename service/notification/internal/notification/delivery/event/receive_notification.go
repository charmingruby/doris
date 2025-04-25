package event

import (
	"context"

	"github.com/charmingruby/doris/lib/proto/gen/notification"
	"google.golang.org/protobuf/proto"
)

func (h *Handler) receiveNotification(ctx context.Context) error {
	return h.sub.Subscribe(ctx, h.topics[receiveNotificationIdentifier], func(message []byte) error {
		var envelope notification.Envelope

		if err := proto.Unmarshal(message, &envelope); err != nil {
			return err
		}

		switch envelope.Type {
		case notification.EnvelopeType_API_KEY_ACTIVATION:
			h.log.Debug("received api key activation", "envelope", &envelope)
		default:
			h.log.Debug("received unknown notification", "envelope", &envelope)
		}

		return nil
	})
}
