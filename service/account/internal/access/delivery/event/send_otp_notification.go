package event

import (
	"context"

	"github.com/charmingruby/doris/lib/core/custom_err"
	"github.com/charmingruby/doris/lib/proto/gen/notification"
	"github.com/charmingruby/doris/service/account/internal/access/core/event"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (h *Handler) SendOTPNotification(ctx context.Context, event *event.SendOTPNotificationMessage) error {
	envelope := notification.Envelope{
		Id:            event.ID,
		To:            event.To,
		RecipientName: event.RecipientName,
		SentAt:        timestamppb.New(event.SentAt),
		Type:          notification.EnvelopeType_OTP,
		Content: &notification.Envelope_Otp{
			Otp: &notification.OTPContent{
				Code: event.Code,
			},
		},
	}

	msgBytes, err := proto.Marshal(&envelope)
	if err != nil {
		return custom_err.NewErrSerializationFailed(err)
	}

	topic := h.topics[otpNotificationIdentifier]

	if err := h.pub.Publish(ctx, topic, msgBytes); err != nil {
		return custom_err.NewErrMessagingPublishFailed(topic, msgBytes, err)
	}

	return nil
}
