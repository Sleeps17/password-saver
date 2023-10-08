package main

import (
	"log"
	"log/slog"
	"net/http"
	"os"
	"password-saver/internal/config"
	"password-saver/internal/http-handlers/delete"
	"password-saver/internal/http-handlers/get"
	"password-saver/internal/http-handlers/getAll"
	"password-saver/internal/http-handlers/reset"
	"password-saver/internal/http-handlers/save"
	"password-saver/internal/logger"
	"password-saver/internal/router"
	"password-saver/internal/storage"
)

func main() {
	cfg := config.MustLoad()

	logsFile, err := os.OpenFile(cfg.LogsPath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		log.Fatal("Cannot open logs file: ", err, cfg.LogsPath)
	}

	logger := logger.Setup(cfg.Env, logsFile)

	logger.Info("Logger succsesfully created", slog.String("logFile", logsFile.Name()))

	storage, err := storage.New(cfg.StoragePath)
	if err != nil {
		logger.Error("Cannot init storage", slog.String("Error", err.Error()))
		os.Exit(1)
	}

	logger.Info("Storage succsesfully created", slog.String("StorageFile", cfg.StoragePath))

	router := router.Setup()
	router.Post("/saver/save", save.New(logger, storage))
	router.Post("/saver/get", get.New(logger, storage))
	router.Post("/saver/delete", delete.New(logger, storage))
	router.Post("/saver/reset", reset.New(logger, storage))
	router.Post("/saver/getAll", getAll.New(logger, storage))

	srv := &http.Server{
		Addr:         cfg.HTTPServer.Addres,
		Handler:      router,
		ReadTimeout:  cfg.HTTPServer.Timeout,
		WriteTimeout: cfg.HTTPServer.Timeout,
		IdleTimeout:  cfg.HTTPServer.IdleTimeout,
	}

	if err := srv.ListenAndServe(); err != nil {
		logger.Error("Failed to start server", slog.String("Error:", err.Error()))
	}

	logger.Error("Server stopped")
}
