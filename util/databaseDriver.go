package util

import (
	"fmt"
	"log"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/windbnb/accomodation-service/model"
)

var (
	accomodations = []model.Accomodation{
		{Name: "Vila Marija",  Address: "Maksima Gorkog 17a, Novi Sad", HasWifi: true, HasKitchen: true, HasAirConditioning: true, HasFreeParking: false, MinimimGuests: 2, MaximumGuests: 5},
		{Name: "Lanterna",  Address: "Ljubice Ravasi 32, Novi Sad", HasWifi: true, HasKitchen: false, HasAirConditioning: false, HasFreeParking: true, MinimimGuests: 4, MaximumGuests: 4},
	}
	accomodationImages = []model.AccomodationImage{
		{ImageName: "slika1.jpg", AccomodationID: 1},
		{ImageName: "slika2.jpg", AccomodationID: 1},
		{ImageName: "slika3.jpg", AccomodationID: 2},
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
	db.AutoMigrate(&model.Accomodation{})
	db.AutoMigrate(&model.AccomodationImage{})

	for _, accomodation := range accomodations {
		db.Create(&accomodation)
	}

	for _, accomodationImage := range accomodationImages {
		db.Create(&accomodationImage)
	}

	return db
}
