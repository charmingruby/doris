package service

import (
	"testing"

	"github.com/charmingruby/doris/lib/instrumentation"
	"github.com/charmingruby/doris/service/hub/internal/identity/core/model"
	"github.com/charmingruby/doris/service/hub/test/memory"
	"github.com/stretchr/testify/suite"
)

type Suite struct {
	suite.Suite

	apiKeyRepo *memory.APIKeyRepository
	otpRepo    *memory.OTPRepository
	evtHandler *memory.EventHandler
	svc        *Service
}

func (s *Suite) SetupTest() {
	logger := instrumentation.New(instrumentation.LOG_LEVEL_DEBUG)

	s.apiKeyRepo = memory.NewAPIKeyRepository()

	s.otpRepo = memory.NewOTPRepository()

	pub := memory.NewPublisher()

	s.evtHandler = memory.NewEventHandler(*pub)

	s.svc = New(logger, s.apiKeyRepo, s.otpRepo, s.evtHandler)
}

func (s *Suite) SetupSubTest() {
	s.apiKeyRepo.Items = []model.APIKey{}
	s.evtHandler.Pub.Messages = []memory.Message{}
	s.apiKeyRepo.IsHealthy = true
}

func TestSuite(t *testing.T) {
	suite.Run(t, new(Suite))
}
