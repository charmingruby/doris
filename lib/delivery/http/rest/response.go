package rest

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func NewPayloadErrResponse(c *gin.Context, reasons []string) {
	c.JSON(http.StatusBadRequest, gin.H{
		"errors": reasons,
	})
}

func NewUncaughtErrResponse(c *gin.Context, err error) {
	c.JSON(http.StatusInternalServerError, gin.H{
		"error": err.Error(),
	})
}

func NewCreatedResponse(c *gin.Context, id, resource string) {
	c.JSON(http.StatusCreated, gin.H{
		"message": fmt.Sprintf("%s created successfully", resource),
		"data": gin.H{
			"id": id,
		},
	})
}

func NewResourceAlreadyExistsResponse(c *gin.Context, resource string) {
	c.JSON(http.StatusConflict, gin.H{
		"error": fmt.Sprintf("%s already exists", resource),
	})
}
