package controller

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"user_api/service"
)

var _ LoginController = &loginController{}

type LoginController interface {
	Login(rw http.ResponseWriter, r *http.Request)
}

type loginController struct {
	log     *log.Logger
	userSrv service.UserService
}

func NewUserController(userSrv service.UserService, logger *log.Logger) LoginController {
	return &loginController{
		log:     logger,
		userSrv: userSrv,
	}
}

func (l *loginController) Login(rw http.ResponseWriter, r *http.Request) {

	rBody := &struct {
		UserName string `json:"user_name"`
		Password string `json:"password"`
	}{}

	if ok := ParseJSON(rw, r.Body, rBody); !ok {
		return
	}

	userIns, err := l.userSrv.Login(rBody.UserName, rBody.Password)
	if err != nil || userIns == nil {
		RenderJson(rw, http.StatusBadRequest, struct {
			Message string `json:"message"`
		}{
			"Invalid credentials",
		})
		return
	}

	RenderJson(rw, http.StatusOK, struct {
		Message string `json:"message"`
	}{
		fmt.Sprintf("Hey %s, Welcome!", userIns.Name),
	})
	return
}

func RenderJson(w http.ResponseWriter, status int, res interface{}) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)

	// We don't have to write body, If status code is 204 (No Content)
	if status == http.StatusNoContent {
		return
	}

	if err := json.NewEncoder(w).Encode(res); err != nil {
		log.Printf("ERROR: renderJson - %q\n", err)
	}
}

func ParseJSON(w http.ResponseWriter, params io.ReadCloser, data interface{}) bool {
	if params != nil {
		defer params.Close()
	}

	err := json.NewDecoder(params).Decode(data)
	if err == nil {
		return true
	}

	e := struct {
		Message string `json:"message"`
		Err     string  `json:"err"`
	}{
		Message: "Invalid JSON",
		Err:     err.Error(),
	}

	RenderJson(w, http.StatusBadRequest, e)
	return false
}
