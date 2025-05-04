package event

import "time"

type SendNewAPIKeyDelegationMessage struct {
	ID      string
	NewTier string
	OldTier string
	SentAt  time.Time
}
