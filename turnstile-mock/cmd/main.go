package main

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/spf13/viper"
	"log"
	"os"
	"os/signal"
	"syscall"
	turnstile_mock "turnstile-mock"
	"turnstile-mock/pkg/client"
	"turnstile-mock/pkg/logging"
	"turnstile-mock/pkg/repository"
	"turnstile-mock/pkg/service"
)

func main() {
	if err := initConfig(); err != nil {
		log.Fatalf("error initializing configs: %s", err.Error())
	}

	if err := godotenv.Load(); err != nil {
		log.Fatalf("error loading env variables: %s", err.Error())
	}

	logger := logging.GetLogger(viper.GetString("logger.trace_level"))

	db, err := repository.NewPostgresDB(repository.Config{
		Host:     viper.GetString("db.host"),
		Port:     viper.GetString("db.port"),
		Username: viper.GetString("db.username"),
		DBName:   viper.GetString("db.dbname"),
		SSLMode:  viper.GetString("db.sslmode"),
		Password: os.Getenv("DB_PASSWORD"),
	})
	if err != nil {
		logger.Fatalf("failed to initialize db: %s", err.Error())
	}

	repos := repository.NewRepository(logger, db)
	services := service.NewService(logger, repos)
	cli := client.NewClient(logger, services, client.Config{
		MaxConnsPerHost: viper.GetInt("client.max_conns_per_host"),
		Timeout:         viper.GetDuration("client.timeout"),
	})

	if err = cli.PostPassage(client.ConfigTimeout{
		DeltaTime: viper.GetDuration("timout.delta_time"),
	}); err != nil {
		logger.Fatalf("failed to  POST Passage: %s", err.Error())
	}

	srv := new(turnstile_mock.Server)
	go func() {
		if err := srv.Run(viper.GetString("port"), gin.New()); err != nil {
			logger.Fatalf("error occured while running http server: %s", err.Error())
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	<-quit

	if err := srv.Shutdown(context.Background()); err != nil {
		logger.Errorf("error occured on server shutting down: %s", err.Error())
	}

	if err = db.Close(); err != nil {
		logger.Errorf("error occured on db connection close: %s", err.Error())
	}
}

func initConfig() error {
	viper.AddConfigPath("configs")
	viper.SetConfigName("config")
	return viper.ReadInConfig()
}
