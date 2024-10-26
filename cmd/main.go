package main

import (
	"auth-service/configs"
	"auth-service/internal"
	"log"
	"log/slog"
	"os"

	"github.com/gofiber/fiber/v2"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{AddSource: true}))

	var httpConfig configs.HttpConfig
	var dbConfig configs.DbConfig

	err := httpConfig.MustConfig()
	if err != nil {
		log.Fatal(err)
	}
	err = dbConfig.MustConfig()
	if err != nil {
		log.Fatal(err)
	}

	app := fiber.New()

	authRepository := internal.NewAuthRepository(logger)
	authRepository.RegisterRouts(app)

	app.Listen(httpConfig.Host + ":" + httpConfig.Port)
}
