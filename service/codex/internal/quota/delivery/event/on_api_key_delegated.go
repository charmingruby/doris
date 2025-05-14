package event

import (
	"context"

	"github.com/charmingruby/doris/lib/core/custom_err"
	"github.com/charmingruby/doris/lib/delivery/proto/gen/account"
	"github.com/charmingruby/doris/service/scribe/internal/quota/core/usecase"
	"google.golang.org/protobuf/proto"
)

func (h *Handler) onAPIKeyDelegated(ctx context.Context) error {
	topic := h.topics[apiKeyDelegatedIdentifier]

	return h.sub.Subscribe(ctx, topic, func(message []byte) error {
		var a account.ApiKeyDelegatedEvent

		if err := proto.Unmarshal(message, &a); err != nil {
			h.logger.Error("failed to unmarshal message", "error", err)

			return custom_err.NewErrSerializationFailed(err)
		}

		h.logger.Debug("event received", "topic", topic, "message", &a)

		if err := h.uc.RecalculateQuotaUsageOnTierChange(ctx, usecase.RecalculateQuotaUsageOnTierChangeInput{
			CorrelationID: a.GetId(),
			OldTier:       a.OldTier.String(),
			NewTier:       a.NewTier.String(),
		}); err != nil {
			h.logger.Error("failed to dispatch notification", "error", err)

			return err
		}

		return nil
	})
}
