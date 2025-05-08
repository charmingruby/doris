package event

import "time"

type APIKeyDelegatedMessage struct {
	ID      string
	NewTier string
	OldTier string
	SentAt  time.Time
}
