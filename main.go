package main

import (
	"auth/config"
	"auth/internal/api/http"
	"auth/internal/api/service"
	"auth/internal/auth"
	"auth/internal/store"
	"fmt"
	"log"
	"log/slog"
	"os"

	"github.com/gofiber/fiber/v2"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{AddSource: true}))

	var authConfig config.AuthConfig
	var httpConfig config.HttpConfig
	var dbConfig config.DbConfig

	err := authConfig.MustConfig()
	if err != nil {
		log.Fatal(err)
	}
	err = httpConfig.MustConfig()
	if err != nil {
		log.Fatal(err)
	}
	err = dbConfig.MustConfig()
	if err != nil {
		log.Fatal(err)
	}

	dsn := fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=disable", dbConfig.Host, dbConfig.Port, dbConfig.Username, dbConfig.Name, dbConfig.Password)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}

	app := fiber.New()

	authService := auth.NewAuthService(&authConfig, logger)
	storeService := store.NewDbService(db, logger)
	httpService := service.NewHttpService(authService, storeService, logger)
	authRepository := http.NewAuthRepository(httpService, logger)

	authRepository.RegisterRouts(app)

	app.Listen(httpConfig.Host + ":" + httpConfig.Port)
}
