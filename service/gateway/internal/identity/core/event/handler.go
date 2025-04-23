package event

import (
	"github.com/charmingruby/doris/lib/delivery/messaging"
)

type Handler struct {
	pub messaging.Publisher
}

func NewHandler(pub messaging.Publisher) *Handler {
	return &Handler{pub: pub}
}
