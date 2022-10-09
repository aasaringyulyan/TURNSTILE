package internal

import (
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"log-receiver-mock/internal/entrypoints"
	"log-receiver-mock/pkg/httpserver"
	"log-receiver-mock/pkg/logging"
	"os"
	"os/signal"
	"syscall"
)

func Run() {
	logger := logging.GetLogger(viper.GetString("logger.trace_level"))

	// HTTP Server
	engine := gin.New()
	entrypoints.Init(engine, logger)
	httpServer := httpserver.New(engine, viper.GetString("port"))

	// Shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	<-quit
	// Shutdown HTTP Server
	if err := httpServer.Shutdown(); err != nil {
		logger.Errorf("error occured on server shutting down: %s", err.Error())
	}
}
