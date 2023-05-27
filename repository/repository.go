package repository

import (
	"errors"
	"strconv"
	"time"

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

func (r *Repository) SavePrice(price model.Price) model.Price {
	r.Db.Create(&price)
	return price
}

func (r *Repository) SaveAvailableTerm(availableTerm model.AvailableTerm) model.AvailableTerm {
	r.Db.Create(&availableTerm)
	return availableTerm
}

func (r *Repository) SaveReservedTerm(reservedTerm model.ReservedTerm) model.ReservedTerm {
	r.Db.Create(&reservedTerm)
	return reservedTerm
}

func (r *Repository) UpdatePrice(price model.Price) model.Price {
	r.Db.Save(&price)
	return price
}

func (r *Repository) UpdateAvailableTerm(availableTerm model.AvailableTerm) model.AvailableTerm {
	r.Db.Save(&availableTerm)
	return availableTerm
}

func (r *Repository) UpdateAccommodation(accommodation model.Accomodation) model.Accomodation {
	r.Db.Save(&accommodation)
	return accommodation
}

func (r *Repository) FindAccomodationById(id uint) (model.Accomodation, error) {
	var accomodation model.Accomodation

	r.Db.First(&accomodation, id)

	if accomodation.ID == 0 {
		return model.Accomodation{}, errors.New("there is no accomodation with id " + strconv.FormatUint(uint64(id), 10))
	}

	return accomodation, nil
}

func (r *Repository) FindPriceById(id uint64) (model.Price, error) {
	var price model.Price

	r.Db.First(&price, id)

	if price.ID == 0 {
		return model.Price{}, errors.New("there is no price with id " + strconv.FormatUint(uint64(id), 10))
	}

	return price, nil
}

func (r *Repository) FindAvailableTermById(id uint64) (model.AvailableTerm, error) {
	var availableTerm model.AvailableTerm

	r.Db.First(&availableTerm, id)

	if availableTerm.ID == 0 {
		return model.AvailableTerm{}, errors.New("there is no available term with id " + strconv.FormatUint(uint64(id), 10))
	}

	return availableTerm, nil
}

func (r *Repository) FindAvailableTermAfter(accommodationId uint, after time.Time) []model.AvailableTerm {
	availableTerms := &[]model.AvailableTerm{}

	r.Db.Where("accomodation_id = ? and (start_date <= ? or end_date <= ?)", accommodationId, after, after).Find(availableTerms)

	return *availableTerms
}

func (r *Repository) FindReservedTermById(id uint64) (model.ReservedTerm, error) {
	var reservedTerm model.ReservedTerm

	r.Db.First(&reservedTerm, id)

	if reservedTerm.ID == 0 {
		return model.ReservedTerm{}, errors.New("there is no reserved term with id " + strconv.FormatUint(uint64(id), 10))
	}

	return reservedTerm, nil
}

func (r *Repository) DeletePrice(id uint64) error {
	var price model.Price

	r.Db.First(&price, id)

	if price.ID == 0 {
		return errors.New("there is no price with id " + strconv.FormatUint(uint64(id), 10))
	}

	r.Db.Delete(&model.Price{}, id)
	return nil
}

func (r *Repository) DeleteAvailableTerm(id uint64) error {
	var availableTerm model.AvailableTerm

	r.Db.First(&availableTerm, id)

	if availableTerm.ID == 0 {
		return errors.New("there is no available term with id " + strconv.FormatUint(uint64(id), 10))
	}

	r.Db.Delete(&model.AvailableTerm{}, id)
	return nil
}

func (r *Repository) DeleteReservedTerm(id uint64) error {
	var reservedTerm model.ReservedTerm

	r.Db.First(&reservedTerm, id)

	if reservedTerm.ID == 0 {
		return errors.New("there is no reserved term with id " + strconv.FormatUint(uint64(id), 10))
	}

	r.Db.Delete(&model.ReservedTerm{}, id)
	return nil
}

func (r *Repository) FindAccomodationByGuestsAndAddress(numberOfGuests uint, address string) []model.Accomodation {
	accomodations := &[]model.Accomodation{}

	r.Db.Find(&accomodations, "address LIKE ? AND minimim_guests <= ? AND maximum_guests >= ?", "%"+address+"%", numberOfGuests, numberOfGuests)

	return *accomodations
}

func (r *Repository) IsReserved(accomodationId uint, startDate time.Time, endDate time.Time) bool {
	count := int64(0)

	r.Db.Model(&model.ReservedTerm{}).Where("accomodation_id = ? AND start_date <= ? AND end_date >= ?", accomodationId, endDate, startDate).Count(&count)

	if count > 0 {
		return true
	}
	return false
}

func (r *Repository) IsAvailable(accomodationId uint, startDate time.Time, endDate time.Time) bool {
	count := int64(0)

	r.Db.Model(&model.AvailableTerm{}).Where("accomodation_id = ? AND start_date <= ? AND end_date >= ?", accomodationId, endDate, startDate).Count(&count)

	if count > 0 {
		return true
	}
	return false
}

func (r *Repository) FindPricesForAccomodation(accomodationId uint, startDate time.Time, endDate time.Time) []model.Price {
	prices := &[]model.Price{}

	r.Db.Find(&prices, "accomodation_id = ? AND start_date <= ? AND end_date >= ? AND active = true", accomodationId, endDate, startDate)

	return *prices
}
