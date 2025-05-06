package config

import (
	"github.com/charmingruby/doris/lib/config"
)

type Config config.Config[CustomConfig]

type CustomConfig struct {
	RestServerHost           string `env:"REST_SERVER_HOST" envDefault:"localhost"`
	RestServerPort           string `env:"REST_SERVER_PORT" envDefault:"3000"`
	DatabaseHost             string `env:"DATABASE_HOST,required"`
	DatabasePort             string `env:"DATABASE_PORT,required"`
	DatabaseUser             string `env:"DATABASE_USER,required"`
	DatabasePassword         string `env:"DATABASE_PASSWORD,required"`
	DatabaseName             string `env:"DATABASE_NAME,required"`
	DatabaseSSL              string `env:"DATABASE_SSL,required"`
	JWTSecret                string `env:"JWT_SECRET,required"`
	JWTIssuer                string `env:"JWT_ISSUER,required"`
	NatsStream               string `env:"NATS_STREAM"`
	SendOTPNotificationTopic string `env:"SEND_OTP_NOTIFICATION_TOPIC"`
	APIKeyActivationTopic    string `env:"API_KEY_ACTIVATION_TOPIC"`
	APIKeyDelegationTopic    string `env:"API_KEY_DELEGATION_TOPIC"`
}

func New() (Config, error) {
	cfg, err := config.New[CustomConfig]()
	if err != nil {
		return Config{}, err
	}

	return Config(cfg), nil
}
