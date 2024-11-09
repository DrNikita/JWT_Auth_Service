package http

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/gofiber/fiber/v2"
)

type httpRepository struct {
	httpService *HttpService
	logger      *slog.Logger
	ctx         *context.Context
}

func NewAuthRepository(httpService *HttpService, logger *slog.Logger, ctx *context.Context) *httpRepository {
	return &httpRepository{
		httpService: httpService,
		logger:      logger,
		ctx:         ctx,
	}
}

func (hr *httpRepository) RegisterRouts(app *fiber.App) {
	app.Post("/login", hr.login)
	app.Post("/register", hr.registration)
}

func (hr *httpRepository) login(c *fiber.Ctx) error {
	var loginUser LogiinUserRequest

	err := c.BodyParser(&loginUser)
	if err != nil {
		c.Status(http.StatusBadRequest)
		c.JSON(err)
		return err
	}

	token, err := hr.httpService.LoginUser(loginUser)
	if err != nil {
		c.Status(http.StatusBadRequest)
		c.JSON(err)
		return err
	}

	c.Status(http.StatusOK)
	c.JSON(token)
	return nil
}

func (hr *httpRepository) registration(c *fiber.Ctx) error {
	var user RegisterUserRequest

	err := c.BodyParser(&user)
	if err != nil {
		c.Status(http.StatusBadRequest)
		c.JSON(err)
		return err
	}

	_, err = hr.httpService.RegisterUser(user)
	if err != nil {
		c.Status(http.StatusBadRequest)
		c.JSON(err)
		return err
	}

	c.Status(http.StatusOK)
	c.JSON(user)
	return nil
}
