package usecase

import (
	"testing"

	"github.com/charmingruby/doris/lib/instrumentation"
	"github.com/charmingruby/doris/service/codex/internal/quota/core/model"
	"github.com/charmingruby/doris/service/codex/test/memory"
	"github.com/stretchr/testify/suite"
)

type Suite struct {
	suite.Suite

	quotaRepo          *memory.QuotaRepository
	quotaUsageRepo     *memory.QuotaUsageRepository
	transactionManager *memory.TransactionManager
	uc                 *UseCase
}

func (s *Suite) SetupTest() {
	logger := instrumentation.New(instrumentation.LOG_LEVEL_DEBUG)
	s.quotaRepo = memory.NewQuotaRepository()
	s.quotaUsageRepo = memory.NewQuotaUsageRepository()
	s.transactionManager = memory.NewQuotaTransactionManager(s.quotaRepo, s.quotaUsageRepo)

	s.uc = New(logger, s.quotaRepo, s.quotaUsageRepo, s.transactionManager)
}

func (s *Suite) SetupSubTest() {
	s.quotaRepo.Items = []model.Quota{}
	s.quotaRepo.IsHealthy = true

	s.quotaUsageRepo.Items = []model.QuotaUsage{}
	s.quotaUsageRepo.IsHealthy = true
}

func TestSuite(t *testing.T) {
	suite.Run(t, new(Suite))
}
