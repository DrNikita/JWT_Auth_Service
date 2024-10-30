package store

import (
	"log/slog"

	"gorm.io/gorm"
)

type StoreService struct {
	db     *gorm.DB
	logger *slog.Logger
}

func NewDbService(db *gorm.DB, logger *slog.Logger) *StoreService {
	return &StoreService{
		db:     db,
		logger: logger,
	}
}

func (ss *StoreService) CreateUser(user *User) (*int64, error) {
	result := ss.db.Create(user)
	if result.Error != nil {
		return nil, result.Error
	}
	return user.Id, nil
}

func (ss *StoreService) UpdateUser(user *User) (*int64, error) {
	result := ss.db.Create(user)
	if result.Error != nil {
		return nil, result.Error
	}
	return user.Id, nil
}

func (ss *StoreService) DeactivateUser(user *User) (*int64, error) {
	result := ss.db.Create(user)
	if result.Error != nil {
		return nil, result.Error
	}
	return user.Id, nil
}

func (ss *StoreService) CreateVideoStory(videoHistory *VideoHistory) (*int64, error) {
	return nil, nil
}
