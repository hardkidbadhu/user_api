package model

import "time"

type User struct {
	Id          int
	UserName    string
	Password    string
	Name        string
	Age         int
	Gender      string
	PhoneNumber string
	EmailID     string
	CreatedTime *time.Time
	UpdatedTime *time.Time
}
