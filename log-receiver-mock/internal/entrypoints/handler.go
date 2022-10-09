package entrypoints

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"log-receiver-mock/internal/models"
	"log-receiver-mock/pkg/logging"
	"net/http"
)

type Handler struct {
	logger logging.Logger
}

func Init(engine *gin.Engine, logger logging.Logger) *Handler {
	h := &Handler{
		logger: logger,
	}

	h.InitRoutes(engine)
	return h
}

func (h *Handler) InitRoutes(engine *gin.Engine) {
	engine.Use(gin.Logger())
	engine.Use(gin.Recovery())

	engine.POST("/logs", h.receiveLogs)
}

func (h *Handler) receiveLogs(ctx *gin.Context) {
	var logs models.PassageLogsForApi

	body, err := ioutil.ReadAll(ctx.Request.Body)
	if err != nil {
		errorResponse(ctx, http.StatusInternalServerError, err.Error(), h.logger)
		return
	}

	err = json.Unmarshal(body, &logs)
	if err != nil {
		errorResponse(ctx, http.StatusInternalServerError, err.Error(), h.logger)
		return
	}

	fmt.Printf("Receive %d logs:\n", len(logs.Logs))
	for i, v := range logs.Logs {
		fmt.Printf("%d) cardID = %d\n", i, v.CardID)
	}

	ctx.JSON(http.StatusOK, nil)
}
