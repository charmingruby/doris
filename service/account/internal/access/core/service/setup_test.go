package service

import (
	"testing"

	"github.com/charmingruby/doris/lib/instrumentation"
	"github.com/charmingruby/doris/service/account/internal/access/core/model"
	"github.com/charmingruby/doris/service/account/test/memory"
	"github.com/stretchr/testify/suite"
)

type Suite struct {
	suite.Suite

	apiKeyRepo  *memory.APIKeyRepository
	otpRepo     *memory.OTPRepository
	txManager   *memory.TransactionManager
	tokenClient *memory.TokenClient
	evtHandler  *memory.EventHandler
	svc         *Service
}

func (s *Suite) SetupTest() {
	logger := instrumentation.New(instrumentation.LOG_LEVEL_DEBUG)

	s.apiKeyRepo = memory.NewAPIKeyRepository()

	s.otpRepo = memory.NewOTPRepository()

	pub := memory.NewPublisher()

	s.evtHandler = memory.NewEventHandler(*pub)

	s.txManager = memory.NewTransactionManager(s.apiKeyRepo, s.otpRepo)

	s.tokenClient = memory.NewTokenClient()

	s.svc = New(logger, s.apiKeyRepo, s.otpRepo, s.txManager, s.tokenClient, s.evtHandler)
}

func (s *Suite) SetupSubTest() {
	s.apiKeyRepo.Items = []model.APIKey{}
	s.apiKeyRepo.IsHealthy = true

	s.otpRepo.Items = []model.OTP{}
	s.otpRepo.IsHealthy = true

	s.evtHandler.Pub.Messages = []memory.Message{}

	s.tokenClient.Items = make(map[string]memory.TokenPayload)
}

func TestSuite(t *testing.T) {
	suite.Run(t, new(Suite))
}
