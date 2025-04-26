package event

import (
	"github.com/charmingruby/doris/lib/delivery/messaging"
)

const (
	apiKeyActivationIdentifier = iota
)

type Handler struct {
	pub messaging.Publisher

	// identifier -> topic
	topics map[int]string
}

type TopicInput struct {
	APIKeyActivationTopic string
}

func NewHandler(pub messaging.Publisher, in TopicInput) *Handler {
	topics := make(map[int]string, 1)

	topics[apiKeyActivationIdentifier] = in.APIKeyActivationTopic

	return &Handler{
		pub:    pub,
		topics: topics,
	}
}
