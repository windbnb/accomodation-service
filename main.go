package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/windbnb/accomodation-service/handler"
	"github.com/windbnb/accomodation-service/repository"
	"github.com/windbnb/accomodation-service/router"
	"github.com/windbnb/accomodation-service/service"
	"github.com/windbnb/accomodation-service/util"
)

func main() {
	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	db := util.ConnectToDatabase()
	router := router.ConfigureRouter(&handler.Handler{Service: &service.AccomodationService{Repo: &repository.Repository{Db: db}}})

	srv := &http.Server{Addr: "0.0.0.0:8082", Handler: router}
	go func() {
		log.Println("server starting")
		if err := srv.ListenAndServe(); err != nil {
			if err != http.ErrServerClosed {
				log.Fatal(err)
			}
		}
	}()

	<-quit

	defer db.Close()
	log.Println("service shutting down ...")

	// gracefully stop server
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal(err)
	}
	log.Println("server stopped")
}
