package main

import (
	"context"
	data_generator_mock "data-generator-mock"
	"data-generator-mock/models"
	"data-generator-mock/pkg/handler"
	"data-generator-mock/pkg/logging"
	"data-generator-mock/pkg/repository"
	"data-generator-mock/pkg/service"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/spf13/viper"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

//var deltaTime = 30 * time.Second

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
	err = repos.AddEmployee(models.Employee{
		CardNumber: 414141,
		EmployeeID: 2222,
		Rv:         27493955359,
		IsDeleted:  false,
	})
	if err != nil {
		logger.Fatalf("failed to add first Employee: %s", err.Error())
	}

	services := service.NewService(logger, repos)

	//n := 5000
	// Тест 2
	//err = services.DataGenerator.GenSlice(n)
	//if err != nil {
	//	return
	//}

	// Тест 1
	//for i := 0; i < n; i++ {
	//	err = services.DataGenerator.GenNewEmployee()
	//	if err != nil {
	//		logger.Fatalf("failed to generate new Employee: %s", err.Error())
	//	}
	//}

	// Было
	ticker := time.NewTicker(viper.GetDuration("timout.delta_time"))
	go func() {
		for {
			select {
			case <-ticker.C:
				err = services.DataGenerator.GenNewEmployee()
				if err != nil {
					logger.Fatalf("failed to generate new Employee: %s", err.Error())
				}
			}
		}
	}()
	handlers := handler.NewHandler(logger, services)

	srv := new(data_generator_mock.Server)
	go func() {
		if err := srv.Run(viper.GetString("port"), handlers.InitRoutes()); err != nil {
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
