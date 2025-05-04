package event

import "time"

type APIKeyDelegationMessage struct {
	ID      string
	NewTier string
	OldTier string
	SentAt  time.Time
}
