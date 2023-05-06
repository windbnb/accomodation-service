package model

import (
	"time"
)

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
	UserId             uint     `json:"userId"`
}

type ErrorResponse struct {
	Message    string `json:"message"`
	StatusCode int    `json:"statusCode"`
}

type UserRole string

const (
	HOST  UserRole = "HOST"
	GUEST UserRole = "GUEST"
)

type UserResponseDTO struct {
	Id       uint     `json:"id"`
	Email    string   `json:"email"`
	Name     string   `json:"name"`
	Surname  string   `json:"surname"`
	Address  string   `json:"address"`
	Username string   `json:"username"`
	Role     UserRole `json:"role"`
}

type CreatePriceDTO struct {
	StartDate      time.Time     `json:"startDate"`
	EndDate        time.Time     `json:"endDate"`
	Value          float32       `json:"value"`
	PriceType      PriceType     `json:"priceType"`
	PriceDuration  PriceDuration `json:"priceDuration"`
	AccomodationID uint          `json:"accomodationId"`
}

type UpdatePriceDTO struct {
	Id        uint      `json:"id"`
	StartDate time.Time `json:"startDate"`
	EndDate   time.Time `json:"endDate"`
	Value     float32   `json:"value"`
}

type PriceDTO struct {
	Id             uint          `json:"id"`
	StartDate      time.Time     `json:"startDate"`
	EndDate        time.Time     `json:"endDate"`
	Value          float32       `json:"value"`
	PriceType      PriceType     `json:"priceType"`
	PriceDuration  PriceDuration `json:"priceDuration"`
	AccomodationID uint          `json:"accomodationId"`
}

type CreateAvailableTermDTO struct {
	StartDate      time.Time `json:"startDate"`
	EndDate        time.Time `json:"endDate"`
	AccomodationID uint      `json:"accomodationId"`
}

type UpdateAvailableTermDTO struct {
	Id        uint      `json:"id"`
	StartDate time.Time `json:"startDate"`
	EndDate   time.Time `json:"endDate"`
}

type AvailableTermDTO struct {
	Id             uint      `json:"id"`
	StartDate      time.Time `json:"startDate"`
	EndDate        time.Time `json:"endDate"`
	AccomodationID uint      `json:"accomodationId"`
}

type CreateReservedTermDTO struct {
	StartDate      time.Time `json:"startDate"`
	EndDate        time.Time `json:"endDate"`
	AccomodationID uint      `json:"accomodationId"`
}

type ReservedTermDTO struct {
	Id             uint      `json:"id"`
	StartDate      time.Time `json:"startDate"`
	EndDate        time.Time `json:"endDate"`
	AccomodationID uint      `json:"accomodationId"`
}
