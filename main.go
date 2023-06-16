package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/opentracing/opentracing-go"
	"github.com/rs/cors"
	"github.com/windbnb/accomodation-service/handler"
	"github.com/windbnb/accomodation-service/repository"
	"github.com/windbnb/accomodation-service/router"
	"github.com/windbnb/accomodation-service/service"
	"github.com/windbnb/accomodation-service/tracer"
	"github.com/windbnb/accomodation-service/util"
)

func main() {
	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	db := util.ConnectToDatabase()

	tracer, closer := tracer.Init("accomodation-service")
	opentracing.SetGlobalTracer(tracer)
	router := router.ConfigureRouter(&handler.Handler{
		Tracer:  tracer,
		Closer:  closer,
		Service: &service.AccomodationService{Repo: &repository.Repository{Db: db}}})

	servicePath, servicePathFound := os.LookupEnv("SERVICE_PATH")
	if !servicePathFound {
		servicePath = "http://localhost:8082"
	}

	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:3005"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowCredentials: true,
		Debug:            true,
		AllowedHeaders:   []string{"Accept", "Content-Type", "Content-Length", "Accept-Encoding", "X-CSRF-Token", "Authorization"},
	})

	srv := &http.Server{Addr: servicePath, Handler: c.Handler(router)}

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
