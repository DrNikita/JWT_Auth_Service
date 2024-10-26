package http

import (
	"log/slog"

	"auth/internal/db"

	"github.com/gofiber/fiber/v2"
)

type authRepository struct {
	dbService *db.DbService
	logger    *slog.Logger
}

func NewAuthRepository(dbService *db.DbService, logger *slog.Logger) *authRepository {
	return &authRepository{
		logger: logger,
	}
}

func (ar *authRepository) RegisterRouts(app *fiber.App) {
	app.Get("/", handler)
}

func handler(c *fiber.Ctx) error {
	return nil
}
