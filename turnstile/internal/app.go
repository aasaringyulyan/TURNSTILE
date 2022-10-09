package internal

import (
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"os"
	"os/signal"
	"syscall"
	"turnstile/internal/entrypoints"
	"turnstile/internal/service"
	"turnstile/internal/service/gripper"
	"turnstile/internal/service/infrastructure"
	"turnstile/internal/service/infrastructure/repo"
	"turnstile/internal/service/scheduler"
	"turnstile/internal/service/sender"
	"turnstile/pkg/client"
	"turnstile/pkg/httpserver"
	"turnstile/pkg/logging"
	"turnstile/pkg/sqlite"
)

func Run() {
	logger := logging.GetLogger(viper.GetString("logger.trace_level"))

	// Init db
	db, err := sqlite.New(sqlite.Config{
		FileName: viper.GetString("db.filename"),
	})
	if err != nil {
		logger.Fatalf("failed to initialize db: %s", err.Error())
	}

	// Repository
	empRepo := repo.NewEmployeeRepo(logger, db, viper.GetString("db.employeeTableName"))
	logRepo := repo.NewLogRepo(logger, db, viper.GetString("db.logTableName"))

	// Infrastructure
	c := client.NewClient(client.Config{
		MaxConnsPerHost: viper.GetInt("client.max_conns_per_host"),
		Timeout:         viper.GetDuration("client.timeout"),
	})
	dataClient := infrastructure.NewDataClient(logger, c, viper.GetString("client.dataClientUrl"))
	logClient := infrastructure.NewLogClient(logger, c, viper.GetString("client.logClientUrl"))

	// Service
	//	Gripper
	dg := gripper.New(dataClient, empRepo, viper.GetString("rvFileName"))
	// 	LogSender
	ls := sender.New(logClient, logRepo)
	// 	Scheduler
	dailyTime := viper.GetDuration("scheduler.dailyTime")
	logTime := viper.GetDuration("scheduler.logTime")
	retryTime := viper.GetDuration("scheduler.retryTime")
	sch := scheduler.New(logger, dg, ls, dailyTime, logTime, retryTime)
	// 	Main service
	tService := service.New(logger, sch, empRepo, logRepo, uint64(viper.GetInt("turnstileId")))

	// First load
	err = tService.StartScheduler()
	if err != nil {
		logger.Fatalf("failed first LoadData: %s", err.Error())
	}

	// HTTP Server
	engine := gin.New()
	entrypoints.Init(engine, logger, tService)
	httpServer := httpserver.New(engine, viper.GetString("port"))

	// Shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	<-quit
	// Shutdown HTTP Server
	if err = httpServer.Shutdown(); err != nil {
		logger.Errorf("error occured on server shutting down: %s", err.Error())
	}
	// Shutdown db
	if err = db.Close(); err != nil {
		logger.Errorf("error occured on db connection close: %s", err.Error())
	}
}
