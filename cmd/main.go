package main

import (
	"context"
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
	"github.com/pkg/errors"

	"week4/internal/api"
	"week4/internal/config"
	custumLog "week4/internal/logger"
	"week4/internal/repos"
	"week4/internal/service"
)

func main() {
	if err := godotenv.Load(".env"); err != nil {
		log.Fatal(errors.Wrap(err, "Error loading .env file"))
	}

	// Загружаем конфигурацию из переменных окружения
	var cfg config.AppConfig
	if err := envconfig.Process("", &cfg); err != nil {
		log.Fatal(errors.Wrap(err, "failed to load configuration"))
	}

	// Инициализация логгера
	logger, err := custumLog.NewLogger(cfg.LogLevel)
	if err != nil {
		log.Fatal(errors.Wrap(err, "error initializing logger"))
	}

	storageType := flag.String("storage", "memory", "type of storage to use: 'memory' or 'postgres'")
	flag.Parse()

	ctx := context.Background()

	var serviceInstance service.Service

	switch *storageType {
	case "postgres":
		repository, err := repos.NewPostgres(ctx, cfg.Postgres)
		if err != nil {
			logger.Fatal(errors.Wrap(err, "error initializing postgres"))
		}

		serviceInstance = service.NewService(repository, logger)

		logger.Infof("db - %v", *storageType)
	case "memory":
		repository := repos.NewMemory()

		serviceInstance = service.NewService(repository, logger)

		logger.Infof("db - %v", *storageType)
	default:
		logger.Fatal(errors.Wrap(err, "unknown storage type"))
	}

	// Инициализация API
	app := api.NewRouters(&api.Routers{Service: serviceInstance}, cfg.Rest.Token)

	// Запуск HTTP-сервера в отдельной горутине
	go func() {
		logger.Infof("Starting server on %s", cfg.Rest.ListenAddress)
		if err := app.Listen(":" + cfg.Rest.ListenAddress); err != nil {
			logger.Fatal(errors.Wrap(err, "failed to start server"))
		}
	}()

	// Ожидание системных сигналов для корректного завершения работы
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt, syscall.SIGTERM)
	<-signalChan

	logger.Info("Shutting down gracefully...")
}
