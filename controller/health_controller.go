package controller

import (
	"log"
	"net/http"
)

var _ HealthController = &healthController{}

type HealthController interface {
	Ping(rw http.ResponseWriter, r *http.Request)
}

type healthController struct {
	log *log.Logger
}

func (h *healthController) Ping(rw http.ResponseWriter, r *http.Request) {
	h.log.Println("Info: controller - Ping")
	rw.WriteHeader(http.StatusOK)
	rw.Write([]byte("pong!."))
	return

}

func NewHealthController(logger *log.Logger) HealthController {
	return &healthController{
		log: logger,
	}
}
