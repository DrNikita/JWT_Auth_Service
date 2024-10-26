package db

import (
	"log/slog"

	"gorm.io/gorm"
)

type DbService struct {
	db     *gorm.DB
	logger *slog.Logger
}

func NewDbService(db *gorm.DB, logger *slog.Logger) *DbService {
	return &DbService{
		db:     db,
		logger: logger,
	}
}

func (ds *DbService) CreateUser(user *User) (int64, error) {
	return 0, nil
}

func (ds *DbService) UpdateUser(user *User) (int64, error) {
	return 0, nil
}

func (ds *DbService) CreateVideo(videoHistory *VideoHistory) (int64, error) {
	return 0, nil
}
