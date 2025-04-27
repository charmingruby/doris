package event

import (
	"context"

	"github.com/charmingruby/doris/lib/proto/gen/notification"
	"github.com/charmingruby/doris/service/notification/internal/notification/core/service"
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
			if err := h.svc.NotifyApiKeyActivation(ctx, service.NotifyApiKeyActivationInput{
				CorrelationID: envelope.Id,
				To:            envelope.To,
				RecipientName: envelope.RecipientName,
			}); err != nil {
				h.logger.Error("failed to notify api key activation", "error", err)
			}
		default:
			h.logger.Error("received unknown notification", "envelope", &envelope)
		}

		return nil
	})
}
