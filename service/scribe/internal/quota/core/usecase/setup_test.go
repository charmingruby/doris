package usecase

import (
	"testing"

	"github.com/charmingruby/doris/lib/instrumentation"
	"github.com/charmingruby/doris/service/scribe/internal/quota/core/model"
	"github.com/charmingruby/doris/service/scribe/test/memory"
	"github.com/stretchr/testify/suite"
)

type Suite struct {
	suite.Suite

	quotaRepo           *memory.QuotaRepository
	quotaLimitRepo      *memory.QuotaLimitRepository
	quotaLimitUsageRepo *memory.QuotaLimitUsageRepository
	uc                  *UseCase
}

func (s *Suite) SetupTest() {
	logger := instrumentation.New(instrumentation.LOG_LEVEL_DEBUG)
	s.quotaRepo = memory.NewQuotaRepository()
	s.quotaLimitRepo = memory.NewQuotaLimitRepository()
	s.quotaLimitUsageRepo = memory.NewQuotaLimitUsageRepository()

	s.uc = New(logger, s.quotaRepo, s.quotaLimitRepo, s.quotaLimitUsageRepo)
}

func (s *Suite) SetupSubTest() {
	s.quotaRepo.Items = []model.Quota{}
	s.quotaRepo.IsHealthy = true

	s.quotaLimitRepo.Items = []model.QuotaLimit{}
	s.quotaLimitRepo.IsHealthy = true

	s.quotaLimitUsageRepo.Items = []model.QuotaLimitUsage{}
	s.quotaLimitUsageRepo.IsHealthy = true
}

func TestSuite(t *testing.T) {
	suite.Run(t, new(Suite))
}
