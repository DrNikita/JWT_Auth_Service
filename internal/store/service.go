package store

import (
	"context"
	"database/sql"
	"log/slog"
)

type StoreService struct {
	db     *sql.DB
	logger *slog.Logger
	ctx    *context.Context
}

func NewDbService(db *sql.DB, logger *slog.Logger, ctx *context.Context) *StoreService {
	return &StoreService{
		db:     db,
		logger: logger,
		ctx:    ctx,
	}
}

func (ss *StoreService) CreateAddress(address *Address) (int64, error) {
	tx, err := ss.db.Begin()
	if err != nil {
		return 0, err
	}
	defer tx.Rollback()

	var addressID int64
	sqlStatement := `
		INSERT INTO public.address
		(settlement_type_id, country, region, district, settlement, street, house_number, flat_number)
		VALUES($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING id
	`
	err = tx.QueryRowContext(*ss.ctx, sqlStatement,
		address.SettlementTypeId, address.Country, address.Region, address.District,
		address.Settlement, address.Street, address.HouseNumber, address.FlatNumber).
		Scan(&addressID)

	if err != nil {
		return addressID, err
	}

	err = tx.Commit()
	if err != nil {
		return 0, err
	}
	return addressID, nil
}

func (ss *StoreService) CreateUser(user *User) (int64, error) {
	tx, err := ss.db.Begin()
	if err != nil {
		return 0, err
	}
	defer tx.Rollback()

	var userID int64
	sqlStatement := `
		INSERT INTO public."user"
		(job_role_id, address_id, "name", second_name, surname, email, "password", birthday, is_active)
		VALUES($1, $2, $3, $4, $5, $6, $7, $8, $9)
		RETURNING id
	`
	err = tx.QueryRowContext(*ss.ctx, sqlStatement,
		user.JobRoleId, user.AddressId, user.Name, user.SecondName, user.Surname,
		user.Email, user.Password, user.Birthday, user.IsActive).
		Scan(&userID)

	if err != nil {
		return userID, err
	}

	err = tx.Commit()
	if err != nil {
		return 0, err
	}
	return userID, nil
}

func (ss *StoreService) FindUserByEmail(email string) (*User, error) {
	tx, err := ss.db.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	var user User
	sqlStatement := `
		SELECT id, job_role_id, address_id, "name", second_name, surname,
		email, "password", birthday, is_active
		FROM public."user"
		WHERE "user".email = $1
	`
	err = tx.QueryRowContext(*ss.ctx, sqlStatement, email).
		Scan(
			&user.Id, &user.JobRoleId, &user.AddressId, &user.Name, &user.SecondName,
			&user.Surname, &user.Email, &user.Password, &user.Birthday, &user.IsActive,
		)
	if err != nil {
		return &user, err
	}

	err = tx.Commit()
	if err != nil {
		return nil, err
	}
	return &user, nil
}
