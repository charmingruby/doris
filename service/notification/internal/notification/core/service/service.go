package service

import (
	"github.com/charmingruby/doris/lib/instrumentation/logger"
)

type Service struct{}

func New(log *logger.Logger) *Service {
	return &Service{}
}
