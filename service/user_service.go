package service

import (
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"log"
	"math/rand"
	"time"
	"user_api/model"
	"user_api/repository"
)

var _ UserService = &userService{}
var letterRunes = []rune("ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")

type UserService interface {
	Login(username, password string) (*model.User, error)
	Register(user *model.User) error
	List() ([]*model.User, error)
}

type userService struct {
	log  *log.Logger
	repo repository.UserRepo
}

func NewUserService(repo repository.UserRepo, logger *log.Logger) UserService {
	return &userService{
		log:  logger,
		repo: repo,
	}
}

func (u *userService) List() ([]*model.User, error) {
	return u.repo.UserList()
}

func (u *userService) Register(user *model.User) error {
	if len(user.PasswordStr) < 10 {
		return errors.New("please enter valid password, password must be 10 chars long")
	}

	buf := []byte(user.PasswordStr)

	encrypted := sha256.New()
	encrypted.Write(buf)

	encPass := base64.StdEncoding.EncodeToString(encrypted.Sum(nil))

	now := time.Now().UTC().Format(time.RFC3339)
	user.Password = encPass
	user.CreatedTime = now
	user.UpdatedTime = now

	return u.repo.SaveNewUser(user)
}

func (u *userService) Login(username, password string) (*model.User, error) {
	buf := []byte(password)

	encrypted := sha256.New()
	encrypted.Write(buf)

	encPass := base64.StdEncoding.EncodeToString(encrypted.Sum(nil))

	rand.Seed(time.Now().UnixNano())

	b := make([]rune, 36)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]

	}

	sessionToken := string(b)
	userIns, err := u.repo.GetUser(username, encPass)
	if err != nil || userIns == nil{
		log.Printf("error: service - Login - %s", err.Error())
		return nil, err
	}
	userIns.Token = sessionToken
	return userIns, nil
}
