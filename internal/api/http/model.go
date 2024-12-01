package http

import (
	"auth/internal/auth"
	"time"
)

type LogiinUserRequest struct {
	Email    string
	Password string
}

type LoginUserResponse struct {
	Token  *auth.Token      `json:",omitempty"`
	Claims *auth.UserClaims `json:",omitempty"`
	Error  error            `json:",omitempty"`
}

type RegisterUserRequest struct {
	JobRoleId    int
	Address      Address
	Name         string
	SecondName   string
	Surname      string
	Email        string
	Password     string
	Birthday     int64
	BirthdayDate time.Time
}

type Address struct {
	SettlementTypeId int
	Country          string
	Region           string
	District         string
	Settlement       string
	Street           string
	HouseNumber      string
	FlatNumber       string
}
