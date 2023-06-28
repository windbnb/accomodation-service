package repository

import (
	"context"
	"errors"
	"strconv"
	"strings"
	"time"

	"github.com/jinzhu/gorm"
	"github.com/windbnb/accomodation-service/model"
	"github.com/windbnb/accomodation-service/tracer"
)

type IRepository interface {
	SaveAccomodation(accomodation model.Accomodation, ctx context.Context) model.Accomodation
	SaveAccomodationImage(image model.AccomodationImage) model.AccomodationImage
	DeleteHostAccomodation(hostId uint, ctx context.Context) error
	SavePrice(price model.Price) model.Price
	SaveAvailableTerm(availableTerm model.AvailableTerm, ctx context.Context) model.AvailableTerm
	SaveReservedTerm(reservedTerm model.ReservedTerm, ctx context.Context) model.ReservedTerm
	UpdatePrice(price model.Price, ctx context.Context) model.Price
	UpdateAvailableTerm(availableTerm model.AvailableTerm, ctx context.Context) model.AvailableTerm
	UpdateAccommodation(accommodation model.Accomodation, ctx context.Context) model.Accomodation
	FindAccomodationById(id uint, ctx context.Context) (model.Accomodation, error)
	FindPriceById(id uint64, ctx context.Context) (model.Price, error)
	FindAvailableTermById(id uint64, ctx context.Context) (model.AvailableTerm, error)
	FindAvailableTermAfter(accommodationId uint, after time.Time, ctx context.Context) []model.AvailableTerm
	FindReservedTermById(id uint64, ctx context.Context) (model.ReservedTerm, error)
	DeletePrice(id uint64, ctx context.Context) error
	DeleteAvailableTerm(id uint64) error
	DeleteReservedTerm(id uint64) error
	FindAccomodationByGuestsAndAddress(numberOfGuests uint, address string, ctx context.Context) []model.Accomodation
	IsReserved(accomodationId uint, startDate time.Time, endDate time.Time, ctx context.Context) bool
	IsAvailable(accomodationId uint, startDate time.Time, endDate time.Time, ctx context.Context) bool
	FindPricesForAccomodation(accomodationId uint, startDate time.Time, endDate time.Time) []model.Price
	FindImagesForAccomodation(accomodationId uint) []string
	FindAccomodationsForHost(hostId uint, ctx context.Context) []model.Accomodation
	GetAvailableTermsForAccomodation(accomodationId uint, ctx context.Context) []model.AvailableTerm
	GetPricesForAccomodation(accomodationId uint, ctx context.Context) []model.Price
}

type Repository struct {
	Db *gorm.DB
}

func (r *Repository) SaveAccomodation(accomodation model.Accomodation, ctx context.Context) model.Accomodation {
	span := tracer.StartSpanFromContext(ctx, "saveAccomodationRepository")
	defer span.Finish()
	r.Db.Create(&accomodation)
	return accomodation
}

func (r *Repository) SaveAccomodationImage(image model.AccomodationImage) model.AccomodationImage {
	r.Db.Create(&image)
	return image
}

func (r *Repository) DeleteHostAccomodation(hostId uint, ctx context.Context) error {
	span := tracer.StartSpanFromContext(ctx, "saveAccomodationRepository")
	defer span.Finish()
	accomodationIdsSubQuery := r.Db.Table("accomodations").Where("user_id = ?", hostId).Select("id").SubQuery()
	if err := r.Db.Where("accomodation_id IN (?)", accomodationIdsSubQuery).Delete(&model.AccomodationImage{}).Error; err != nil {
		tracer.LogError(span, err)
		return err
	}

	result := r.Db.Where("user_id = ?", hostId).Delete(&model.Accomodation{})
	if result.Error != nil {
		tracer.LogError(span, result.Error)
		return result.Error
	} else if result.RowsAffected == 0 {
		err := errors.New("there are no accomodations for host with given id")
		tracer.LogError(span, err)
		return err
	}
	return nil
}

func (r *Repository) SavePrice(price model.Price) model.Price {
	r.Db.Create(&price)
	return price
}

func (r *Repository) SaveAvailableTerm(availableTerm model.AvailableTerm, ctx context.Context) model.AvailableTerm {
	span := tracer.StartSpanFromContext(ctx, "saveAvailableTermRepository")
	defer span.Finish()

	r.Db.Create(&availableTerm)
	return availableTerm
}

func (r *Repository) SaveReservedTerm(reservedTerm model.ReservedTerm, ctx context.Context) model.ReservedTerm {
	span := tracer.StartSpanFromContext(ctx, "saveReservedTermRepository")
	defer span.Finish()

	r.Db.Create(&reservedTerm)
	return reservedTerm
}

func (r *Repository) UpdatePrice(price model.Price, ctx context.Context) model.Price {
	span := tracer.StartSpanFromContext(ctx, "updatePriceRepository")
	defer span.Finish()

	r.Db.Save(&price)
	return price
}

func (r *Repository) UpdateAvailableTerm(availableTerm model.AvailableTerm, ctx context.Context) model.AvailableTerm {
	span := tracer.StartSpanFromContext(ctx, "updateAvailableTermRepository")
	defer span.Finish()

	r.Db.Save(&availableTerm)
	return availableTerm
}

func (r *Repository) UpdateAccommodation(accommodation model.Accomodation, ctx context.Context) model.Accomodation {
	span := tracer.StartSpanFromContext(ctx, "updateAccomodationRepository")
	defer span.Finish()
	r.Db.Save(&accommodation)
	return accommodation
}

func (r *Repository) FindAccomodationById(id uint, ctx context.Context) (model.Accomodation, error) {
	span := tracer.StartSpanFromContext(ctx, "findAccomodationByIdRepository")
	defer span.Finish()
	var accomodation model.Accomodation

	r.Db.First(&accomodation, id)

	if accomodation.ID == 0 {
		err := errors.New("there is no accomodation with id " + strconv.FormatUint(uint64(id), 10))
		tracer.LogError(span, err)
		return model.Accomodation{}, err
	}

	return accomodation, nil
}

func (r *Repository) FindPriceById(id uint64, ctx context.Context) (model.Price, error) {
	span := tracer.StartSpanFromContext(ctx, "findPriceByIdRepository")
	defer span.Finish()
	var price model.Price

	r.Db.First(&price, id)

	if price.ID == 0 {
		err := errors.New("there is no price with id " + strconv.FormatUint(uint64(id), 10))
		tracer.LogError(span, err)
		return model.Price{}, err
	}

	return price, nil
}

func (r *Repository) FindAvailableTermById(id uint64, ctx context.Context) (model.AvailableTerm, error) {
	span := tracer.StartSpanFromContext(ctx, "findAvailableTermByIdRepository")
	defer span.Finish()
	var availableTerm model.AvailableTerm

	r.Db.First(&availableTerm, id)

	if availableTerm.ID == 0 {
		err := errors.New("there is no available term with id " + strconv.FormatUint(uint64(id), 10))
		tracer.LogError(span, err)
		return model.AvailableTerm{}, err
	}

	return availableTerm, nil
}

func (r *Repository) FindAvailableTermAfter(accommodationId uint, after time.Time, ctx context.Context) []model.AvailableTerm {
	span := tracer.StartSpanFromContext(ctx, "findAvailableTermAfterRepository")
	defer span.Finish()
	availableTerms := &[]model.AvailableTerm{}

	r.Db.Where("accomodation_id = ? and (start_date <= ? or end_date <= ?)", accommodationId, after, after).Find(availableTerms)

	return *availableTerms
}

func (r *Repository) FindReservedTermById(id uint64, ctx context.Context) (model.ReservedTerm, error) {
	span := tracer.StartSpanFromContext(ctx, "findReservedTermByIdRepository")
	defer span.Finish()
	var reservedTerm model.ReservedTerm

	r.Db.First(&reservedTerm, id)

	if reservedTerm.ID == 0 {
		err := errors.New("there is no reserved term with id " + strconv.FormatUint(uint64(id), 10))
		tracer.LogError(span, err)
		return model.ReservedTerm{}, err
	}

	return reservedTerm, nil
}

func (r *Repository) DeletePrice(id uint64, ctx context.Context) error {
	span := tracer.StartSpanFromContext(ctx, "deletePriceRepository")
	defer span.Finish()
	var price model.Price

	r.Db.First(&price, id)

	if price.ID == 0 {
		err := errors.New("there is no price with id " + strconv.FormatUint(uint64(id), 10))
		tracer.LogError(span, err)
		return err
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

func (r *Repository) FindAccomodationByGuestsAndAddress(numberOfGuests uint, address string, ctx context.Context) []model.Accomodation {
	span := tracer.StartSpanFromContext(ctx, "findAccomodationByGuestsAndAddressRepository")
	defer span.Finish()
	accomodations := &[]model.Accomodation{}

	r.Db.Find(&accomodations, "LOWER(address) LIKE ? AND minimim_guests <= ? AND maximum_guests >= ?", "%"+strings.ToLower(address)+"%", numberOfGuests, numberOfGuests)

	return *accomodations
}

func (r *Repository) IsReserved(accomodationId uint, startDate time.Time, endDate time.Time, ctx context.Context) bool {
	span := tracer.StartSpanFromContext(ctx, "isReservedRepository")
	defer span.Finish()
	count := int64(0)

	r.Db.Model(&model.ReservedTerm{}).Where("accomodation_id = ? AND start_date <= ? AND end_date >= ?", accomodationId, endDate, startDate).Count(&count)

	if count > 0 {
		return true
	}
	return false
}

func (r *Repository) IsAvailable(accomodationId uint, startDate time.Time, endDate time.Time, ctx context.Context) bool {
	span := tracer.StartSpanFromContext(ctx, "isAvailableRepository")
	defer span.Finish()
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

func (r *Repository) FindImagesForAccomodation(accomodationId uint) []string {
	accomodationImages := &[]model.AccomodationImage{}

	r.Db.Find(&accomodationImages, "accomodation_id = ?", accomodationId)

	var imageNames []string
	for _, accomodationImage := range *accomodationImages {
		imageNames = append(imageNames, accomodationImage.ImageName)
	}
	return imageNames
}

func (r *Repository) FindAccomodationsForHost(hostId uint, ctx context.Context) []model.Accomodation {
	span := tracer.StartSpanFromContext(ctx, "findAccomodationsForHostRepository")
	defer span.Finish()
	accomodations := &[]model.Accomodation{}

	r.Db.Find(&accomodations, "user_id = ?", hostId)
	return *accomodations
}

func (r *Repository) GetAvailableTermsForAccomodation(accomodationId uint, ctx context.Context) []model.AvailableTerm {
	span := tracer.StartSpanFromContext(ctx, "getAvailableTermsForAccomodationRepository")
	defer span.Finish()
	availableTerms := &[]model.AvailableTerm{}

	r.Db.Find(&availableTerms, "accomodation_id = ?", accomodationId)
	return *availableTerms
}

func (r *Repository) GetPricesForAccomodation(accomodationId uint, ctx context.Context) []model.Price {
	span := tracer.StartSpanFromContext(ctx, "getPricesForAccomodationRepository")
	defer span.Finish()
	prices := &[]model.Price{}

	r.Db.Find(&prices, "accomodation_id = ?", accomodationId)
	return *prices
}

