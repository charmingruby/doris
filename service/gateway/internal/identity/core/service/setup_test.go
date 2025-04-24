package service

import (
	"testing"

	"github.com/charmingruby/doris/lib/instrumentation/logger"
	"github.com/charmingruby/doris/service/gateway/internal/identity/core/event"
	"github.com/charmingruby/doris/service/gateway/internal/identity/core/model"
	"github.com/charmingruby/doris/service/gateway/test/memory"
	"github.com/stretchr/testify/suite"
)

type Suite struct {
	suite.Suite

	apiKeyRepo *memory.APIKeyRepository
	pub        *memory.Publisher
	evtHandler *event.Handler
	svc        *Service
}

func (s *Suite) SetupTest() {
	logger := logger.New(logger.LOG_LEVEL_DEBUG)

	s.apiKeyRepo = memory.NewAPIKeyRepository()

	s.pub = memory.NewPublisher()

	s.evtHandler = event.NewHandler(s.pub, event.HandlerInput{
		APIKeyRequestTopic: "notifications.send",
	})

	s.svc = New(logger, s.apiKeyRepo, s.evtHandler)
}

func (s *Suite) SetupSubTest() {
	s.apiKeyRepo.Items = []model.APIKey{}
	s.apiKeyRepo.IsHealthy = true
	s.pub.Messages = []memory.Message{}
}

func TestSuite(t *testing.T) {
	suite.Run(t, new(Suite))
}
