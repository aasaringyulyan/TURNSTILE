package entrypoints

import (
	"github.com/gin-gonic/gin"
	"turnstile/pkg/logging"
)

type response struct {
	Message string `json:"message"`
}

func errorResponse(c *gin.Context, statusCode int, message string, logger logging.Logger) {
	logger.Error(message)
	c.AbortWithStatusJSON(statusCode, response{message})
}
