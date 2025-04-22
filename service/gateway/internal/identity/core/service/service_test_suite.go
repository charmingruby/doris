package service

import (
	"testing"

	"github.com/charmingruby/doris/lib/instrumentation/logger"
	"github.com/charmingruby/doris/service/gateway/internal/identity/core/model"
	"github.com/charmingruby/doris/service/gateway/test/memory"
	"github.com/stretchr/testify/suite"
)

type Suite struct {
	suite.Suite
	apiKeyRepo *memory.APIKeyRepository
	svc        *Service
}

func (s *Suite) SetupTest() {
	logger := logger.New(logger.LOG_LEVEL_DEBUG)

	s.apiKeyRepo = memory.NewAPIKeyRepository()
	s.svc = New(logger, s.apiKeyRepo)
}

func (s *Suite) SetupSubTest() {
	s.apiKeyRepo.Items = []model.APIKey{}
	s.apiKeyRepo.IsHealthy = true
}

func TestSuite(t *testing.T) {
	suite.Run(t, new(Suite))
}
