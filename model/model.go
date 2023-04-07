package model

import (
	"github.com/jinzhu/gorm"
)

type Accomodation struct {
	gorm.Model
	Name string
	Address string
	HasWifi bool
	HasKitchen bool
	HasAirConditioning bool
	HasFreeParking bool
	MinimimGuests uint
	MaximumGuests uint
	Images []AccomodationImage
}

func (accomodation *Accomodation) ToDTO() AccomodationDTO {
	return AccomodationDTO{Id: accomodation.ID, Name: accomodation.Name, Address: accomodation.Address, HasWifi: accomodation.HasWifi, HasKitchen: accomodation.HasKitchen, HasAirConditioning: accomodation.HasAirConditioning, HasFreeParking: accomodation.HasFreeParking, MinimimGuests: accomodation.MinimimGuests, MaximumGuests: accomodation.MaximumGuests, Images: []string{}}
}

type AccomodationImage struct {
	gorm.Model
	ImageName string
	AccomodationID uint
}
