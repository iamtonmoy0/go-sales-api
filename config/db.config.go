package database

import (
	"log"

	"github.com/iamtonmoy0/go-sales-api/models"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func Database() *gorm.DB {
	db, err := gorm.Open(sqlite.Open("gorm.db"), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}

	db.AutoMigrate(
		&models.Cashier{},
		&models.Category{},
		&models.Discount{},
		&models.Order{},
		&models.ProductResponseOrder{},
		&models.ProductOrder{},
		&models.RevenueResponse{},
		&models.SoldResponse{},
		&models.Payment{},
		&models.PaymentType{},
		&models.Product{},
		&models.ProductResult{},
	)

	return db
}
