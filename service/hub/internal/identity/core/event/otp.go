package event

import "time"

type OTP struct {
	ID            string
	To            string
	RecipientName string
	Code          string
	SentAt        time.Time
}
