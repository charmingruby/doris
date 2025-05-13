package email

import (
	"context"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ses"
	"github.com/aws/aws-sdk-go-v2/service/ses/types"
	"github.com/charmingruby/doris/service/notification/internal/notification/core/model"
)

type SES struct {
	client      *ses.Client
	sourceEmail string
}

func NewSES(region, sourceEmail string) (*SES, error) {
	cfg, err := config.LoadDefaultConfig(context.Background())
	if err != nil {
		return nil, err
	}

	cfg.Region = region

	client := ses.NewFromConfig(cfg)

	return &SES{
		client:      client,
		sourceEmail: sourceEmail,
	}, nil
}

func (s *SES) Send(ctx context.Context, notification model.Notification) error {
	ctx, cancel := context.WithTimeout(ctx, time.Second*5)
	defer cancel()

	in := &ses.SendEmailInput{
		Source: aws.String(s.sourceEmail),
		Destination: &types.Destination{
			ToAddresses: []string{notification.To},
		},
		Message: &types.Message{
			Subject: &types.Content{
				Data:    aws.String(getEmailSubject(notification.NotificationType)),
				Charset: aws.String("UTF-8"),
			},
			Body: &types.Body{
				Text: &types.Content{
					Data:    aws.String(notification.Content),
					Charset: aws.String("UTF-8"),
				},
			},
		},
	}

	if _, err := s.client.SendEmail(ctx, in); err != nil {
		return err
	}

	return nil
}

func getEmailSubject(notificationType model.NotificationType) string {
	switch notificationType {
	case model.OTPNotification:
		return "Your Verification Code"
	default:
		return "Notification"
	}
}
