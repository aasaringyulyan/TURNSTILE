package client

import (
	"errors"
	"fmt"
	"turnstile-mock/pkg/logging"
)

type errorResponse struct {
	Message string `json:"message"`
}

func newErrorResponse(statusCode int, message string, logger logging.Logger) error {
	logger.Error(message)
	return errors.New(fmt.Sprintf("statusCode: %v, errorResponse: %s", statusCode, errorResponse{message}))
}
