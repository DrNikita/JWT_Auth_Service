package http

import (
	"log/slog"
	"strconv"

	"auth/internal/db"

	"github.com/gofiber/fiber/v2"
)

type authRepository struct {
	dbService *db.DbService
	logger    *slog.Logger
}

func NewAuthRepository(dbService *db.DbService, logger *slog.Logger) *authRepository {
	return &authRepository{
		dbService: dbService,
		logger:    logger,
	}
}

func (ar *authRepository) RegisterRouts(app *fiber.App) {
	app.Post("/register", ar.registerUser)
}

func (ar *authRepository) registerUser(c *fiber.Ctx) error {
	var user db.User

	err := c.BodyParser(&user)
	if err != nil {
		return err
	}

	id, err := ar.dbService.CreateUser(&user)
	if err != nil {
		return err
	}

	err = c.SendString(strconv.Itoa(int(*id)))
	if err != nil {
		return err
	}

	return nil
}
