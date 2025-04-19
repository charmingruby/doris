package config

import (
	"github.com/caarlos0/env/v6"
	"github.com/joho/godotenv"
)

type Base struct {
	LogLevel    string `env:"LOG_LEVEL" envDefault:"info"`
	Environment string `env:"ENVIRONMENT,required" envDefault:"dev"`
}

type Config[T any] struct {
	Base
	Custom T
}

func New[T any]() (Config[T], error) {
	err := godotenv.Load()

	if mustLoadEnvironment() && err != nil {
		return Config[T]{}, err
	}

	cfg := Config[T]{}

	if err := env.Parse(&cfg); err != nil {
		return Config[T]{}, err
	}

	return cfg, nil
}
