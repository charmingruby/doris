package config

import (
	"github.com/charmingruby/doris/lib/config"
)

type Config config.Config[CustomConfig]

type CustomConfig struct {
	RestServerHost           string `env:"REST_SERVER_HOST" envDefault:"localhost"`
	RestServerPort           string `env:"REST_SERVER_PORT" envDefault:"3001"`
	AWSRegion                string `env:"AWS_REGION,required"`
	NotificatiosnDynamoTable string `env:"NOTIFICATIONS_DYNAMO_TABLE,required"`
	NatsStream               string `env:"NATS_STREAM,required"`
	SendOTPNotificationTopic string `env:"SEND_OTP_NOTIFICATION_TOPIC,required"`
	SourceEmail              string `env:"SOURCE_EMAIL,required"`
}

func New() (Config, error) {
	cfg, err := config.New[CustomConfig]()
	if err != nil {
		return Config{}, err
	}

	return Config(cfg), nil
}
