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

func (ds *DbService) CreateUser(user *User) (*int64, error) {
	result := ds.db.Create(user)
	if result.Error != nil {
		return nil, result.Error
	}
	return user.Id, nil
}

func (ds *DbService) UpdateUser(user *User) (*int64, error) {
	result := ds.db.Create(user)
	if result.Error != nil {
		return nil, result.Error
	}
	return user.Id, nil
}

func (ds *DbService) DeactivateUser(user *User) (*int64, error) {
	result := ds.db.Create(user)
	if result.Error != nil {
		return nil, result.Error
	}
	return user.Id, nil
}

func (ds *DbService) CreateVideoStory(videoHistory *VideoHistory) (*int64, error) {
	return nil, nil
}
