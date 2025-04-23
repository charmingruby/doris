package event

import (
	"github.com/charmingruby/doris/lib/delivery/messaging"
)

const (
	requestAPIKeyIdentifier = iota
)

type Handler struct {
	pub messaging.Publisher

	// identifier -> topic
	topics map[int]string
}

type HandlerInput struct {
	RequestAPIKeyTopic string
}

func NewHandler(pub messaging.Publisher, in HandlerInput) *Handler {
	topics := make(map[int]string, 1)

	topics[requestAPIKeyIdentifier] = in.RequestAPIKeyTopic

	return &Handler{
		pub:    pub,
		topics: topics,
	}
}
