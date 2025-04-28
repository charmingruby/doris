package event

import (
	"github.com/charmingruby/doris/lib/delivery/messaging"
)

const (
	otpNotificationIdentifier = iota
)

type Handler struct {
	pub    messaging.Publisher
	topics map[int]string
}

type TopicInput struct {
	OTPNotification string
}

func NewHandler(pub messaging.Publisher, in TopicInput) *Handler {
	topics := make(map[int]string, 1)

	topics[otpNotificationIdentifier] = in.OTPNotification

	return &Handler{
		pub:    pub,
		topics: topics,
	}
}
