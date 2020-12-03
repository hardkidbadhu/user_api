package repository

import (
	"database/sql"
	"fmt"
	"log"
	"user_api/model"
)

const (
	FindUser = `select * from user_details where user_name=? AND password=?`
	InsertUser = `INSERT INTO user_details(user_name,password,name,age,gender,phone_number,email_id,created_time,updated_time) 
             VALUES(?,?,?,?,?,?,?,?,?) `
	FindAllUser = `select * from user_details`
)

var _ UserRepo = &userRepo{}

type UserRepo interface {
	SaveNewUser(user *model.User) error
	GetUser(username, password string) (*model.User, error)
	UserList() ([]*model.User, error)
}

type userRepo struct {
	db  *sql.DB
	log *log.Logger
}

func (u *userRepo) SaveNewUser(user *model.User) error {
	res, err := u.db.Exec(InsertUser, user.UserName, user.Password, user.Name, user.Age,
		user.Gender, user.PhoneNumber, user.EmailID, user.CreatedTime, user.UpdatedTime)
	if err != nil {
		u.log.Printf("error: SaveNewUser - GetUser - %s", err.Error())
		return err
	}

	id, _ := res.LastInsertId()

	u.log.Println("info - inserted id", id)
	return nil
}

func (u *userRepo) UserList() ([]*model.User, error) {
	rows, err := u.db.Query(FindAllUser)
	if err != nil {
		u.log.Printf("error: list all available users  - UserList - %s", err.Error())
		return nil, err
	}

	defer rows.Close()
	result := []*model.User{}

	for rows.Next() {
		user := new(model.User)

		if err := rows.Scan(
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
			u.log.Printf("error: repository - UserList(unmarshalling scanned details) - %s", rows.Err())
			return nil, rows.Err()
		}

		result = append(result, user)
	}

	fmt.Println(result)
	return result, nil
}


func (u *userRepo) GetUser(username, password string) (*model.User, error) {
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

func NewUserRepo(db  *sql.DB, log *log.Logger) UserRepo {
	return &userRepo{
		db: db,
		log: log,
	}
}