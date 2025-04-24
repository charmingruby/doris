package event

import "time"

type APIKeyRequest struct {
	ID               string
	To               string
	ConfirmationCode string
	SentAt           time.Time
}
