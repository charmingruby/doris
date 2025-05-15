package codex

import (
	"github.com/charmingruby/doris/lib/delivery/http/rest"
	"github.com/charmingruby/doris/lib/instrumentation"
	"github.com/charmingruby/doris/lib/validation"
	"github.com/charmingruby/doris/service/codex/internal/codex/core/repository"
	"github.com/charmingruby/doris/service/codex/internal/codex/core/usecase"
	"github.com/charmingruby/doris/service/codex/internal/codex/delivery/http/rest/endpoint"
	"github.com/charmingruby/doris/service/codex/internal/codex/persistence"
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
)

type Datasource struct {
	codexRepo repository.CodexRepository
}

func NewDatasource(db *sqlx.DB) (*Datasource, error) {
	codexRepo, err := persistence.NewCodexRepository(db)
	if err != nil {
		return nil, err
	}

	return &Datasource{
		codexRepo: codexRepo,
	}, nil
}

func NewUseCase(
	logger *instrumentation.Logger,
	datasource *Datasource,
) *usecase.UseCase {
	return usecase.New(
		logger,
		datasource.codexRepo,
	)
}

func NewHTTPHandler(logger *instrumentation.Logger, r *gin.Engine, mw *rest.Middleware, val *validation.Validator, uc *usecase.UseCase) {
	endpoint.New(logger, r, mw, val, uc).Register()
}
