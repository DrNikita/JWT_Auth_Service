package http

import (
	"auth/config"
	"auth/internal/auth"
	"context"
	"errors"
	"log/slog"
	"net/http"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/csrf"
	"github.com/gofiber/fiber/v2/middleware/encryptcookie"
	"github.com/gofiber/fiber/v2/middleware/monitor"
	"github.com/gofiber/fiber/v2/middleware/pprof"
	"github.com/halilylm/prometheusfiber"
)

type httpRepository struct {
	httpService *HttpService
	authService *auth.AuthRepository
	authConfig  *config.AuthConfig
	logger      *slog.Logger
	ctx         *context.Context
}

func NewAuthRepository(httpService *HttpService, authService *auth.AuthRepository, authConfig *config.AuthConfig, logger *slog.Logger, ctx *context.Context) *httpRepository {
	return &httpRepository{
		httpService: httpService,
		authService: authService,
		authConfig:  authConfig,
		logger:      logger,
		ctx:         ctx,
	}
}

func (hr *httpRepository) RegisterRouts(app *fiber.App) {
	app.Get("/monitor", monitor.New(monitor.Config{Title: "Auth service monitor"}))

	middleware := prometheusfiber.NewPrometheus(
		prometheusfiber.WithSubSystem("fiber"),
		prometheusfiber.WithMetricPath("/metrics"),
	)
	middleware.Use(app)
	middleware.SetMetricsPath(app)

	app.Use(pprof.New())

	app.Post("/login", hr.login)
	app.Post("/register", hr.registerUser)
	app.Post("/verify-header-token", hr.verifyToken)
	app.Post("/verify-cookie-token", hr.verifyCookieToken)

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
		hr.logger.Error("failed to login", "err", err.Error())
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
	c.JSON(LoginUserResponse{
		Token: token,
	})
	return nil
}

func (hr *httpRepository) registerUser(c *fiber.Ctx) error {
	var user RegisterUserRequest

	err := c.BodyParser(&user)
	if err != nil {
		hr.logger.Error("failed to parse body", "err", err.Error())
		c.Status(http.StatusBadRequest)
		c.JSON(err)
		return err
	}

	_, err = hr.httpService.RegisterUser(user)
	if err != nil {
		hr.logger.Error("failed to register user", "err", err.Error())
		c.Status(http.StatusBadRequest)
		c.JSON(err)
		return err
	}

	c.Status(http.StatusOK)
	c.JSON(user)
	return nil
}

func (hr *httpRepository) verifyCookieToken(c *fiber.Ctx) error {
	token := &auth.Token{
		Access:  c.Cookies("access_token"),
		Refresh: c.Cookies("refresh_token"),
	}

	claims, err := hr.authService.VerifyAccessToken(token.Access)
	if err != nil {
		claims, err := hr.authService.VerifyRefreshToken(token.Refresh)
		if err != nil {
			c.Status(http.StatusBadRequest)
			c.JSON(LoginUserResponse{
				Error: err,
			})
			c.Status(http.StatusUnauthorized)
			c.JSON(LoginUserResponse{
				Error: err,
			})
			return err
		}

		newAccessToken, err := hr.authService.CreateAccessToken(claims)
		if err != nil {
			c.Status(http.StatusBadRequest)
			c.JSON(LoginUserResponse{
				Error: err,
			})
		}

		token.Access = newAccessToken
		c.Status(http.StatusCreated)
		c.JSON(LoginUserResponse{
			Token:  token,
			Claims: claims,
		})
		return nil
	}

	c.Status(http.StatusOK)
	c.JSON(LoginUserResponse{
		Token:  token,
		Claims: claims,
	})
	return nil
}

func (hr *httpRepository) verifyToken(c *fiber.Ctx) error {
	token, err := parseHeaderToken(c.GetReqHeaders())
	if err != nil {
		c.Status(http.StatusBadRequest)
		c.JSON(err)
		return err
	}

	claims, err := hr.authService.VerifyAccessToken(token.Access)
	if err != nil {
		claims, err := hr.authService.VerifyRefreshToken(token.Refresh)
		if err != nil {
			c.Status(http.StatusBadRequest)
			c.JSON(LoginUserResponse{
				Error: err,
			})
			c.Status(http.StatusUnauthorized)
			c.JSON(LoginUserResponse{
				Error: err,
			})
			return err
		}

		newAccessToken, err := hr.authService.CreateAccessToken(claims)
		if err != nil {
			c.Status(http.StatusBadRequest)
			c.JSON(LoginUserResponse{
				Error: err,
			})
		}

		token.Access = newAccessToken
		c.Status(http.StatusCreated)
		c.JSON(LoginUserResponse{
			Token:  token,
			Claims: claims,
		})
		return nil
	}

	c.Status(http.StatusOK)
	c.JSON(LoginUserResponse{
		Token:  token,
		Claims: claims,
	})
	return nil
}

func parseHeaderToken(headers map[string][]string) (*auth.Token, error) {
	token := new(auth.Token)

	accessToken, ok := headers["Access_token"]
	if ok && len(accessToken) > 0 {
		token.Access = accessToken[0]
	} else {
		return nil, errors.New("failed to get access token from headers")
	}
	refreshToken, ok := headers["Refresh_token"]
	if ok && len(refreshToken) > 0 {
		token.Refresh = refreshToken[0]
	} else {
		return nil, errors.New("failed to get refresh token from headers")
	}

	return token, nil
}
