package service

import (
	"github.com/charmingruby/doris/lib/instrumentation"
)

type Service struct {
	logger *instrumentation.Logger
}

func New(logger *instrumentation.Logger) *Service {
	return &Service{
		logger: logger,
	}
}
