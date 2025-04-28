package event

import (
	"github.com/charmingruby/doris/lib/delivery/messaging"
)

const (
	otpIdentifier = iota
)

type Handler struct {
	pub messaging.Publisher

	// identifier -> topic
	topics map[int]string
}

type TopicInput struct {
	OTPTopic string
}

func NewHandler(pub messaging.Publisher, in TopicInput) *Handler {
	topics := make(map[int]string, 1)

	topics[otpIdentifier] = in.OTPTopic

	return &Handler{
		pub:    pub,
		topics: topics,
	}
}
