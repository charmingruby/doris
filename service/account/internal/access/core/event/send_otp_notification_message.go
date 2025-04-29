package event

import "time"

type SendOTPNotificationMessage struct {
	ID            string
	To            string
	RecipientName string
	Code          string
	SentAt        time.Time
}
