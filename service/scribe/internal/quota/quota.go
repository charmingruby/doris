package quota

import (
	"github.com/charmingruby/doris/lib/delivery/http/rest"
	"github.com/charmingruby/doris/lib/instrumentation"
	"github.com/charmingruby/doris/lib/validation"
	"github.com/charmingruby/doris/service/scribe/internal/quota/core/repository"
	"github.com/charmingruby/doris/service/scribe/internal/quota/core/usecase"
	"github.com/charmingruby/doris/service/scribe/internal/quota/delivery/http/rest/endpoint"
	"github.com/charmingruby/doris/service/scribe/internal/quota/persistence"
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
)

type Datasource struct {
	quotaRepo repository.QuotaRepository
}

func NewDatasource(db *sqlx.DB) (*Datasource, error) {
	quotaRepo, err := persistence.NewQuotaRepository(db)
	if err != nil {
		return nil, err
	}

	return &Datasource{
		quotaRepo: quotaRepo,
	}, nil
}

func NewUseCase(
	logger *instrumentation.Logger,
	datasource *Datasource,
) *usecase.UseCase {
	return usecase.New(
		logger,
		datasource.quotaRepo,
	)
}

func NewHTTPHandler(logger *instrumentation.Logger, r *gin.Engine, mw *rest.Middleware, val *validation.Validator, uc *usecase.UseCase) {
	endpoint.New(logger, r, mw, val, uc).Register()
}
