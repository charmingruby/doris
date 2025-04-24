package event

import "time"

type APIKeyActivation struct {
	ID             string
	To             string
	ActivationCode string
	SentAt         time.Time
}
