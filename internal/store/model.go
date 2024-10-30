package store

import "time"

type User struct {
	Id           *int64   `gorm:"primaryKey;autoIncrement:true"`
	JobRoleId    *int     `gorm:"default:1"`
	JobRole      *JobRole `gorm:"foreignKey:JobRoleId"`
	AddressId    *int64   `gorm:"foreignKey:JobRoleId;default:null"`
	Address      *Address `gorm:"foreignKey:AddressId"`
	Name         *string
	SecondName   *string
	Surname      *string
	Email        *string
	Password     *string
	Birthday     *int64
	BirthdayDate time.Time `gorm:"-"`
	IsActive     *string   `gorm:"default:true"`
}

func (User) TableName() string {
	return "user"
}

type Address struct {
	Id               *int64 `gorm:"primaryKey"`
	SettlementTypeId *int
	SettlementType   *SettlementType `gorm:"foreignKey:SettlementTypeId"`
	Country          *string
	Region           *string
	District         *string
	Settlement       *string
	Street           *string
	HouseNumber      *string
	FlatNumber       *string
}

type VideoHistory struct {
	Id        *int64 `gorm:"primaryKey"`
	UerId     *int64
	User      *User `gorm:"foreignKey:UserId"`
	VideoName *string
	CreatedAt *int64
}

type Role struct {
	Id   *int `gorm:"primaryKey"`
	Name *string
}

type JobRole struct {
	Id      *int `gorm:"primaryKey"`
	Role_id *int
	Name    *string
}

type SettlementType struct {
	Id   *int `gorm:"primaryKey"`
	Name *string
}
