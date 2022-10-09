package handler

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func (h *Handler) getByRv(ctx *gin.Context) {
	logger := h.logger.Logger

	rv, err := getRv(ctx)
	if err != nil {
		newErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}
	logger.Info("Got rv")

	employees, err := h.services.DataGenerator.GetByRv(rv)
	if err != nil {
		newErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"data": employees,
	})
}
