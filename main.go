package main

import (
	"auth/config"
	"auth/internal/api/http"
	"auth/internal/auth"
	"auth/internal/store"
	"context"
	"database/sql"
	"fmt"
	"log"
	"log/slog"
	"os"
	"time"

	_ "github.com/lib/pq"

	"github.com/gofiber/fiber/v2"
)

const contextTimeoutMillis = 5000

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

	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		dbConfig.Host, dbConfig.Port, dbConfig.Username, dbConfig.Password, dbConfig.Name)
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}

	app := fiber.New()

	ctx, done := context.WithTimeout(context.Background(), time.Second*contextTimeoutMillis)
	defer done()

	authService := auth.NewAuthService(&authConfig, logger, &ctx)
	storeService := store.NewDbService(db, logger, &ctx)
	httpService := http.NewHttpService(authService, storeService, logger, &ctx)
	authRepository := http.NewAuthRepository(httpService, logger, &ctx)

	authRepository.RegisterRouts(app)

	app.Listen(httpConfig.Host + ":" + httpConfig.Port)
}
