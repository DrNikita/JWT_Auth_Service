package http

import (
	"auth/config"
	"auth/internal/auth"
	"context"
	"log/slog"
	"net/http"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/csrf"
	"github.com/gofiber/fiber/v2/middleware/encryptcookie"
)

// var cookieSecret string

// func Init() {
// 	cookieSecret = encryptcookie.GenerateKey()
// }

type httpRepository struct {
	httpService *HttpService
	authConfig  *config.AuthConfig
	logger      *slog.Logger
	ctx         *context.Context
}

func NewAuthRepository(httpService *HttpService, authConfig *config.AuthConfig, logger *slog.Logger, ctx *context.Context) *httpRepository {
	return &httpRepository{
		httpService: httpService,
		authConfig:  authConfig,
		logger:      logger,
		ctx:         ctx,
	}
}

func (hr *httpRepository) RegisterRouts(app *fiber.App) {
	app.Post("/login", hr.login)
	app.Post("/register", hr.registration)
	app.Post("/verify", hr.verifyToken)

	app.Use(encryptcookie.New(encryptcookie.Config{
		Key:    hr.authConfig.CookieSecret,
		Except: []string{csrf.ConfigDefault.CookieName}, // exclude CSRF cookie
	}))
	app.Use(csrf.New(csrf.Config{
		KeyLookup:      "header:" + csrf.HeaderName,
		CookieSameSite: "Lax",
		CookieSecure:   true,
		CookieHTTPOnly: false,
	}))
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

	c.Cookie(&fiber.Cookie{
		Name:    "access_token",
		Value:   token.Access,
		Expires: time.Now().Add(15 * time.Minute),
	})

	c.Cookie(&fiber.Cookie{
		Name:    "refresh_token",
		Value:   token.Refresh,
		Expires: time.Now().Add(24 * time.Hour),
	})

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

func (hr *httpRepository) verifyToken(c *fiber.Ctx) error {
	c.Status(http.StatusOK)
	token := auth.Token{
		Access:  c.Cookies("access_token"),
		Refresh: c.Cookies("refresh_token"),
	}
	c.JSON(token)
	return nil
}
