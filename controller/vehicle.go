package controller

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/google/uuid"
	"github.com/nkamuo/rasta-server/data/pagination"
	"github.com/nkamuo/rasta-server/dto"
	"github.com/nkamuo/rasta-server/model"
	"github.com/nkamuo/rasta-server/service"

	"github.com/gin-gonic/gin"
)

// GET /vehicles
// Get all vehicles
func FindVehicles(c *gin.Context) {
	var vehicles []model.Vehicle
	var page pagination.Page
	if err := c.ShouldBindQuery(&page); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": err.Error()})
		return
	}

	query := model.DB.Model(&model.Vehicle{}).Preload("Model").Preload("Owner")
	query = query.Scopes(pagination.Paginate(vehicles, &page, query))

	if err := query.Find(&vehicles).Error; nil != err {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": err.Error()})
		return
	}
	page.Rows = vehicles
	c.JSON(http.StatusOK, gin.H{"data": page})
}

func CreateVehicle(c *gin.Context) {

	userService := service.GetUserService()
	vehicleService := service.GetVehicleService()
	modelService := service.GetVehicleModelService()

	var input dto.VehicleCreationInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	vehicleModel, err := modelService.GetById(input.ModelID)
	if nil != err {
		message := fmt.Sprintf("Could not resolve the specified Model with [id:%s]: %s", input.ModelID, err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"message": message, "status": "error"})
		return
	}

	owner, err := userService.GetById(input.OwnerID)
	if nil != err {
		message := fmt.Sprintf("Could not resolve the specified User with [id:%s]: %s", input.OwnerID, err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"message": message, "status": "error"})
		return
	}

	vehicle := model.Vehicle{
		LicensePlateNumber: *input.LicensePlateNumber,
		ModelID:            &vehicleModel.ID,
		OwnerID:            &owner.ID,
		Color:              *input.Color,
		Description:        *input.Description,
		Published:          &input.Published,
	}
	if err := vehicleService.Save(&vehicle); nil != err {
		c.JSON(http.StatusOK, gin.H{"status": "error", "message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": vehicle, "status": "success"})
}

func FindVehicle(c *gin.Context) {
	// vehicleService := service.GetVehicleService()

	id, err := uuid.Parse(c.Param("id"))
	if nil != err {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "Invalid Id provided"})
		return
	}

	var vehicle model.Vehicle

	query := model.DB.Preload("Model").Preload("Owner")
	if err := query.First(&vehicle, "id = ?", id).Error; err != nil {
		message := fmt.Sprintf("Could not find vehicle with [id:%s]", id)
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": message})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": vehicle})
}

func UpdateVehicle(c *gin.Context) {
	userService := service.GetUserService()
	vehicleService := service.GetVehicleService()
	modelService := service.GetVehicleModelService()

	var input dto.VehicleUpdateInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	id, err := uuid.Parse(c.Param("id"))
	if nil != err {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "Invalid Id provided"})
		return
	}
	vehicle, err := vehicleService.GetById(id)
	if nil != err {
		message := fmt.Sprintf("Could not find vehicle with [id:%s]", id)
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": message})
	}

	if nil != input.Color {
		vehicle.Color = *input.Color
	}

	if nil != input.LicensePlateNumber {
		vehicle.LicensePlateNumber = *input.LicensePlateNumber
	}
	if nil != input.Description {
		vehicle.Description = *input.Description
	}

	if nil != input.Published {
		vehicle.Published = input.Published
	}

	if nil != input.ModelID {
		vehicleModel, err := modelService.GetById(*input.ModelID)
		if nil != err {
			message := fmt.Sprintf("Could not resolve the specified Model with [id:%s]: %s", input.ModelID, err.Error())
			c.JSON(http.StatusBadRequest, gin.H{"message": message, "status": "error"})
			return
		}
		vehicle.ModelID = &vehicleModel.ID
	}

	if nil != input.OwnerID {
		owner, err := userService.GetById(*input.OwnerID)
		if nil != err {
			message := fmt.Sprintf("Could not resolve the specified User with [id:%s]: %s", input.OwnerID, err.Error())
			c.JSON(http.StatusBadRequest, gin.H{"message": message, "status": "error"})
			return
		}
		vehicle.OwnerID = &owner.ID
	}

	if err := vehicleService.Save(vehicle); nil != err {
		c.JSON(http.StatusOK, gin.H{"status": "error", "message": err.Error()})
		return
	}

	message := fmt.Sprintf("Updated model \"%s\"", vehicle.ID)
	c.JSON(http.StatusOK, gin.H{"data": vehicle, "status": "success", "message": message})
}

func DeleteVehicle(c *gin.Context) {
	vehicleService := service.GetVehicleService()

	id := c.Param("id")

	var vehicle model.Vehicle

	if err := model.DB.Where("id = ?", id).First(&vehicle).Error; err != nil {
		message := fmt.Sprintf("Could not find vehicle with [id:%s]", id)
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": message})
		return
	}

	if err := vehicleService.Delete(&vehicle); nil != err {
		message := fmt.Sprintf("An error occurred while deleting entry:\"%s\"", err.Error())
		c.JSON(http.StatusOK, gin.H{"data": vehicle, "status": "success", "message": message})
		return
	}

	message := fmt.Sprintf("Deleted vehicle \"%s\"", vehicle.ID)
	c.JSON(http.StatusOK, gin.H{"data": vehicle, "status": "success", "message": message})
}

func ValidateVehicleCategory(category model.VehicleCategory) (err error) {
	switch category {
	case model.PRODUCT_FLAT_TIRE_SERVICE:
		return nil
	case model.PRODUCT_FUEL_DELIVERY_SERVICE:
		return nil
	case model.PRODUCT_TIRE_AIR_SERVICE:
		return nil
	case model.PRODUCT_TOWING_SERVICE:
		return nil
	case model.PRODUCT_JUMP_START_SERVICE:
		return nil
	case model.PRODUCT_KEY_UNLOCK_SERVICE:
		return nil
	}
	return errors.New(fmt.Sprintf("Unsupported Vehicle Category \"%s\"", category))
}
