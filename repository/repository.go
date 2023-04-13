package repository

import (
	"github.com/jinzhu/gorm"
	"github.com/windbnb/accomodation-service/model"
)

type Repository struct {
	Db *gorm.DB
}

func (r *Repository) SaveAccomodation(accomodation model.Accomodation) model.Accomodation {
	r.Db.Create(&accomodation)
	return accomodation
}

func (r *Repository) SaveAccomodationImage(image model.AccomodationImage) model.AccomodationImage {
	r.Db.Create(&image)
	return image
}
