package repository

import (
	"database/sql"
	"log"
	"user_api/model"
)

const (
	FindUser = `select * from user_details where user_name=? AND password=?`
)

var _ UserRepo = &userRepo{}

type UserRepo interface {
	GetUser(username, password string) (*model.User, error)
}

type userRepo struct {
	db  *sql.DB
	log *log.Logger
}

func NewUserRepo(db  *sql.DB, log *log.Logger) UserRepo {
	return &userRepo{
		db: db,
		log: log,
	}
}

func (u userRepo) GetUser(username, password string) (*model.User, error) {
	user := new(model.User)

	row := u.db.QueryRow(FindUser, username, password)
	if row.Err() != nil {
		u.log.Printf("error: repository - GetUser - %s", row.Err())
		return nil, row.Err()
	}

	if err := row.Scan(
		&user.Id,
		&user.UserName,
		&user.Password,
		&user.Name,
		&user.Age,
		&user.Gender,
		&user.PhoneNumber,
		&user.EmailID,
		&user.CreatedTime,
		&user.UpdatedTime,
		); err != nil {
		u.log.Printf("error: repository - GetUser(unmarshalling scanned details) - %s", row.Err())
		return nil, row.Err()
	}

	log.Println("user", user)
	return user, nil
}
