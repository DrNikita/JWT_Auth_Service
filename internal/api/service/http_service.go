package service

import (
	"auth/internal/auth"
	"auth/internal/store"
	"log/slog"
)

type HttpService struct {
	storeService store.StoreService
	logger       *slog.Logger
}

func NewHttpService(storeService *store.StoreService, logger *slog.Logger) *HttpService {
	return &HttpService{
		storeService: *storeService,
		logger:       logger,
	}
}

func (hs *HttpService) RegisterUser(user *store.User) (string, error) {
	_, err := hs.storeService.CreateUser(user)
	if err != nil {
		return "", err
	}

	jwt, err := auth.CreateJWT(user)
	if err != nil {
		return "", err
	}

	return jwt, nil
}
