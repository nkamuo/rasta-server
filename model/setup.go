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

func ConnectDatabase(config *initializers.Config) (err error) {

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
		&UserVerificationRequest{},
		&Company{},
		&Product{},
		&Respondent{},
		&ProductRespondentAssignment{},
		&RespondentSession{},
		&RespondentSessionAssignedProduct{},
		//
		&CompanyWallet{},
		&RespondentWallet{},
		&RespondentEarning{},
		&CompanyEarning{},
		&CompanyWithdrawal{},
		&RespondentWithdrawal{},
		//
		&RespondentAccessProductBalance{},
		&RespondentAccessProductSubscription{},
		&RespondentAccessProductPurchase{},
		&RespondentAccessProductPrice{},
		//
		&MotoristRequestSituation{},
		&OrderMotoristRequestSituation{},
		//
		&RespondentOrderCharge{},
		&OrderFulfilment{},
		&Order{},
		&Request{},
		&OrderAdjustment{},
		&OrderPayment{},
		&RequestVehicleInfo{},
		&RequestFuelTypeInfo{},
		&Payment{},
		&PaymentMethod{},
		&TowingPlaceRate{},
		// REQUEST_TYPE - SPECIFIC INFORMATION
		&FuelType{},
		&FuelTypePlaceRate{},
		// LOCATIONS
		&Place{},
		&Location{},
		&RespondentServiceReview{},
		&UserPassword{},
		// &LocationCoordinates{},
		&RespondentSessionLocationEntry{},
		//
		&ImageDocument{},
	)
	if err != nil {
		return err
	}

	DB = database
	MigrateMotoristSituations(database)
	return err
}
