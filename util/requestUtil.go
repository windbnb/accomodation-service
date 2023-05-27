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
	priceType := r.MultipartForm.Value["priceType"][0]

	defaultAcceptReservationType := model.MANUAL

	return model.Accomodation{
		Name:                  name,
		Address:               address,
		HasWifi:               hasWifi,
		HasKitchen:            hasKitchen,
		HasAirConditioning:    hasAirConditioning,
		HasFreeParking:        hasFreeParking,
		MinimimGuests:         uint(minimimGuests),
		MaximumGuests:         uint(maximumGuests),
		UserId:                0,
		AcceptReservationType: defaultAcceptReservationType,
    PriceType:             model.PriceType(priceType)}
}

func FromCreatePriceDTOToPrice(price model.CreatePriceDTO) model.Price {

	return model.Price{
		StartDate:      price.StartDate,
		EndDate:        price.EndDate,
		Value:          price.Value,
		PriceDuration:  price.PriceDuration,
		AccomodationID: price.AccomodationID}
}

func FromUpdatePriceDTOToPrice(price model.UpdatePriceDTO) model.Price {

	return model.Price{
		StartDate: price.StartDate,
		EndDate:   price.EndDate,
		Value:     price.Value}
}

func FromCreateAvailableTermDTOToAvailableTerm(availableTerm model.CreateAvailableTermDTO) model.AvailableTerm {

	return model.AvailableTerm{
		StartDate:      availableTerm.StartDate,
		EndDate:        availableTerm.EndDate,
		AccomodationID: availableTerm.AccomodationID}
}

func FromUpdateAvailableTermDTOToAvailableTerm(availableTerm model.UpdateAvailableTermDTO) model.AvailableTerm {

	return model.AvailableTerm{
		StartDate: availableTerm.StartDate,
		EndDate:   availableTerm.EndDate}
}

func FromCreateReservedTermDTOToReservedTerm(reservedTerm model.CreateReservedTermDTO) model.ReservedTerm {

	return model.ReservedTerm{
		StartDate:      reservedTerm.StartDate,
		EndDate:        reservedTerm.EndDate,
		AccomodationID: reservedTerm.AccomodationID}
}
