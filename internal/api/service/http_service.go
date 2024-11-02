package service

import (
	"auth/internal/auth"
	"auth/internal/store"
	"log/slog"
)

type HttpService struct {
	authService  *auth.AuthService
	storeService *store.StoreService
	logger       *slog.Logger
}

func NewHttpService(authService *auth.AuthService, storeService *store.StoreService, logger *slog.Logger) *HttpService {
	return &HttpService{
		authService:  authService,
		storeService: storeService,
		logger:       logger,
	}
}

func (hs *HttpService) RegisterUser(user *store.User) (*auth.Token, error) {
	_, err := hs.storeService.CreateUser(user)
	if err != nil {
		return nil, err
	}

	jwt, err := hs.authService.CreateToken(user)
	if err != nil {
		return nil, err
	}

	return jwt, nil
}
