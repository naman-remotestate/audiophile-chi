package database

import (
	"audiophile/models"
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var (
	DB *gorm.DB
)

type SSLMode string

const (
	SSLModeEnable  SSLMode = "enable"
	SSLModeDisable SSLMode = "disable"
)

func ConnectToDatabase(host, port, databaseName, user, password string, sslMode SSLMode) error {
	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s", host, port, user, password, databaseName, sslMode)
	db, err := gorm.Open(postgres.Open(connStr), &gorm.Config{})
	if err != nil {
		return err
	}
	DB = db
	err = DB.AutoMigrate(&models.Users{}, &models.Roles{}, &models.Address{},
		&models.ProductVariants{}, &models.ProductImages{},
		&models.Carts{}, &models.CartDetails{}, &models.Inventory{}, &models.Orders{}, &models.Session{}, models.Images{})
	return err
}
