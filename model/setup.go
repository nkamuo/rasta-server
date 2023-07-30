package model

import (
	"fmt"
	"log"
	"os"

	"github.com/nkamuo/rasta-server/initializers"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDatabase(config *initializers.Config) {

	// dsn := fmt.Sprintf("user:pass@tcp(127.0.0.1:3306)/dbname?charset=utf8mb4&parseTime=True&loc=Local")
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", config.DBUserName, config.DBUserPassword, config.DBHost, config.DBPort, config.DBName)

	database, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to the Database! \n", err.Error())
		os.Exit(1)
	}

	// DB.Logger = log.Default.LogMode(logger.Info)

	// log.Debug("Running Migrations")

	err = database.AutoMigrate(
		&User{},
		&Company{},
		&Product{},
		&Respondent{},
		&ProductRespondentAssignment{},

		&Order{},
		&OrderItem{},
		&Payment{},
		&PaymentMethod{},
		// REQUEST_TYPE - SPECIFIC INFORMATION
		&FuelType{},
		// LOCATIONS
		&Place{},
		&Location{},
		&LocationCoordinates{},
	)
	if err != nil {
		return
	}

	DB = database
}
