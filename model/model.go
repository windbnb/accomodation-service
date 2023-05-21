package model

import (
	"time"

	"github.com/jinzhu/gorm"
)

type Accomodation struct {
	gorm.Model
	Name                  string
	Address               string
	HasWifi               bool
	HasKitchen            bool
	HasAirConditioning    bool
	HasFreeParking        bool
	MinimimGuests         uint
	MaximumGuests         uint
	Images                []AccomodationImage
	UserId                uint
	Prices                []Price
	AcceptReservationType AcceptReservationType
}

type PriceType string

const (
	PER_GUEST             PriceType = "PER GUEST"
	PER_ACCOMODATION_UNIT PriceType = "PER ACCOMODATION UNIT"
)

type PriceDuration string

const (
	REGULAR PriceDuration = "REGULAR"
	WEEKEND PriceDuration = "WEEKEND"
	HOLIDAY PriceDuration = "HOLIDAY"
)

type Price struct {
	gorm.Model
	StartDate      time.Time
	EndDate        time.Time
	Value          float32
	PriceType      PriceType
	PriceDuration  PriceDuration
	AccomodationID uint
}

type ReservedTerm struct {
	gorm.Model
	StartDate      time.Time
	EndDate        time.Time
	AccomodationID uint
}

type AvailableTerm struct {
	gorm.Model
	StartDate      time.Time
	EndDate        time.Time
	AccomodationID uint
}

func (accomodation *Accomodation) ToDTO() AccomodationDTO {
	return AccomodationDTO{Id: accomodation.ID,
		Name:                  accomodation.Name,
		Address:               accomodation.Address,
		HasWifi:               accomodation.HasWifi,
		HasKitchen:            accomodation.HasKitchen,
		HasAirConditioning:    accomodation.HasAirConditioning,
		HasFreeParking:        accomodation.HasFreeParking,
		MinimimGuests:         accomodation.MinimimGuests,
		MaximumGuests:         accomodation.MaximumGuests,
		Images:                []string{},
		UserId:                accomodation.UserId,
		AcceptReservationType: accomodation.AcceptReservationType}
}

type AccomodationImage struct {
	gorm.Model
	ImageName      string
	AccomodationID uint
}

type AcceptReservationType string

const (
	MANUAL        AcceptReservationType = "MANUAL"
	AUTOMATICALLY AcceptReservationType = "AUTOMATICALLY"
)

func (price *Price) ToDTO() PriceDTO {
	return PriceDTO{Id: price.ID,
		StartDate:      price.StartDate,
		EndDate:        price.EndDate,
		Value:          price.Value,
		PriceType:      price.PriceType,
		PriceDuration:  price.PriceDuration,
		AccomodationID: price.AccomodationID}
}

func (availableTerm *AvailableTerm) ToDTO() AvailableTermDTO {
	return AvailableTermDTO{Id: availableTerm.ID,
		StartDate:      availableTerm.StartDate,
		EndDate:        availableTerm.EndDate,
		AccomodationID: availableTerm.AccomodationID}
}

func (reservedTerm *ReservedTerm) ToDTO() ReservedTermDTO {
	return ReservedTermDTO{Id: reservedTerm.ID,
		StartDate:      reservedTerm.StartDate,
		EndDate:        reservedTerm.EndDate,
		AccomodationID: reservedTerm.AccomodationID}
}
