package router

import (
	"database/sql"
	"fmt"
	"time"
	"user_api/constants"
	"user_api/controller"
	"user_api/middleware"
	"user_api/repository"
	"user_api/service"

	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
)

func NewRouter(logger *log.Logger) *http.ServeMux {
	mux := http.NewServeMux()


	source := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", constants.DBUser, constants.DBPassword, "127.0.0.1",
		"3306", constants.DBName)
	db, err := sql.Open("mysql", source)
	if err != nil {
		logger.Fatalf("error: db connection")
	}

	db.SetConnMaxLifetime(time.Minute * 3)
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(10)

	pingCtrl := controller.NewHealthController(logger)

	userRepo := repository.NewUserRepo(db, logger)
	srv := service.NewUserService(userRepo, logger)
	loginCtrl := controller.NewUserController(srv, logger)

	mux.Handle("/api/v1/ping", middleware.LoggingHandler(middleware.RecoverHandler(http.HandlerFunc(pingCtrl.Ping))))
	mux.Handle("/api/v1/login", middleware.LoggingHandler(middleware.RecoverHandler(http.HandlerFunc(loginCtrl.Login))))
	mux.Handle("/api/v1/register", middleware.LoggingHandler(middleware.RecoverHandler(http.HandlerFunc(loginCtrl.Register))))
	mux.Handle("/api/v1/list/users", middleware.LoggingHandler(middleware.RecoverHandler(middleware.PostLogin(http.HandlerFunc(loginCtrl.ListAllUsers)))))

	return mux
}
