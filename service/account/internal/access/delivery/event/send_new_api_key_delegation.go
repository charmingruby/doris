package event

import (
	"context"

	"github.com/charmingruby/doris/service/account/internal/access/core/event"
)

func (h *Handler) SendNewAPIKeyDelegation(ctx context.Context, event event.SendNewAPIKeyDelegationMessage) error {
	h.logger.Debug("sent new api key delegation event")

	return nil
}
