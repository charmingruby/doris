package event

import (
	"github.com/charmingruby/doris/lib/delivery/messaging"
	"github.com/charmingruby/doris/lib/instrumentation"
)

const (
	otpNotificationIdentifier = iota
	newAPIKeyDelegationIdentifier
)

type Handler struct {
	logger *instrumentation.Logger
	pub    messaging.Publisher
	topics map[int]string
}

type TopicInput struct {
	OTPNotification     string
	NewAPIKeyDelegation string
}

func NewHandler(logger *instrumentation.Logger, pub messaging.Publisher, in TopicInput) *Handler {
	topics := make(map[int]string, 2)

	topics[otpNotificationIdentifier] = in.OTPNotification
	topics[newAPIKeyDelegationIdentifier] = in.NewAPIKeyDelegation

	return &Handler{
		logger: logger,
		pub:    pub,
		topics: topics,
	}
}
