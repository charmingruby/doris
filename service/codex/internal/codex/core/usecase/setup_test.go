package usecase

// import (
// 	"testing"

// 	"github.com/charmingruby/doris/lib/instrumentation"
// 	"github.com/charmingruby/doris/service/codex/internal/codex/core/model"
// 	"github.com/charmingruby/doris/service/codex/test/memory"
// 	"github.com/stretchr/testify/suite"
// )

// type Suite struct {
// 	suite.Suite

// 	codexRepo *memory.CodexRepository
// 	uc        *UseCase
// }

// func (s *Suite) SetupTest() {
// 	logger := instrumentation.New(instrumentation.LOG_LEVEL_DEBUG)
// 	s.codexRepo = memory.NewCodexRepository()

// 	s.uc = New(logger, s.codexRepo)
// }

// func (s *Suite) SetupSubTest() {
// 	s.codexRepo.Items = []model.Codex{}
// 	s.codexRepo.IsHealthy = true
// }

// func TestSuite(t *testing.T) {
// 	suite.Run(t, new(Suite))
// }
