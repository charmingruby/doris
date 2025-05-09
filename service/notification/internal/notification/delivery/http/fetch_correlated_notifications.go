package endpoint

import (
	"context"
	"strconv"

	"github.com/charmingruby/doris/lib/delivery/http/rest"
	"github.com/charmingruby/doris/service/notification/internal/notification/core/usecase"
	"github.com/gin-gonic/gin"
)

func (e *Endpoint) makeFetchCorrelatedNotifications(c *gin.Context) {
	apiKeyID := c.GetString("api-key-id")
	if apiKeyID == "" {
		rest.NewUnauthorizedResponse(c)
		return
	}

	pageQuery := c.Query("page")
	page := 0

	if pageQuery != "" {
		convPage, err := strconv.Atoi(pageQuery)

		if err != nil || convPage < 0 {
			rest.NewBadRequestResponse(c, "invalid page value")
			return
		}

		page = convPage
	}

	notifications, err := e.uc.FetchCorrelatedNotifications(context.Background(), usecase.FetchCorrelatedNotificationsInput{
		CorrelationID: apiKeyID,
		Page:          page,
	})
	if err != nil {
		rest.HandleHTTPError(c, e.logger, err)
		return
	}

	rest.NewOKResponse(c, "retrieved correlated notifications", notifications)
}
