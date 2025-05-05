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
		var n notification.Notification

		if err := proto.Unmarshal(message, &n); err != nil {
			h.logger.Error("failed to unmarshal message", "error", err)

			return custom_err.NewErrSerializationFailed(err)
		}

		if n.Type != notification.NotificationType_OTP {
			h.logger.Error("received unknown notification", "message", &n)

			return custom_err.NewErrSerializationFailed(errors.New("unsupported notification type"))
		}

		h.logger.Debug("received otp notification", "message", &n)

		if err := h.svc.DispatchNotification(ctx, service.DispatchNotificationInput{
			CorrelationID:    n.Id,
			To:               n.To,
			Content:          n.GetOtp().Code,
			RecipientName:    n.RecipientName,
			NotificationType: model.OTPNotification,
		}); err != nil {
			h.logger.Error("failed to dispatch notification", "error", err)

			return err
		}

		h.logger.Debug("notification dispatched", "correlation_id", n.Id)

		return nil
	})
}
