package http

import "time"

type LogiinUserRequest struct {
	Email    string
	Password string
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
	SettlementTypeId int `json:"settlement_type_id"`
	Country          string
	Region           string
	District         string
	Settlement       string
	Street           string
	HouseNumber      string
	FlatNumber       string
}
