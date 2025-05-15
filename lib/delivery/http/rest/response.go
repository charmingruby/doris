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

func NewCreatedResponse(c *gin.Context, msg, id, resource string) {
	resMsg := fmt.Sprintf("%s created successfully", resource)

	if msg != "" {
		resMsg = msg
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": resMsg,
		"data": gin.H{
			"id": id,
		},
	})
}

func NewConflictResponse(c *gin.Context, msg string) {
	c.JSON(http.StatusConflict, gin.H{
		"error": msg,
	})
}

func NewForbiddenResponse(c *gin.Context) {
	c.JSON(http.StatusForbidden, gin.H{
		"error": "insufficient permissions",
	})
}

func NewResourceAlreadyExistsResponse(c *gin.Context, msg string) {
	c.JSON(http.StatusConflict, gin.H{
		"error": msg,
	})
}

func NewTooManyRequestsResponse(c *gin.Context, msg string) {
	c.JSON(http.StatusTooManyRequests, gin.H{
		"error": msg,
	})
}

func NewResourceNotFoundResponse(c *gin.Context, msg string) {
	c.JSON(http.StatusNotFound, gin.H{
		"error": msg,
	})
}

func NewUnprocessableEntityResponse(c *gin.Context, msg string) {
	c.JSON(http.StatusUnprocessableEntity, gin.H{
		"error": msg,
	})
}

func NewOKResponse(c *gin.Context, msg string, data any) {
	res := gin.H{
		"message": msg,
	}

	if data != nil {
		res["data"] = data
	}

	c.JSON(http.StatusOK, res)
}

func NewUnauthorizedResponse(c *gin.Context) {
	c.JSON(http.StatusUnauthorized, gin.H{
		"error": "unauthorized",
	})
}

func NewNoContentResponse(c *gin.Context) {
	c.Status(http.StatusNoContent)
}

func NewBadRequestResponse(c *gin.Context, msg string) {
	c.JSON(http.StatusBadRequest, gin.H{
		"error": msg,
	})
}

type AcceptedResponseIdentifier struct {
	Key   string   `json:"key,omitempty"`
	Value []string `json:"value,omitempty"`
}

func NewAcceptedResponse(c *gin.Context, res []AcceptedResponseIdentifier) {
	fmtRes := gin.H{}

	for _, r := range res {
		if len(r.Value) == 0 {
			continue
		}

		fmtRes[r.Key] = r.Value
	}

	c.JSON(http.StatusAccepted, gin.H{
		"data": fmtRes,
	})
}
