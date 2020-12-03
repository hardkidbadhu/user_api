package service

import (
	"crypto/sha256"
	"encoding/base64"
	"log"
	"user_api/model"
	"user_api/repository"
)

var _ UserService = &userService{}

type UserService interface {
	Login(username, password string) (*model.User, error)
}

type userService struct {
	log  *log.Logger
	repo repository.UserRepo
}

func NewUserService(repo repository.UserRepo, logger *log.Logger) UserService {
	return &userService{
		log: logger,
		repo: repo,
	}
}

func (u *userService) Login(username, password string) (*model.User, error) {
	buf := []byte(password)

	encrypted := sha256.New()
	encrypted.Write(buf)

	encPass := base64.StdEncoding.EncodeToString(encrypted.Sum(nil))

	return u.repo.GetUser(username, encPass)
}


