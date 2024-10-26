package main

import (
	"auth/config"
	"auth/internal/api/http"
	dbService "auth/internal/db"
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

	var httpConfig config.HttpConfig
	var dbConfig config.DbConfig

	err := httpConfig.MustConfig()
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

	dbService := dbService.NewDbService(db, logger)

	authRepository := http.NewAuthRepository(dbService, logger)
	authRepository.RegisterRouts(app)

	app.Listen(httpConfig.Host + ":" + httpConfig.Port)
}
