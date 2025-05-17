package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/sol1corejz/keep-coin/cmd/config"
	sso "github.com/sol1corejz/keep-coin/internal/clients/sso/grpc"
	"github.com/sol1corejz/keep-coin/internal/handlers"
	"log/slog"
	"os"
	"time"
)

func main() {
	config.ParseFlags()

	log := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))

	ssoCLient, err := sso.New(
		"localhost:44044",
		10*time.Hour,
		log,
		10,
	)
	if err != nil {
		log.Error("Failed to initialize sso grpc client", err.Error())
	}

	if err := run(ssoCLient); err != nil {
		log.Error("Failed to run server", err.Error())
	}
}

func run(ssoClient *sso.Client) error {
	app := fiber.New()

	handler := handlers.New(ssoClient)

	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowMethods: "GET,POST,OPTIONS",
	}))

	app.Post("/register", handler.RegisterGRPC)
	app.Post("/login", handler.LoginGRPC)

	log.Info("Running server on address", config.RunAddress)
	return app.Listen(config.RunAddress)
}
