package service

import (
	"github.com/windbnb/accomodation-service/model"
	"github.com/windbnb/accomodation-service/repository"
)

type AccomodationService struct {
	Repo *repository.Repository
}

func (s *AccomodationService) SaveAccomodation(accomodation model.Accomodation) (model.Accomodation) {
	return s.Repo.SaveAccomodation(accomodation)
}

func (s *AccomodationService) SaveAccomodationImage(image model.AccomodationImage) (model.AccomodationImage) {
	return s.Repo.SaveAccomodationImage(image)
}

func (s *AccomodationService) DeleteHostAccomodation(hostId uint) error {
	return s.Repo.DeleteHostAccomodation(hostId)
}
