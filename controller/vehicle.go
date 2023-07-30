package controller

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/google/uuid"
	"github.com/nkamuo/rasta-server/dto"
	"github.com/nkamuo/rasta-server/model"
	"github.com/nkamuo/rasta-server/service"

	"github.com/gin-gonic/gin"
)

// GET /vehicles
// Get all vehicles
func FindVehicles(c *gin.Context) {
	var vehicles []model.Vehicle
	if err := model.DB.Find(&vehicles).Error; nil != err {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": vehicles})
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
		LicensePlaceNumber: *input.LicensePlaceNumber,
		Description:        input.Description,
		ModelID:            vehicleModel.ID,
		OwnerID:            owner.ID,
	}
	if err := vehicleService.Save(&vehicle); nil != err {
		c.JSON(http.StatusOK, gin.H{"status": "error", "message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": vehicle, "status": "success"})
}

func FindVehicle(c *gin.Context) {
	vehicleService := service.GetVehicleService()

	id, err := uuid.Parse(c.Param("id"))
	if nil != err {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "Invalid Id provided"})
		return
	}

	vehicle, err := vehicleService.GetById(id)
	if err != nil {
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

	if nil != input.LicensePlaceNumber {
		vehicle.LicensePlaceNumber = *input.LicensePlaceNumber
	}
	if nil != input.Description {
		vehicle.Description = *input.Description
	}

	if nil != input.ModelID {
		vehicleModel, err := modelService.GetById(*input.ModelID)
		if nil != err {
			message := fmt.Sprintf("Could not resolve the specified Model with [id:%s]: %s", input.ModelID, err.Error())
			c.JSON(http.StatusBadRequest, gin.H{"message": message, "status": "error"})
			return
		}
		vehicle.ModelID = vehicleModel.ID
	}

	if nil != input.OwnerID {
		owner, err := userService.GetById(*input.OwnerID)
		if nil != err {
			message := fmt.Sprintf("Could not resolve the specified User with [id:%s]: %s", input.OwnerID, err.Error())
			c.JSON(http.StatusBadRequest, gin.H{"message": message, "status": "error"})
			return
		}
		vehicle.OwnerID = owner.ID
	}

	if err := vehicleService.Save(vehicle); nil != err {
		c.JSON(http.StatusOK, gin.H{"status": "error", "message": err.Error()})
		return
	}

	message := fmt.Sprintf("Updated model \"%s\"", vehicle.ID)
	c.JSON(http.StatusOK, gin.H{"data": vehicle, "status": "success", "message": message})
}

func DeleteVehicle(c *gin.Context) {
	id := c.Param("id")

	var vehicle model.Vehicle

	if err := model.DB.Where("id = ?", id).First(&vehicle).Error; err != nil {
		message := fmt.Sprintf("Could not find vehicle with [id:%s]", id)
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": message})
		return
	}
	model.DB.Delete(&vehicle)
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
