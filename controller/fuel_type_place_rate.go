package controller

import (
	"fmt"
	"net/http"

	"github.com/google/uuid"
	"github.com/nkamuo/rasta-server/data/pagination"
	"github.com/nkamuo/rasta-server/dto"
	"github.com/nkamuo/rasta-server/model"
	"github.com/nkamuo/rasta-server/service"

	"github.com/gin-gonic/gin"
)

func FindFuelTypePlaceRates(c *gin.Context) {
	var fuelTypePlaceRates []model.FuelTypePlaceRate
	var page pagination.Page
	if err := c.ShouldBindQuery(&page); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": err.Error()})
		return
	}

	if err := model.DB.Scopes(pagination.Paginate(fuelTypePlaceRates, &page, model.DB)).Find(&fuelTypePlaceRates).Error; nil != err {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": err.Error()})
		return
	}
	page.Rows = fuelTypePlaceRates
	c.JSON(http.StatusOK, gin.H{"data": page})
}

func CreateFuelTypePlaceRate(c *gin.Context) {

	placeService := service.GetPlaceService()
	fuelTypeService := service.GetFuelTypeService()
	fuelTypePlaceRateService := service.GetFuelTypePlaceRateService()

	var input dto.FuelTypePlaceRateCreationInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	_, err := fuelTypeService.GetById(input.FuelTypeID)
	if err != nil {
		message := fmt.Sprintf("Could not find fuel type with [id:%s]: %s", input.FuelTypeID, err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": message})
		return
	}

	_, err = placeService.GetById(input.PlaceID)
	if err != nil {
		message := fmt.Sprintf("Could not find place with [id:%s]: %s", input.PlaceID, err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": message})
		return
	}

	var fuelTypePlaceRate *model.FuelTypePlaceRate

	if err := model.DB.
		Where("fuel_type_id = ? AND place_id = ?", input.FuelTypeID, input.PlaceID).
		First(fuelTypePlaceRate).Error; nil != err {
		if err.Error() != "record not found" {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
	}

	if nil == fuelTypePlaceRate {
		fuelTypePlaceRate = &model.FuelTypePlaceRate{}
	}

	fuelTypePlaceRate.Description = input.Description
	fuelTypePlaceRate.FuelTypeID = &input.FuelTypeID
	fuelTypePlaceRate.PlaceID = &input.PlaceID
	fuelTypePlaceRate.Active = &input.Active

	// Description: input.Description,
	// Rate:        input.Rate,
	// Active: &input.Active,

	if err := fuelTypePlaceRateService.Save(fuelTypePlaceRate); nil != err {
		c.JSON(http.StatusOK, gin.H{"status": "error", "message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": fuelTypePlaceRate, "status": "success"})
}

func FindFuelTypePlaceRate(c *gin.Context) {
	fuelTypePlaceRateService := service.GetFuelTypePlaceRateService()

	id, err := uuid.Parse(c.Param("id"))
	if nil != err {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "Invalid Id provided"})
		return
	}

	fuelTypePlaceRate, err := fuelTypePlaceRateService.GetById(id)
	if err != nil {
		message := fmt.Sprintf("Could not find fuelTypePlaceRate with [id:%s]", id)
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": message})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": fuelTypePlaceRate})
}

func UpdateFuelTypePlaceRate(c *gin.Context) {
	fuelTypePlaceRateService := service.GetFuelTypePlaceRateService()

	var input dto.FuelTypePlaceRateUpdateInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	id, err := uuid.Parse(c.Param("id"))
	if nil != err {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "Invalid Id provided"})
		return
	}
	fuelTypePlaceRate, err := fuelTypePlaceRateService.GetById(id)
	if nil != err {
		message := fmt.Sprintf("Could not find fuelTypePlaceRate with [id:%s]", id)
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": message})
	}

	if nil != input.Rate {
		fuelTypePlaceRate.Rate = *input.Rate
	}
	if nil != input.Active {
		fuelTypePlaceRate.Active = *&input.Active
	}
	if nil != input.Description {
		fuelTypePlaceRate.Description = input.Description
	}

	if err := fuelTypePlaceRateService.Save(fuelTypePlaceRate); nil != err {
		c.JSON(http.StatusOK, gin.H{"status": "error", "message": err.Error()})
		return
	}

	message := fmt.Sprintf("Updated model \"%s\"", fuelTypePlaceRate.ID)
	c.JSON(http.StatusOK, gin.H{"data": fuelTypePlaceRate, "status": "success", "message": message})
}

func DeleteFuelTypePlaceRate(c *gin.Context) {
	id := c.Param("id")

	var fuelTypePlaceRate model.FuelTypePlaceRate

	if err := model.DB.Where("id = ?", id).First(&fuelTypePlaceRate).Error; err != nil {
		message := fmt.Sprintf("Could not find fuelTypePlaceRate with [id:%s]", id)
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": message})
		return
	}
	model.DB.Delete(&fuelTypePlaceRate)
	message := fmt.Sprintf("Deleted fuelTypePlaceRate \"%s\"", fuelTypePlaceRate.ID)
	c.JSON(http.StatusOK, gin.H{"data": fuelTypePlaceRate, "status": "success", "message": message})
}
