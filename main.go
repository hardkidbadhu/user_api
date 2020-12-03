package main

import (
	"user_api/constants"
	"user_api/middleware"
	"user_api/model"
	"user_api/router"

	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func init() {
	log.Println("initiating cache")
	middleware.CacheMap = make(map[string]model.User)
}

func main() {

	logger := log.New(os.Stdout, "user_service: ", log.LstdFlags|log.Lshortfile)
	logger.SetFlags(log.LstdFlags | log.Lshortfile)

	//init router
	mux := router.NewRouter(logger)

	srv := http.Server{
		Addr: constants.Host + ":" + constants.Port,
		Handler: mux,
		ReadTimeout:    time.Duration(60) * time.Second,
		WriteTimeout:   time.Duration(80) * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	// Graceful shut down of server
	graceful := make(chan os.Signal)
	signal.Notify(graceful, syscall.SIGINT)
	signal.Notify(graceful, syscall.SIGTERM)
	go func() {
		<-graceful
		logger.Println("Shutting down server...")
		ctx, cancelFunc := context.WithTimeout(context.Background(), 20 * time.Second)
		defer cancelFunc()
		if err := srv.Shutdown(ctx); err != nil {
			logger.Fatalf("Could not do graceful shutdown: %v\n", err)
		}
	}()

	logger.Println("Listening server on ", constants.Host + ":" + constants.Port)
	if err := srv.ListenAndServe(); err != http.ErrServerClosed {
		logger.Fatalf("server: fatal error - %s", err.Error())
	}

}
