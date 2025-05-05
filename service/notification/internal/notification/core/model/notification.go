package model

import (
	"time"

	"github.com/charmingruby/doris/lib/core/id"
)

type NotificationType string

const (
	OTPNotification NotificationType = "otp"
)

type NotificationInput struct {
	CorrelationID    string           `json:"correlation_id"`
	To               string           `json:"to"`
	RecipientName    string           `json:"recipient_name"`
	Content          string           `json:"content"`
	NotificationType NotificationType `json:"notification_type"`
}

func NewNotification(in NotificationInput) *Notification {
	return &Notification{
		ID:               id.New(),
		CorrelationID:    in.CorrelationID,
		To:               in.To,
		RecipientName:    in.RecipientName,
		Content:          in.Content,
		NotificationType: in.NotificationType,
		CreatedAt:        time.Now(),
	}
}

type Notification struct {
	ID               string           `json:"id"`
	CorrelationID    string           `json:"correlation_id"`
	To               string           `json:"to"`
	RecipientName    string           `json:"recipient_name"`
	Content          string           `json:"content"`
	NotificationType NotificationType `json:"notification_type"`
	CreatedAt        time.Time        `json:"created_at"`
}
