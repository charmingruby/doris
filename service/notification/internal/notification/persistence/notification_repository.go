package persistence

import (
	"context"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/charmingruby/doris/lib/core/pagination"
	"github.com/charmingruby/doris/lib/persistence/dynamo"
	"github.com/charmingruby/doris/service/notification/internal/notification/core/model"
)

var (
	ttl = time.Now().Add(30 * 24 * time.Hour).Unix()
)

type NotificationRepository struct {
	client    *dynamodb.Client
	tableName string
	indexes   tableIndex
}

type tableIndex struct {
	correlationID string
}

type NotificationRepositoryInput struct {
	TableName          string
	CorrelationIDIndex string
}

func NewNotificationRepository(client *dynamodb.Client, in NotificationRepositoryInput) *NotificationRepository {
	return &NotificationRepository{
		client:    client,
		tableName: in.TableName,
		indexes: tableIndex{
			correlationID: in.CorrelationIDIndex,
		},
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

func (r *NotificationRepository) FindManyByCorrelationID(ctx context.Context, correlationID string, page int) ([]model.Notification, error) {
	ctx, cancel := context.WithTimeout(ctx, time.Second*5)
	defer cancel()

	limit := int32(pagination.MAX_ITEMS_PER_PAGE)

	input := &dynamodb.QueryInput{
		TableName:              aws.String(r.tableName),
		IndexName:              aws.String(r.indexes.correlationID),
		KeyConditionExpression: aws.String("correlationId = :correlationId"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":correlationId": &types.AttributeValueMemberS{Value: correlationID},
		},
		ScanIndexForward: aws.Bool(false),
		Limit:            aws.Int32(limit),
	}

	var lastEvaluatedKey map[string]types.AttributeValue
	for i := 1; i < page; i++ {
		if lastEvaluatedKey != nil {
			input.ExclusiveStartKey = lastEvaluatedKey
		}

		result, err := r.client.Query(ctx, input)
		if err != nil {
			return nil, fmt.Errorf("failed to query notifications: %w", err)
		}

		lastEvaluatedKey = result.LastEvaluatedKey

		if lastEvaluatedKey == nil {
			return nil, nil
		}
	}

	if lastEvaluatedKey != nil {
		input.ExclusiveStartKey = lastEvaluatedKey
	}

	result, err := r.client.Query(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("failed to query notifications: %w", err)
	}

	var notifications []model.Notification
	for _, item := range result.Items {
		var notification model.Notification

		if err := attributevalue.UnmarshalMap(item, &notification); err != nil {
			return nil, fmt.Errorf("failed to unmarshal notification: %w", err)
		}

		if pkAttr, ok := item["PK"].(*types.AttributeValueMemberS); ok {
			notificationType := dynamo.ExtractKeyValue(pkAttr.Value)

			convNotificationType, err := model.ParseNotificationType(notificationType)
			if err != nil {
				notification.NotificationType = model.UnknownNotification
			} else {
				notification.NotificationType = convNotificationType
			}
		}
		if skAttr, ok := item["SK"].(*types.AttributeValueMemberS); ok {
			notification.ID = dynamo.ExtractKeyValue(skAttr.Value)
		}

		notifications = append(notifications, notification)
	}

	return notifications, nil
}
