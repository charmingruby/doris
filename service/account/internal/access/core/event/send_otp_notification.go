package event

import "time"

type SendOTPNotification struct {
	ID            string
	To            string
	RecipientName string
	Code          string
	SentAt        time.Time
}
