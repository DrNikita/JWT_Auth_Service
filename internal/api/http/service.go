package http

import (
	"auth/internal/auth"
	"auth/internal/store"
	"context"
	"log/slog"

	"golang.org/x/crypto/bcrypt"
)

type HttpService struct {
	authService  *auth.AuthService
	storeService *store.StoreService
	logger       *slog.Logger
	ctx          *context.Context
}

func NewHttpService(authService *auth.AuthService, storeService *store.StoreService, logger *slog.Logger, ctx *context.Context) *HttpService {
	return &HttpService{
		authService:  authService,
		storeService: storeService,
		logger:       logger,
		ctx:          ctx,
	}
}

func (hs *HttpService) RegisterUser(user RegisterUserRequest) (int64, error) {
	hashedPwd, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return 0, err
	}

	addressID, err := hs.storeService.CreateAddress(&store.Address{
		SettlementTypeId: user.Address.SettlementTypeId,
		Country:          user.Address.Country,
		Region:           user.Address.Region,
		District:         user.Address.District,
		Settlement:       user.Address.Settlement,
		Street:           user.Address.Street,
		HouseNumber:      user.Address.HouseNumber,
		FlatNumber:       user.Address.FlatNumber,
	})
	if err != nil {
		return 0, err
	}

	userID, err := hs.storeService.CreateUser(&store.User{
		JobRoleId:  user.JobRoleId,
		AddressId:  addressID,
		Name:       user.Name,
		SecondName: user.SecondName,
		Surname:    user.Surname,
		Email:      user.Email,
		Password:   string(hashedPwd),
		Birthday:   user.Birthday,
		IsActive:   true,
	})
	if err != nil {
		return 0, err
	}

	return userID, nil
}

func (hs *HttpService) LoginUser(loginData LogiinUserRequest) (*auth.Token, error) {
	user, err := hs.storeService.FindUserByEmail(loginData.Email)
	if err != nil {
		return nil, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(loginData.Password))
	if err != nil {
		return nil, err
	}

	jwt, err := hs.authService.CreateToken(user)
	if err != nil {
		return nil, err
	}

	return jwt, nil
}
