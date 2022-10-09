package entrypoints

import (
	"github.com/gin-gonic/gin"
	ginprometheus "github.com/zsais/go-gin-prometheus"
	"net/http"
	"turnstile/internal/models"
	"turnstile/internal/service"
	"turnstile/pkg/logging"
)

type Handler struct {
	tService service.Turnstile
	logger   logging.Logger
}

func Init(engine *gin.Engine, logger logging.Logger, services service.Turnstile) *Handler {
	h := &Handler{
		tService: services,
		logger:   logger,
	}

	h.InitRoutes(engine)
	return h
}

func (h *Handler) InitRoutes(engine *gin.Engine) {
	engine.Use(gin.Logger())
	engine.Use(gin.Recovery())

	//Gin Web Framework Prometheus metrics exporter
	p := ginprometheus.NewPrometheus("gin")
	p.Use(engine)

	engine.POST("/check", h.getPassage)
	engine.POST("/logs", h.getLogs)
}

func (h *Handler) getPassage(ctx *gin.Context) {
	var passage models.PassageCheck

	if err := ctx.ShouldBindJSON(&passage); err != nil {
		errorResponse(ctx, http.StatusBadRequest, "invalid input body", h.logger)
		return
	}

	isAccepted, err := h.tService.CheckHandler(passage)
	if err != nil {
		errorResponse(ctx, http.StatusInternalServerError, err.Error(), h.logger)
		return
	}

	var status int

	switch isAccepted {
	case true:
		status = http.StatusOK
	case false:
		status = http.StatusUnauthorized
	}

	ctx.JSON(status, gin.H{
		"accepted": isAccepted,
	})
}

func (h *Handler) getLogs(ctx *gin.Context) {
	var logs models.PassageLogsLinux

	if err := ctx.ShouldBindJSON(&logs); err != nil {
		errorResponse(ctx, http.StatusBadRequest, "invalid input body", h.logger)
		return
	}

	if len(logs.Logs) < 1 {
		errorResponse(ctx, http.StatusBadRequest, "invalid input body", h.logger)
		return
	}

	// Request to service
	err := h.tService.LogHandler(logs)
	if err != nil {
		errorResponse(ctx, http.StatusInternalServerError, err.Error(), h.logger)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code": http.StatusOK,
	})
}
