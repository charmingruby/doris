package config

import (
	"github.com/charmingruby/doris/lib/config"
)

type Config config.Config[CustomConfig]

type CustomConfig struct {
	RestServerHost           string `env:"REST_SERVER_HOST" envDefault:"localhost"`
	RestServerPort           string `env:"REST_SERVER_PORT" envDefault:"3001"`
	NatsStream               string `env:"NATS_STREAM,required"`
	SendOTPNotificationTopic string `env:"SEND_OTP_NOTIFICATION_TOPIC,required"`
	SourceEmail              string `env:"SOURCE_EMAIL,required"`
	AWSRegion                string `env:"AWS_REGION,required"`
	NotificatiosnDynamoTable string `env:"NOTIFICATIONS_DYNAMO_TABLE,required"`
	CorrelationIDDynamoIndex string `env:"CORRELATION_ID_DYNAMO_INDEX,required"`
	JWTSecret                string `env:"JWT_SECRET,required"`
	JWTIssuer                string `env:"JWT_ISSUER,required"`
}

func New() (Config, error) {
	cfg, err := config.New[CustomConfig]()
	if err != nil {
		return Config{}, err
	}

	return Config(cfg), nil
}
