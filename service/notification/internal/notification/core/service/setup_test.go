package service

import (
	"testing"

	"github.com/charmingruby/doris/lib/instrumentation"
	"github.com/charmingruby/doris/service/notification/internal/notification/core/model"
	"github.com/charmingruby/doris/service/notification/test/memory"
	"github.com/stretchr/testify/suite"
)

type Suite struct {
	suite.Suite

	notificationRepo *memory.NotificationRepository
	notifierClient   *memory.Notifier
	svc              *Service
}

func (s *Suite) SetupTest() {
	logger := instrumentation.New(instrumentation.LOG_LEVEL_DEBUG)

	s.notificationRepo = memory.NewNotificationRepository()

	s.notifierClient = memory.NewNotifier()

	s.svc = New(logger, s.notificationRepo, s.notifierClient)
}

func (s *Suite) SetupSubTest() {
	s.notificationRepo.Items = []model.Notification{}
	s.notificationRepo.IsHealthy = true

	s.notifierClient.Items = []model.Notification{}
	s.notifierClient.IsHealthy = true
}

func TestSuite(t *testing.T) {
	suite.Run(t, new(Suite))
}
