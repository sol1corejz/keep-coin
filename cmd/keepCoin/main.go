package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/sol1corejz/keep-coin/cmd/config"
	"github.com/sol1corejz/keep-coin/internal/handlers"
	"github.com/sol1corejz/keep-coin/internal/logger"
	"github.com/sol1corejz/keep-coin/internal/storage"
	"go.uber.org/zap"
)

func main() {
	config.ParseFlags()

	if err := logger.Initialize(config.LogLevel); err != nil {
		logger.Log.Fatal("Failed to initialize logger", zap.Error(err))
	}

	if err := storage.Init(); err != nil {
		logger.Log.Fatal("Failed to init storage", zap.Error(err))
		return
	}

	if err := run(); err != nil {
		logger.Log.Fatal("Failed to run server", zap.Error(err))
	}
}

func run() error {
	app := fiber.New()
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowMethods: "GET,POST,OPTIONS",
	}))

	app.Post("/register", handlers.RegisterHandler)
	app.Post("/login", handlers.LoginHandler)

	logger.Log.Info("Running server", zap.String("address", config.RunAddress))
	return app.Listen(config.RunAddress)
}
