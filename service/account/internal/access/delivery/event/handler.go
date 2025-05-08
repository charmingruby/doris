package event

import (
	"github.com/charmingruby/doris/lib/delivery/messaging"
	"github.com/charmingruby/doris/lib/instrumentation"
)

const (
	sendOTPNotificationIdentifier = iota
	apiKeyDelegatedIdentifier
	apiKeyActivatedIdentifier
)

type Handler struct {
	logger *instrumentation.Logger
	pub    messaging.Publisher
	topics map[int]string
}

type TopicInput struct {
	SendOTPNotification string
	APIKeyDelegated     string
	APIKeyActivated     string
}

func NewHandler(logger *instrumentation.Logger, pub messaging.Publisher, in TopicInput) *Handler {
	topics := make(map[int]string, 3)

	topics[sendOTPNotificationIdentifier] = in.SendOTPNotification
	topics[apiKeyDelegatedIdentifier] = in.APIKeyDelegated
	topics[apiKeyActivatedIdentifier] = in.APIKeyActivated

	return &Handler{
		logger: logger,
		pub:    pub,
		topics: topics,
	}
}
