package event

import "time"

type APIKeyActivation struct {
	ID             string
	To             string
	RecipientName  string
	ActivationCode string
	SentAt         time.Time
}
