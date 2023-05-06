package service

import (
	"errors"

	"github.com/windbnb/accomodation-service/model"
	"github.com/windbnb/accomodation-service/repository"
)

type AccomodationService struct {
	Repo *repository.Repository
}

func (s *AccomodationService) SaveAccomodation(accomodation model.Accomodation) model.Accomodation {
	return s.Repo.SaveAccomodation(accomodation)
}

func (s *AccomodationService) SaveAccomodationImage(image model.AccomodationImage) model.AccomodationImage {
	return s.Repo.SaveAccomodationImage(image)
}

func (s *AccomodationService) DeleteHostAccomodation(hostId uint) error {
	return s.Repo.DeleteHostAccomodation(hostId)
}

func (s *AccomodationService) SavePrice(price model.Price) model.Price {
	return s.Repo.SavePrice(price)
}

func (s *AccomodationService) SaveAvailableTerm(availableTerm model.AvailableTerm) model.AvailableTerm {
	return s.Repo.SaveAvailableTerm(availableTerm)
}

func (s *AccomodationService) SaveReservedTerm(reservedTerm model.ReservedTerm) model.ReservedTerm {
	return s.Repo.SaveReservedTerm(reservedTerm)
}

func (s *AccomodationService) UpdatePrice(price model.Price, id uint64) (model.Price, error) {
	var priceToUpdate model.Price

	priceToUpdate, _ = s.FindPriceById(id)
	if priceToUpdate.ID != 0 {
		priceToUpdate.StartDate = price.StartDate
		priceToUpdate.EndDate = price.EndDate
		priceToUpdate.Value = price.Value
		s.Repo.UpdatePrice(priceToUpdate)
	}

	return priceToUpdate, nil
}

func (s *AccomodationService) UpdateAvailableTerm(availableTerm model.AvailableTerm, id uint64) (model.AvailableTerm, error) {
	var availableTermToUpdate model.AvailableTerm

	availableTermToUpdate, _ = s.FindAvailableTermById(id)
	if availableTermToUpdate.ID != 0 {
		availableTermToUpdate.StartDate = availableTerm.StartDate
		availableTermToUpdate.EndDate = availableTerm.EndDate
		s.Repo.UpdateAvailableTerm(availableTermToUpdate)
	}

	return availableTermToUpdate, nil
}

func (s *AccomodationService) DeletePrice(id uint64) error {
	var priceToDelete model.Price

	priceToDelete, _ = s.FindPriceById(id)
	if priceToDelete.ID == 0 {
		return errors.New("price with given id does not exist")
	}

	return s.Repo.DeletePrice(id)
}

func (s *AccomodationService) DeleteAvailableTerm(id uint64) error {
	var availableTermToDelete model.AvailableTerm

	availableTermToDelete, _ = s.FindAvailableTermById(id)
	if availableTermToDelete.ID == 0 {
		return errors.New("available term with given id does not exist")
	}

	return s.Repo.DeleteAvailableTerm(id)
}

func (s *AccomodationService) DeleteReservedTerm(id uint64) error {
	var reservedTermToDelete model.ReservedTerm

	reservedTermToDelete, _ = s.FindReservedTermById(id)
	if reservedTermToDelete.ID == 0 {
		return errors.New("reserved term with given id does not exist")
	}

	return s.Repo.DeleteReservedTerm(id)
}

func (service *AccomodationService) FindAccomodationById(id uint) (model.Accomodation, error) {
	accomodation, err := service.Repo.FindAccomodationById(id)

	if err != nil {
		return model.Accomodation{}, errors.New("accomodation with given id does not exist")
	}

	return accomodation, nil
}

func (s *AccomodationService) FindPriceById(id uint64) (model.Price, error) {
	price, err := s.Repo.FindPriceById(id)

	if err != nil {
		return model.Price{}, errors.New("price with given id does not exist")
	}

	return price, nil
}

func (service *AccomodationService) FindAvailableTermById(id uint64) (model.AvailableTerm, error) {
	availableTerm, err := service.Repo.FindAvailableTermById(id)

	if err != nil {
		return model.AvailableTerm{}, errors.New("available term with given id does not exist")
	}

	return availableTerm, nil
}

func (service *AccomodationService) FindReservedTermById(id uint64) (model.ReservedTerm, error) {
	reservedTerm, err := service.Repo.FindReservedTermById(id)

	if err != nil {
		return model.ReservedTerm{}, errors.New("reserved term with given id does not exist")
	}

	return reservedTerm, nil
}
