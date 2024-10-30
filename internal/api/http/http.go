package http

import (
	"log/slog"

	"auth/internal/api/service"
	"auth/internal/store"

	"github.com/gofiber/fiber/v2"
)

type httpRepository struct {
	httpService *service.HttpService
	logger      *slog.Logger
}

func NewAuthRepository(httpService *service.HttpService, logger *slog.Logger) *httpRepository {
	return &httpRepository{
		httpService: httpService,
		logger:      logger,
	}
}

func (hr *httpRepository) RegisterRouts(app *fiber.App) {
	app.Post("/register", hr.registerUser)
}

func (hr *httpRepository) registerUser(c *fiber.Ctx) error {
	var user *store.User

	err := c.BodyParser(&user)
	if err != nil {
		return err
	}

	jwt, err := hr.httpService.RegisterUser(user)
	if err != nil {
		c.SendString(err.Error())
	}

	c.SendString(jwt)

	return nil
}
