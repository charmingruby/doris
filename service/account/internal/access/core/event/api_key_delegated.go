package event

import "time"

type APIKeyDelegated struct {
	ID      string
	NewTier string
	OldTier string
	SentAt  time.Time
}
