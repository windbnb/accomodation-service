package router

import (
	"github.com/gorilla/mux"
	"github.com/windbnb/accomodation-service/handler"
)

func ConfigureRouter(handler *handler.Handler) *mux.Router {
	router := mux.NewRouter()
	router.HandleFunc("/api/accomodation", handler.CreateAccomodation).Methods("POST")
	router.HandleFunc("/api/accomodation/{id}", handler.FindAccommodationById).Methods("GET")
	router.HandleFunc("/api/accomodation/{id}/acceptReservationType", handler.UpdateAccommodationAcceptReservationType).Methods("PUT")

	router.HandleFunc("/api/accomodation/image/{filename}", handler.ImageHandler).Methods("GET")

	router.HandleFunc("/api/accomodation/delete-all/{hostId}", handler.DeleteHostAccomodation).Methods("DELETE")
	router.HandleFunc("/api/accomodation/price", handler.CreatePrice).Methods("POST")
	router.HandleFunc("/api/accomodation/price/{id}", handler.UpdatePrice).Methods("PUT")
	router.HandleFunc("/api/accomodation/price/{id}", handler.DeletePrice).Methods("DELETE")

	router.HandleFunc("/api/accomodation/availableTerm", handler.CreateAvailableTerm).Methods("POST")
	router.HandleFunc("/api/accomodation/availableTerm/{id}", handler.UpdateAvailableTerm).Methods("PUT")
	router.HandleFunc("/api/accomodation/availableTerm/{id}", handler.DeleteAvailableTerm).Methods("DELETE")

	router.HandleFunc("/api/accomodation/reservedTerm", handler.CreateReservedTerm).Methods("POST")
	router.HandleFunc("/api/accomodation/reservedTerm/{id}", handler.DeleteReservedTerm).Methods("DELETE")

	return router
}
