package event

import (
	"github.com/charmingruby/doris/lib/delivery/messaging"
	"github.com/charmingruby/doris/lib/instrumentation"
)

const (
	codexDocumentUploadedIdentifier = iota
)

type Handler struct {
	logger *instrumentation.Logger
	pub    messaging.Publisher
	topics map[int]string
}

type TopicInput struct {
	CodexDocumentUploaded string
}

func NewHandler(logger *instrumentation.Logger, pub messaging.Publisher, in TopicInput) *Handler {
	topics := make(map[int]string, 1)

	topics[codexDocumentUploadedIdentifier] = in.CodexDocumentUploaded

	return &Handler{
		logger: logger,
		pub:    pub,
		topics: topics,
	}
}
