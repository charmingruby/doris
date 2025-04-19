package config

import (
	"os"
)

const (
	ENVIRONMENT_DEVELOPMENT string = "dev"
	ENVIRONMENT_PRODUCTION  string = "prod"
)

func mustLoadEnvironment() bool {
	environment, ok := os.LookupEnv("ENVIRONMENT")

	if !ok {
		environment = ENVIRONMENT_DEVELOPMENT
	}

	switch environment {
	case ENVIRONMENT_DEVELOPMENT:
		return true
	case ENVIRONMENT_PRODUCTION:
		return false
	default:
		return true
	}
}
