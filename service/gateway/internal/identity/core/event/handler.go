package event

import (
	"github.com/charmingruby/doris/lib/delivery/messaging"
)

const (
	apiKeyRequestIdentifier = iota
)

type Handler struct {
	pub messaging.Publisher

	// identifier -> topic
	topics map[int]string
}

type HandlerInput struct {
	APIKeyRequestTopic string
}

func NewHandler(pub messaging.Publisher, in HandlerInput) *Handler {
	topics := make(map[int]string, 1)

	topics[apiKeyRequestIdentifier] = in.APIKeyRequestTopic

	return &Handler{
		pub:    pub,
		topics: topics,
	}
}
