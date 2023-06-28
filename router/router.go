package router

import (
	"github.com/gorilla/mux"
	"github.com/windbnb/accomodation-service/handler"
	"github.com/windbnb/accomodation-service/metrics"
)

func ConfigureRouter(handler *handler.Handler) *mux.Router {
	router := mux.NewRouter()
	router.HandleFunc("/api/accomodation/create", metrics.MetricProxy(handler.CreateAccomodation)).Methods("POST")
	router.HandleFunc("/api/accomodation/{id}", metrics.MetricProxy(handler.FindAccommodationById)).Methods("GET")
	router.HandleFunc("/api/accomodation/{id}/acceptReservationType", metrics.MetricProxy(handler.UpdateAccommodationAcceptReservationType)).Methods("PUT")
	router.HandleFunc("/api/accomodation/search/available", metrics.MetricProxy(handler.SearchAccomodation)).Methods("POST")
	router.HandleFunc("/api/accomodation/for-host/{hostId}", metrics.MetricProxy(handler.FindAccommodationsForHost)).Methods("GET")

	router.HandleFunc("/api/accomodation/image/{filename}", handler.ImageHandler).Methods("GET")

	router.HandleFunc("/api/accomodation/delete-all/{hostId}", metrics.MetricProxy(handler.DeleteHostAccomodation)).Methods("DELETE")
	router.HandleFunc("/api/accomodation/price", metrics.MetricProxy(handler.CreatePrice)).Methods("POST")
	router.HandleFunc("/api/accomodation/price/{id}", metrics.MetricProxy(handler.UpdatePrice)).Methods("PUT")
	router.HandleFunc("/api/accomodation/price/{id}", metrics.MetricProxy(handler.DeletePrice)).Methods("DELETE")
	router.HandleFunc("/api/accomodation/price/for-accomodation/{id}", metrics.MetricProxy(handler.GetPricesForAccomodation)).Methods("GET")

	router.HandleFunc("/api/accomodation/availableTerm", metrics.MetricProxy(handler.CreateAvailableTerm)).Methods("POST")
	router.HandleFunc("/api/accomodation/availableTerm/{id}", metrics.MetricProxy(handler.UpdateAvailableTerm)).Methods("PUT")
	router.HandleFunc("/api/accomodation/availableTerm/{id}", metrics.MetricProxy(handler.DeleteAvailableTerm)).Methods("DELETE")
	router.HandleFunc("/api/accomodation/availableTerm/for-accomodation/{id}", metrics.MetricProxy(handler.GetAvailableTermsForAccomodation)).Methods("GET")

	router.HandleFunc("/api/accomodation/reservedTerm", metrics.MetricProxy(handler.CreateReservedTerm)).Methods("POST")
	router.HandleFunc("/api/accomodation/reservedTerm/{id}", metrics.MetricProxy(handler.DeleteReservedTerm)).Methods("DELETE")

	router.Path("/metrics").Handler(metrics.MetricsHandler())

	router.HandleFunc("/probe/liveness", handler.Healthcheck)
	router.HandleFunc("/probe/readiness", handler.Ready)

	return router
}
