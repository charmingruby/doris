package config

import (
	"github.com/charmingruby/doris/lib/config"
)

type Config config.Config[CustomConfig]

type CustomConfig struct {
	RestServerHost         string `env:"REST_SERVER_HOST" envDefault:"localhost"`
	RestServerPort         string `env:"REST_SERVER_PORT" envDefault:"3000"`
	NatsStream             string `env:"NATS_STREAM"`
	NotificationsSendTopic string `env:"NOTIFICATIONS_SEND_TOPIC"`
	MongoURL               string `env:"MONGO_URL,required"`
	MongoDatabase          string `env:"MONGO_DB,required"`
}

func New() (Config, error) {
	cfg, err := config.New[CustomConfig]()
	if err != nil {
		return Config{}, err
	}

	return Config(cfg), nil
}
