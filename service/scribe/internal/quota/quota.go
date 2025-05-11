package quota

import (
	"github.com/charmingruby/doris/lib/delivery/http/rest"
	"github.com/charmingruby/doris/lib/instrumentation"
	"github.com/charmingruby/doris/lib/validation"
	"github.com/charmingruby/doris/service/scribe/internal/quota/core/repository"
	"github.com/charmingruby/doris/service/scribe/internal/quota/core/usecase"
	"github.com/charmingruby/doris/service/scribe/internal/quota/delivery/http/rest/endpoint"
	"github.com/charmingruby/doris/service/scribe/test/memory"
	"github.com/gin-gonic/gin"
)

type Datasource struct {
	quotaRepo           repository.QuotaRepository
	quotaLimitRepo      repository.QuotaLimitRepository
	quotaLimitUsageRepo repository.QuotaLimitUsageRepository
}

func NewDatasource() (*Datasource, error) {
	quotaRepo := memory.NewQuotaRepository()
	quotaLimitRepo := memory.NewQuotaLimitRepository()
	quotaLimitUsageRepo := memory.NewQuotaLimitUsageRepository()

	return &Datasource{
		quotaRepo:           quotaRepo,
		quotaLimitRepo:      quotaLimitRepo,
		quotaLimitUsageRepo: quotaLimitUsageRepo,
	}, nil
}

func NewUseCase(
	logger *instrumentation.Logger,
	datasource *Datasource,
) *usecase.UseCase {
	return usecase.New(
		logger,
		datasource.quotaRepo,
		datasource.quotaLimitRepo,
		datasource.quotaLimitUsageRepo,
	)
}

func NewHTTPHandler(logger *instrumentation.Logger, r *gin.Engine, mw *rest.Middleware, val *validation.Validator, uc *usecase.UseCase) {
	endpoint.New(logger, r, mw, val, uc).Register()
}
