package model

type AccomodationDTO struct {
	Id                 uint     `json:"id"`
	Name               string   `json:"statusCode"`
	Address            string   `json:"address"`
	HasWifi            bool     `json:"hasWifi"`
	HasKitchen         bool     `json:"hasKitchen"`
	HasAirConditioning bool     `json:"hasAirConditioning"`
	HasFreeParking     bool     `json:"hasFreeParking"`
	MinimimGuests      uint     `json:"minimimGuests"`
	MaximumGuests      uint     `json:"maximumGuests"`
	Images             []string `json:"images"`
}

type ErrorResponse struct {
	Message    string `json:"message"`
	StatusCode int    `json:"statusCode"`
}
