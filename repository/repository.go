package repository

import (
	"errors"

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

func (r *Repository) DeleteHostAccomodation(hostId uint) error {
	accomodationIdsSubQuery := r.Db.Table("accomodations").Where("user_id = ?", hostId).Select("id").SubQuery()
	if err := r.Db.Where("accomodation_id IN (?)", accomodationIdsSubQuery).Delete(&model.AccomodationImage{}).Error; err != nil {
		return err
	}

	result := r.Db.Where("user_id = ?", hostId).Delete(&model.Accomodation{})
	if result.Error != nil {
		return result.Error
	} else if result.RowsAffected == 0 {
		return errors.New("there are no accomodations for host with given id")
	}

	return nil
}
