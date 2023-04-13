package util

import (
	"net/http"
	"net/url"
	"strconv"

	roundrobin "github.com/hlts2/round-robin"
	"github.com/windbnb/accomodation-service/model"
)

var BaseUserServicePathRoundRobin, _ = roundrobin.New(
    &url.URL{Host: "http://localhost:8081"},
)

func ParseMultipartAccomodation(r *http.Request) model.Accomodation {
	name := r.MultipartForm.Value["name"][0]
	address := r.MultipartForm.Value["address"][0]
	hasWifi, _ := strconv.ParseBool(r.MultipartForm.Value["hasWifi"][0])
	hasKitchen, _ := strconv.ParseBool(r.MultipartForm.Value["hasKitchen"][0])
	hasAirConditioning, _ := strconv.ParseBool(r.MultipartForm.Value["hasAirConditioning"][0])
	hasFreeParking, _ := strconv.ParseBool(r.MultipartForm.Value["hasFreeParking"][0])
	minimimGuests, _ := strconv.ParseUint(r.MultipartForm.Value["minimumGuests"][0], 10, 32)
	maximumGuests, _ := strconv.ParseUint(r.MultipartForm.Value["maximumGuests"][0], 10, 32)
	
	return model.Accomodation{Name: name, Address: address, HasWifi: hasWifi, HasKitchen: hasKitchen, HasAirConditioning: hasAirConditioning, HasFreeParking: hasFreeParking, MinimimGuests: uint(minimimGuests), MaximumGuests: uint(maximumGuests), UserId: 0}
}
