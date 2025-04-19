package config

import (
	"github.com/charmingruby/doris/lib/config"
)

type CustomConfig struct {
	RestServerHost string `env:"REST_SERVER_HOST" envDefault:"localhost"`
	RestServerPort string `env:"REST_SERVER_PORT" envDefault:"3000"`
}

func New() (config.Config[CustomConfig], error) {
	return config.New[CustomConfig]()
}
