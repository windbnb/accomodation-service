package util

import (
	"fmt"
	"log"
	"time"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/windbnb/accomodation-service/model"
)

var (
	accomodations = []model.Accomodation{
		{Name: "Vila Marija", Address: "Maksima Gorkog 17a, Novi Sad", HasWifi: true, HasKitchen: true, HasAirConditioning: true, HasFreeParking: false, MinimimGuests: 2, MaximumGuests: 5, UserId: 1, AcceptReservationType: model.MANUAL, PriceType: model.PER_GUEST},
		{Name: "Lanterna", Address: "Ljubice Ravasi 32, Novi Sad", HasWifi: true, HasKitchen: false, HasAirConditioning: false, HasFreeParking: true, MinimimGuests: 4, MaximumGuests: 4, UserId: 1, AcceptReservationType: model.AUTOMATICALLY, PriceType: model.PER_GUEST},
	}
	accomodationImages = []model.AccomodationImage{
		{ImageName: "373488187.jpg", AccomodationID: 1},
		{ImageName: "373487944.jpg", AccomodationID: 1},
		{ImageName: "373486431.jpg", AccomodationID: 1},
		{ImageName: "242225269.jpg", AccomodationID: 2},
		{ImageName: "242218937.jpg", AccomodationID: 2},
		{ImageName: "242216685.jpg", AccomodationID: 2},
	}

	prices = []model.Price{
		{StartDate: time.Date(2023, 1, 1, 10, 0, 0, 0, time.Local), EndDate: time.Date(2024, 1, 1, 10, 0, 0, 0, time.Local),
			Value: 3000, PriceDuration: model.REGULAR, AccomodationID: 1, Active: true},
		{StartDate: time.Date(2023, 1, 1, 10, 0, 0, 0, time.Local), EndDate: time.Date(2024, 1, 1, 10, 0, 0, 0, time.Local),
			Value: 5000, PriceDuration: model.HOLIDAY, AccomodationID: 1, Active: true},
		{StartDate: time.Date(2023, 1, 1, 10, 0, 0, 0, time.Local), EndDate: time.Date(2024, 1, 1, 10, 0, 0, 0, time.Local),
			Value: 3500, PriceDuration: model.REGULAR, AccomodationID: 2, Active: true},
	}

	availableTerms = []model.AvailableTerm{
		{StartDate: time.Date(2023, 1, 1, 10, 0, 0, 0, time.Local), EndDate: time.Date(2023, 5, 1, 10, 0, 0, 0, time.Local),
			AccomodationID: 1},
		{StartDate: time.Date(2023, 5, 1, 10, 0, 0, 0, time.Local), EndDate: time.Date(2024, 1, 1, 10, 0, 0, 0, time.Local),
			AccomodationID: 1},
		{StartDate: time.Date(2023, 1, 1, 10, 0, 0, 0, time.Local), EndDate: time.Date(2024, 1, 1, 10, 0, 0, 0, time.Local),
			AccomodationID: 2},
	}
)

func ConnectToDatabase() *gorm.DB {
	connectionString := "host=accomodation_db user=postgres dbname=AccomodationServiceDB sslmode=disable password=root port=5432"
	dialect := "postgres"

	db, err := gorm.Open(dialect, connectionString)
	if err != nil {
		log.Fatal(err)
	} else {
		fmt.Println("Connection to DB successfull.")
	}

	db.DropTable("accomodations")
	db.DropTable("accomodation_images")
	db.DropTable("prices")
	db.DropTable("reserved_terms")
	db.DropTable("available_terms")
	db.AutoMigrate(&model.Accomodation{})
	db.AutoMigrate(&model.AccomodationImage{})
	db.AutoMigrate(&model.Price{})
	db.AutoMigrate(&model.ReservedTerm{})
	db.AutoMigrate(&model.AvailableTerm{})

	for _, accomodation := range accomodations {
		db.Create(&accomodation)
	}

	for _, accomodationImage := range accomodationImages {
		db.Create(&accomodationImage)
	}

	for _, price := range prices {
		db.Create(&price)
	}

	for _, availableTerm := range availableTerms {
		db.Create(&availableTerm)
	}

	return db
}
