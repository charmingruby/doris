package event

import "time"

type APIKeyActivated struct {
	ID     string
	Tier   string
	SentAt time.Time
}
