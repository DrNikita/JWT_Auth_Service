package store

type User struct {
	Id         int64
	JobRoleId  int
	AddressId  int64
	Name       string
	SecondName string
	Surname    string
	Email      string
	Password   string
	Birthday   int64
	IsActive   bool
}

type Address struct {
	Id               int64
	SettlementTypeId int
	Country          string
	Region           string
	District         string
	Settlement       string
	Street           string
	HouseNumber      string
	FlatNumber       string
}

type VideoHistory struct {
	Id        int64
	UerId     int64
	User      User
	VideoName string
	CreatedAt int64
}

type Role struct {
	Id   int
	Name string
}

type JobRole struct {
	Id      int
	Role_id int
	Name    string
}

type SettlementType struct {
	Id   int
	Name string
}
