package event

import (
	"context"

	"github.com/charmingruby/doris/lib/core/custom_err"
	"github.com/charmingruby/doris/lib/delivery/proto/gen/notification"
	"github.com/charmingruby/doris/service/account/internal/access/core/event"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (h *Handler) DispatchSendOTPNotification(ctx context.Context, event event.SendOTPNotification) error {
	notification := notification.NotificationEvent{
		Id:            event.ID,
		To:            event.To,
		RecipientName: event.RecipientName,
		SentAt:        timestamppb.New(event.SentAt),
		Type:          notification.NotificationType_OTP,
		Content: &notification.NotificationEvent_Otp{
			Otp: &notification.OTPContent{
				Code: event.Code,
			},
		},
	}

	msgBytes, err := proto.Marshal(&notification)
	if err != nil {
		return custom_err.NewErrSerializationFailed(err)
	}

	topic := h.topics[sendOTPNotificationIdentifier]

	if err := h.pub.Publish(ctx, topic, msgBytes); err != nil {
		return custom_err.NewErrMessagingPublishFailed(topic, msgBytes, err)
	}

	h.logger.Debug("event sent", "topic", topic)

	return nil
}
