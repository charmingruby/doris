package event

import (
	"context"

	"github.com/charmingruby/doris/lib/proto/gen/notification"
	"github.com/charmingruby/doris/service/notification/internal/notification/core/service"
	"google.golang.org/protobuf/proto"
)

func (h *Handler) receiveOTPNotification(ctx context.Context) error {
	return h.sub.Subscribe(ctx, h.topics[otpNotificationIdentifier], func(message []byte) error {
		var envelope notification.Envelope

		if err := proto.Unmarshal(message, &envelope); err != nil {
			return err
		}

		switch envelope.Type {
		case notification.EnvelopeType_OTP:
			if err := h.svc.NotifyOTP(ctx, service.NotifyOTPInput{
				CorrelationID: envelope.Id,
				To:            envelope.To,
				Content:       envelope.GetOtp().Code,
				RecipientName: envelope.RecipientName,
			}); err != nil {
				h.logger.Error("failed to notify otp", "error", err)
			}
		default:
			h.logger.Error("received unknown notification", "envelope", &envelope)
		}

		return nil
	})
}
