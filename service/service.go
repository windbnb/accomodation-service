package service

import (
	"context"
	"errors"
	"time"

	"github.com/windbnb/accomodation-service/model"
	"github.com/windbnb/accomodation-service/repository"
	"github.com/windbnb/accomodation-service/tracer"
)

type AccomodationService struct {
	Repo repository.IRepository
}

func (s *AccomodationService) SaveAccomodation(accomodation model.Accomodation, ctx context.Context) model.Accomodation {
	span := tracer.StartSpanFromContext(ctx, "saveAccomodationService")
	defer span.Finish()

	ctx = tracer.ContextWithSpan(context.Background(), span)
	return s.Repo.SaveAccomodation(accomodation, ctx)
}

func (s *AccomodationService) UpdateAccommodationAcceptReservationType(accommodationId uint, acceptReservationType model.AcceptReservationType, hostId uint, ctx context.Context) (*model.Accomodation, error) {
	span := tracer.StartSpanFromContext(ctx, "acceptReservationTypeService")
	defer span.Finish()
	ctx = tracer.ContextWithSpan(context.Background(), span)

	if acceptReservationType != model.MANUAL && acceptReservationType != model.AUTOMATICALLY {
		return nil, errors.New("Given type does not exist")
	}

	accommodation, err := s.Repo.FindAccomodationById(accommodationId, ctx)
	if err != nil {
		tracer.LogError(span, err)
		return nil, errors.New("Given accommodation does not exist.")
	}

	if hostId != accommodation.UserId {
		return nil, errors.New("You don't have access to this entity.")
	}

	accommodation.AcceptReservationType = acceptReservationType
	s.Repo.UpdateAccommodation(accommodation, ctx)

	return &accommodation, nil
}

func (s *AccomodationService) SaveAccomodationImage(image model.AccomodationImage) model.AccomodationImage {
	return s.Repo.SaveAccomodationImage(image)
}

func (s *AccomodationService) DeleteHostAccomodation(hostId uint, ctx context.Context) error {
	span := tracer.StartSpanFromContext(ctx, "deleteHostAccomodationService")
	defer span.Finish()

	ctx = tracer.ContextWithSpan(context.Background(), span)
	return s.Repo.DeleteHostAccomodation(hostId, ctx)
}

func (s *AccomodationService) SavePrice(price model.Price) model.Price {
	return s.Repo.SavePrice(price)
}

func (s *AccomodationService) SaveAvailableTerm(availableTerm model.AvailableTerm, ctx context.Context) model.AvailableTerm {
	span := tracer.StartSpanFromContext(ctx, "saveAvailableTermService")
	defer span.Finish()

	ctx = tracer.ContextWithSpan(context.Background(), span)
	return s.Repo.SaveAvailableTerm(availableTerm, ctx)
}

func (s *AccomodationService) SaveReservedTerm(reservedTerm model.ReservedTerm, ctx context.Context) model.ReservedTerm {
	span := tracer.StartSpanFromContext(ctx, "saveReservedTermService")
	defer span.Finish()

	ctx = tracer.ContextWithSpan(context.Background(), span)
	return s.Repo.SaveReservedTerm(reservedTerm, ctx)
}

func (s *AccomodationService) UpdatePrice(price model.Price, id uint64, ctx context.Context) (model.Price, error) {
	span := tracer.StartSpanFromContext(ctx, "updatePriceService")
	defer span.Finish()
	ctx = tracer.ContextWithSpan(context.Background(), span)
	var priceToUpdate model.Price

	priceToUpdate, _ = s.FindPriceById(id, ctx)
	if priceToUpdate.ID != 0 {
		priceToUpdate.StartDate = price.StartDate
		priceToUpdate.EndDate = price.EndDate
		priceToUpdate.Value = price.Value
		s.Repo.UpdatePrice(priceToUpdate, ctx)
	}

	return priceToUpdate, nil
}

func (s *AccomodationService) UpdateAvailableTerm(availableTerm model.AvailableTerm, id uint64, ctx context.Context) (model.AvailableTerm, error) {
	span := tracer.StartSpanFromContext(ctx, "updateAvailableTermService")
	defer span.Finish()
	ctx = tracer.ContextWithSpan(context.Background(), span)
	var availableTermToUpdate model.AvailableTerm

	availableTermToUpdate, _ = s.FindAvailableTermById(id, ctx)
	if availableTermToUpdate.ID != 0 {
		availableTermToUpdate.StartDate = availableTerm.StartDate
		availableTermToUpdate.EndDate = availableTerm.EndDate
		s.Repo.UpdateAvailableTerm(availableTermToUpdate, ctx)
	}

	return availableTermToUpdate, nil
}

func (s *AccomodationService) DeletePrice(id uint64, ctx context.Context) error {
	span := tracer.StartSpanFromContext(ctx, "deletePriceService")
	defer span.Finish()
	ctx = tracer.ContextWithSpan(context.Background(), span)
	var priceToDelete model.Price

	priceToDelete, _ = s.FindPriceById(id, ctx)
	if priceToDelete.ID == 0 {
		err := errors.New("price with given id does not exist")
		tracer.LogError(span, err)
		return err
	}

	return s.Repo.DeletePrice(id, ctx)
}

func (s *AccomodationService) DeleteAvailableTerm(id uint64, ctx context.Context) error {
	span := tracer.StartSpanFromContext(ctx, "deleteAvailableTermService")
	defer span.Finish()
	ctx = tracer.ContextWithSpan(context.Background(), span)
	var availableTermToDelete model.AvailableTerm

	availableTermToDelete, _ = s.FindAvailableTermById(id, ctx)
	if availableTermToDelete.ID == 0 {
		err := errors.New("available term with given id does not exist")
		tracer.LogError(span, err)
		return err
	}

	return s.Repo.DeleteAvailableTerm(id)
}

func (s *AccomodationService) DeleteReservedTerm(id uint64, ctx context.Context) error {
	span := tracer.StartSpanFromContext(ctx, "deleteReservedTermService")
	defer span.Finish()
	ctx = tracer.ContextWithSpan(context.Background(), span)
	var reservedTermToDelete model.ReservedTerm

	reservedTermToDelete, _ = s.FindReservedTermById(id, ctx)
	if reservedTermToDelete.ID == 0 {
		err := errors.New("reserved term with given id does not exist")
		tracer.LogError(span, err)
		return err
	}

	return s.Repo.DeleteReservedTerm(id)
}

func (service *AccomodationService) FindAccomodationById(id uint, ctx context.Context) (model.Accomodation, error) {
	span := tracer.StartSpanFromContext(ctx, "findAccomodationByIdService")
	defer span.Finish()
	ctx = tracer.ContextWithSpan(context.Background(), span)

	accomodation, err := service.Repo.FindAccomodationById(id, ctx)

	if err != nil {
		tracer.LogError(span, err)
		return model.Accomodation{}, errors.New("accomodation with given id does not exist")
	}

	return accomodation, nil
}

func (s *AccomodationService) FindPriceById(id uint64, ctx context.Context) (model.Price, error) {
	span := tracer.StartSpanFromContext(ctx, "findPriceByIdService")
	defer span.Finish()
	ctx = tracer.ContextWithSpan(context.Background(), span)
	price, err := s.Repo.FindPriceById(id, ctx)

	if err != nil {
		tracer.LogError(span, err)
		return model.Price{}, errors.New("price with given id does not exist")
	}

	return price, nil
}

func (service *AccomodationService) FindAvailableTermById(id uint64, ctx context.Context) (model.AvailableTerm, error) {
	span := tracer.StartSpanFromContext(ctx, "findAvailableTermByIdService")
	defer span.Finish()
	ctx = tracer.ContextWithSpan(context.Background(), span)
	availableTerm, err := service.Repo.FindAvailableTermById(id, ctx)

	if err != nil {
		tracer.LogError(span, err)
		return model.AvailableTerm{}, errors.New("available term with given id does not exist")
	}

	return availableTerm, nil
}

func (service *AccomodationService) FindReservedTermById(id uint64, ctx context.Context) (model.ReservedTerm, error) {
	span := tracer.StartSpanFromContext(ctx, "findReservedTermByIdService")
	defer span.Finish()
	ctx = tracer.ContextWithSpan(context.Background(), span)
	reservedTerm, err := service.Repo.FindReservedTermById(id, ctx)

	if err != nil {
		tracer.LogError(span, err)
		return model.ReservedTerm{}, errors.New("reserved term with given id does not exist")
	}

	return reservedTerm, nil
}

func (service *AccomodationService) FindAvailableTerms(accommodationId uint, ctx context.Context) ([]model.AvailableTerm, error) {
	span := tracer.StartSpanFromContext(ctx, "findAvailableTermsService")
	defer span.Finish()
	ctx = tracer.ContextWithSpan(context.Background(), span)

	_, err := service.Repo.FindAccomodationById(accommodationId, ctx)

	if err != nil {
		tracer.LogError(span, err)
		return []model.AvailableTerm{}, errors.New("accommodation with given id does not exist")
	}

	var availableTerms = service.Repo.FindAvailableTermAfter(accommodationId, time.Now(), ctx)

	return availableTerms, nil
}

func (service *AccomodationService) CalculatePrice(accommodation model.Accomodation, searchAccomodationDTO model.SearchAccomodationDTO) (float32, int) {
	prices := service.Repo.FindPricesForAccomodation(accommodation.ID, searchAccomodationDTO.StartDate, searchAccomodationDTO.EndDate)
	var basePrice float32 = 0
	for _, price := range prices {
		if price.PriceDuration == model.HOLIDAY || price.PriceDuration == model.WEEKEND {
			basePrice = price.Value
			break
		} else {
			basePrice = price.Value
		}
	}

	var totalPrice float32 = 0
	if accommodation.PriceType == model.PER_GUEST {
		totalPrice = float32(searchAccomodationDTO.NumberOfGuests) * basePrice * float32(searchAccomodationDTO.EndDate.Sub(searchAccomodationDTO.StartDate).Hours()/24)
	} else {
		totalPrice = basePrice
	}

	return basePrice, int(totalPrice)
}

func (service *AccomodationService) SearchAccomodations(searchAccomodationDTO model.SearchAccomodationDTO, ctx context.Context) []model.SearchAccomodationReturnDTO {
	span := tracer.StartSpanFromContext(ctx, "searchAccomodationsService")
	defer span.Finish()
	ctx = tracer.ContextWithSpan(context.Background(), span)
	accomodations := service.Repo.FindAccomodationByGuestsAndAddress(searchAccomodationDTO.NumberOfGuests, searchAccomodationDTO.Address, ctx)

	var availableAccomodations []model.SearchAccomodationReturnDTO
	for _, accommodation := range accomodations {
		if service.Repo.IsAvailable(accommodation.ID, searchAccomodationDTO.StartDate, searchAccomodationDTO.EndDate, ctx) == true {
			if service.Repo.IsReserved(accommodation.ID, searchAccomodationDTO.StartDate, searchAccomodationDTO.EndDate, ctx) == false {
				basePrice, totalPrice := service.CalculatePrice(accommodation, searchAccomodationDTO)
				accomodationDTO := accommodation.ToDTO()
				accomodationDTO.Images = service.Repo.FindImagesForAccomodation(accommodation.ID)
				var searchAccomodationReturnDTO model.SearchAccomodationReturnDTO
				searchAccomodationReturnDTO.Accomodation = accomodationDTO
				searchAccomodationReturnDTO.Price = basePrice
				searchAccomodationReturnDTO.TotalPrice = totalPrice
				searchAccomodationReturnDTO.StartDate = searchAccomodationDTO.StartDate
				searchAccomodationReturnDTO.EndDate = searchAccomodationDTO.EndDate
				searchAccomodationReturnDTO.NumberOfGuests = searchAccomodationDTO.NumberOfGuests
				availableAccomodations = append(availableAccomodations, searchAccomodationReturnDTO)
			}
		}
	}
	return availableAccomodations
}

func (service *AccomodationService) FindAccommodationsForHost(hostId uint, ctx context.Context) []model.AccomodationDTO {
	span := tracer.StartSpanFromContext(ctx, "findAccomodationsForHostService")
	defer span.Finish()
	ctx = tracer.ContextWithSpan(context.Background(), span)
	accomodations := service.Repo.FindAccomodationsForHost(hostId, ctx)

	var hostAccomodations []model.AccomodationDTO
	for _, accommodation := range accomodations {
		accomodationDTO := accommodation.ToDTO()
		accomodationDTO.Images = service.Repo.FindImagesForAccomodation(accommodation.ID)
		hostAccomodations = append(hostAccomodations, accomodationDTO)

	}
	return hostAccomodations
}

func (service *AccomodationService) GetAvailableTermsForAccomodation(accomodationId uint, ctx context.Context) []model.AvailableTermDTO {
	span := tracer.StartSpanFromContext(ctx, "getAvailableTermsForAccomodationService")
	defer span.Finish()
	ctx = tracer.ContextWithSpan(context.Background(), span)
	availableTerms := service.Repo.GetAvailableTermsForAccomodation(accomodationId, ctx)

	var availableTermsDTO []model.AvailableTermDTO
	for _, availableTerm := range availableTerms {
		availableTermDTO := availableTerm.ToDTO()
		availableTermsDTO = append(availableTermsDTO, availableTermDTO)

	}
	return availableTermsDTO
}

func (service *AccomodationService) GetPricesForAccomodation(accomodationId uint, ctx context.Context) []model.PriceDTO {
	span := tracer.StartSpanFromContext(ctx, "getPricesForAccomodationService")
	defer span.Finish()
	ctx = tracer.ContextWithSpan(context.Background(), span)
	prices := service.Repo.GetPricesForAccomodation(accomodationId, ctx)

	var pricesDTO []model.PriceDTO
	for _, price := range prices {
		priceDTO := price.ToDTO()
		pricesDTO = append(pricesDTO, priceDTO)

	}
	return pricesDTO
}
