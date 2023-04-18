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
		{Name: "Vila Marija", Address: "Maksima Gorkog 17a, Novi Sad", HasWifi: true, HasKitchen: true, HasAirConditioning: true, HasFreeParking: false, MinimimGuests: 2, MaximumGuests: 5, UserId: 1},
		{Name: "Lanterna", Address: "Ljubice Ravasi 32, Novi Sad", HasWifi: true, HasKitchen: false, HasAirConditioning: false, HasFreeParking: true, MinimimGuests: 4, MaximumGuests: 4, UserId: 1},
	}
	accomodationImages = []model.AccomodationImage{
		{ImageName: "slika1.jpg", AccomodationID: 1},
		{ImageName: "slika2.jpg", AccomodationID: 1},
		{ImageName: "slika3.jpg", AccomodationID: 2},
	}

	prices = []model.Price{
		{StartDate: time.Date(2023, 1, 1, 10, 0, 0, 0, time.Local), EndDate: time.Date(2024, 1, 1, 10, 0, 0, 0, time.Local),
			Value: 3000, PriceType: model.PER_GUEST, PriceDuration: model.REGULAR, AccomodationID: 1},
		{StartDate: time.Date(2023, 1, 1, 10, 0, 0, 0, time.Local), EndDate: time.Date(2024, 1, 1, 10, 0, 0, 0, time.Local),
			Value: 5000, PriceType: model.PER_GUEST, PriceDuration: model.WEEKEND, AccomodationID: 1},
		{StartDate: time.Date(2023, 1, 1, 10, 0, 0, 0, time.Local), EndDate: time.Date(2023, 6, 1, 10, 0, 0, 0, time.Local),
			Value: 3500, PriceType: model.PER_ACCOMODATION_UNIT, PriceDuration: model.REGULAR, AccomodationID: 2},
	}

	availableTerms = []model.AvailableTerm{
		{StartDate: time.Date(2023, 1, 1, 10, 0, 0, 0, time.Local), EndDate: time.Date(2023, 3, 1, 10, 0, 0, 0, time.Local),
			AccomodationID: 1},
		{StartDate: time.Date(2023, 5, 1, 10, 0, 0, 0, time.Local), EndDate: time.Date(2024, 8, 1, 10, 0, 0, 0, time.Local),
			AccomodationID: 1},
		{StartDate: time.Date(2023, 1, 1, 10, 0, 0, 0, time.Local), EndDate: time.Date(2024, 1, 1, 10, 0, 0, 0, time.Local),
			AccomodationID: 2},
	}
)

func ConnectToDatabase() *gorm.DB {
	connectionString := "host=localhost user=postgres dbname=AccomodationServiceDB sslmode=disable password=root port=5432"
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
