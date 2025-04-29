package event

import (
	"context"
	"errors"

	"github.com/charmingruby/doris/lib/core/custom_err"
	"github.com/charmingruby/doris/lib/proto/gen/notification"
	"github.com/charmingruby/doris/service/notification/internal/notification/core/model"
	"github.com/charmingruby/doris/service/notification/internal/notification/core/service"
	"google.golang.org/protobuf/proto"
)

func (h *Handler) receiveOTPNotification(ctx context.Context) error {
	return h.sub.Subscribe(ctx, h.topics[otpNotificationIdentifier], func(message []byte) error {
		var envelope notification.Envelope

		if err := proto.Unmarshal(message, &envelope); err != nil {
			h.logger.Error("failed to unmarshal envelope", "error", err)

			return custom_err.NewErrSerializationFailed(err)
		}

		if envelope.Type != notification.EnvelopeType_OTP {
			h.logger.Error("received unknown notification", "envelope", &envelope)

			return custom_err.NewErrSerializationFailed(errors.New("unsupported envelope type"))
		}

		if err := h.svc.DispatchNotification(ctx, service.DispatchNotificationInput{
			CorrelationID: envelope.Id,
			To:            envelope.To,
			Content:       envelope.GetOtp().Code,
			RecipientName: envelope.RecipientName,
			MessageType:   model.OTPMessageType,
		}); err != nil {
			h.logger.Error("failed to dispatch notification", "error", err)

			return err
		}

		h.logger.Debug("notification dispatched", "correlation_id", envelope.Id)

		return nil
	})
}
