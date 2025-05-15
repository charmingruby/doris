package event

import (
	"context"

	"github.com/charmingruby/doris/lib/core/custom_err"
	"github.com/charmingruby/doris/lib/delivery/proto/gen/account"
	"github.com/charmingruby/doris/service/codex/internal/quota/core/usecase"
	"google.golang.org/protobuf/proto"
)

func (h *Handler) onAPIKeyActivated(ctx context.Context) error {
	topic := h.topics[apiKeyActivatedIdentifier]

	return h.sub.Subscribe(ctx, topic, func(message []byte) error {
		var a account.ApiKeyActivatedEvent

		if err := proto.Unmarshal(message, &a); err != nil {
			h.logger.Error("failed to unmarshal message", "error", err)

			return custom_err.NewErrSerializationFailed(err)
		}

		h.logger.Debug("event received", "topic", topic, "message", &a)

		if err := h.uc.CreateInitialQuotaUsages(ctx, usecase.CreateInitialQuotaUsagesInput{
			Tier:          a.Tier.String(),
			CorrelationID: a.GetId(),
		}); err != nil {
			h.logger.Error("failed to dispatch notification", "error", err)

			return err
		}

		return nil
	})
}
