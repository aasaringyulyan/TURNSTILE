package handler

import (
	"data-generator-mock/pkg/logging"
	"data-generator-mock/pkg/service"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	logger   logging.Logger
	services *service.Service
}

func NewHandler(logger logging.Logger, services *service.Service) *Handler {
	return &Handler{
		logger:   logger,
		services: services,
	}
}

func (h *Handler) InitRoutes() *gin.Engine {
	router := gin.New()

	api := router.Group("/data")
	{
		api.GET("/get_by_rv", h.getByRv)
	}

	return router
}
