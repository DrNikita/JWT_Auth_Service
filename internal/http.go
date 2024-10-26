package internal

import (
	"log/slog"

	"github.com/gofiber/fiber/v2"
)

type authRepository struct {
	logger *slog.Logger
}

func NewAuthRepository(logger *slog.Logger) *authRepository {
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
