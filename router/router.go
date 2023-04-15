package router

import (
	"github.com/gorilla/mux"
	"github.com/windbnb/accomodation-service/handler"
)

func ConfigureRouter(handler *handler.Handler) *mux.Router {
	router := mux.NewRouter()
	router.HandleFunc("/api/accomodation", handler.CreateAccomodation).Methods("POST")

	router.HandleFunc("/api/accomodation/image/{filename}", handler.ImageHandler).Methods("GET")

	router.HandleFunc("/api/accomodation/delete-all/{hostId}", handler.DeleteHostAccomodation).Methods("DELETE")

	return router
}