package service_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/windbnb/accomodation-service/model"
	"github.com/windbnb/accomodation-service/repository"
	"github.com/windbnb/accomodation-service/service"
	"github.com/windbnb/accomodation-service/util"
)

func TestUpdateAccommodationAcceptReservationType_AccomodationDoesNotExist(t *testing.T) {
	mockRepo := &MockRepo{
		FindAccomodationByIdFn: func(id uint, ctx context.Context) (model.Accomodation, error) {
			return model.Accomodation{}, errors.New("Given accommodation does not exist.")
		},
	}

	accommodationService := service.AccomodationService{
		Repo: mockRepo,
	}

	accomodation, err := accommodationService.UpdateAccommodationAcceptReservationType(1, model.MANUAL, 1, context.Background())

	assert.Empty(t, accomodation)
	assert.EqualError(t, err, "Given accommodation does not exist.")

}

func TestUpdateAccommodationAcceptReservationType_UserDoesNotHaveAccess(t *testing.T) {
	mockRepo := &MockRepo{
		FindAccomodationByIdFn: func(id uint, ctx context.Context) (model.Accomodation, error) {
			return model.Accomodation{
				UserId: 1,
			}, nil
		},
	}

	accommodationService := service.AccomodationService{
		Repo: mockRepo,
	}

	accomodation, err := accommodationService.UpdateAccommodationAcceptReservationType(1, model.MANUAL, 2, context.Background())

	assert.Empty(t, accomodation)
	assert.EqualError(t, err, "You don't have access to this entity.")

}

func TestUpdateAccommodationAcceptReservationType_Successfull(t *testing.T) {
	mockRepo := &MockRepo{
		FindAccomodationByIdFn: func(id uint, ctx context.Context) (model.Accomodation, error) {
			return model.Accomodation{
				Name:               "Vila",
				Address:            "Bulevar oslobodjenja 1, Novi Sad",
				HasWifi:            true,
				HasKitchen:         true,
				HasAirConditioning: true,
				HasFreeParking:     true,
				MinimimGuests:      1,
				MaximumGuests:      5,
				Images: []model.AccomodationImage{
					{ImageName: "slika1.jpg", AccomodationID: 1},
				},
				UserId: 2,
				Prices: []model.Price{
					{StartDate: time.Date(2023, 1, 1, 10, 0, 0, 0, time.Local), EndDate: time.Date(2024, 1, 1, 10, 0, 0, 0, time.Local),
						Value: 3000, PriceDuration: model.REGULAR, AccomodationID: 1, Active: true}},
				PriceType:             model.PER_GUEST,
				AcceptReservationType: model.MANUAL,
			}, nil
		},
		UpdateAccommodationFn: func(accomodation model.Accomodation, ctx context.Context) model.Accomodation {
			return accomodation
		},
	}

	accommodationService := service.AccomodationService{
		Repo: mockRepo,
	}

	accomodation := model.Accomodation{
		Name:               "Vila",
		Address:            "Bulevar oslobodjenja 1, Novi Sad",
		HasWifi:            true,
		HasKitchen:         true,
		HasAirConditioning: true,
		HasFreeParking:     true,
		MinimimGuests:      1,
		MaximumGuests:      5,
		Images: []model.AccomodationImage{
			{ImageName: "slika1.jpg", AccomodationID: 1},
		},
		UserId: 2,
		Prices: []model.Price{
			{StartDate: time.Date(2023, 1, 1, 10, 0, 0, 0, time.Local), EndDate: time.Date(2024, 1, 1, 10, 0, 0, 0, time.Local),
				Value: 3000, PriceDuration: model.REGULAR, AccomodationID: 1, Active: true}},
		PriceType:             model.PER_GUEST,
		AcceptReservationType: model.AUTOMATICALLY,
	}

	updatedAccomodation, err := accommodationService.UpdateAccommodationAcceptReservationType(1, model.AUTOMATICALLY, 2, context.Background())
	assert.Equal(t, &accomodation, updatedAccomodation)
	assert.NoError(t, err)

}

func TestUpdateAccommodationAcceptReservationType_AccomodationDoesNotExist_Integration(t *testing.T) {
	db := util.ConnectToDatabase()
	defer db.Close()
	accomodationService := service.AccomodationService{Repo: &repository.Repository{Db: db}}

	updatedAccomodation, err := accomodationService.UpdateAccommodationAcceptReservationType(10, model.AUTOMATICALLY, 2, context.Background())

	assert.Empty(t, updatedAccomodation)
	assert.EqualError(t, err, "Given accommodation does not exist.")
}

func TestUpdateAccommodationAcceptReservationType_UserDoesNotHaveAccess_Integration(t *testing.T) {
	db := util.ConnectToDatabase()
	defer db.Close()
	accomodationService := service.AccomodationService{Repo: &repository.Repository{Db: db}}

	updatedAccomodation, err := accomodationService.UpdateAccommodationAcceptReservationType(1, model.AUTOMATICALLY, 2, context.Background())

	assert.Empty(t, updatedAccomodation)
	assert.EqualError(t, err, "You don't have access to this entity.")

}

func TestUpdateAccommodationAcceptReservationType_Successfull_Integration(t *testing.T) {
	db := util.ConnectToDatabase()
	defer db.Close()
	accomodationService := service.AccomodationService{Repo: &repository.Repository{Db: db}}

	updatedAccomodation, err := accomodationService.UpdateAccommodationAcceptReservationType(1, model.AUTOMATICALLY, 1, context.Background())
	accomodation, _ := accomodationService.FindAccomodationById(1, context.Background())
	assert.Equal(t, accomodation.AcceptReservationType, updatedAccomodation.AcceptReservationType)
	assert.NoError(t, err)

}

type MockRepo struct {
	repository.Repository
	UpdateAccommodationFn  func(accomodation model.Accomodation, ctx context.Context) model.Accomodation
	FindAccomodationByIdFn func(id uint, ctx context.Context) (model.Accomodation, error)
}

func (m *MockRepo) UpdateAccommodation(accomodation model.Accomodation, ctx context.Context) model.Accomodation {
	return m.UpdateAccommodationFn(accomodation, ctx)
}

func (m *MockRepo) FindAccomodationById(id uint, ctx context.Context) (model.Accomodation, error) {
	return m.FindAccomodationByIdFn(id, ctx)
}
