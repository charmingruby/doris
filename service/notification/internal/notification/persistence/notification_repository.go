package persistence

import (
	"context"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/charmingruby/doris/service/notification/internal/notification/core/model"
)

var (
	ttl = time.Now().Add(30 * 24 * time.Hour).Unix()
)

type NotificationRepository struct {
	client    *dynamodb.Client
	tableName string
}

func NewNotificationRepository(client *dynamodb.Client, tableName string) *NotificationRepository {
	return &NotificationRepository{
		client:    client,
		tableName: tableName,
	}
}

func (r *NotificationRepository) Create(ctx context.Context, notification model.Notification) error {
	ctx, cancel := context.WithTimeout(ctx, time.Second*5)
	defer cancel()

	item, err := attributevalue.MarshalMap(map[string]any{
		"PK":            fmt.Sprintf("NOTIFICATION#%s", notification.NotificationType),
		"SK":            fmt.Sprintf("METADATA#%s", notification.ID),
		"correlationId": notification.CorrelationID,
		"to":            notification.To,
		"recipientName": notification.RecipientName,
		"content":       notification.Content,
		"timestamp":     notification.CreatedAt.Unix(),
		"createdAt":     notification.CreatedAt.Format(time.RFC3339),
		"ttl":           ttl,
	})
	if err != nil {
		return fmt.Errorf("failed to marshal notification: %w", err)
	}

	if _, err := r.client.PutItem(ctx, &dynamodb.PutItemInput{
		TableName: aws.String(r.tableName),
		Item:      item,
	}); err != nil {
		return fmt.Errorf("failed to put item in dynamodb: %w", err)
	}

	return nil
}
