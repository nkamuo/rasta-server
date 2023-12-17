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
	"github.com/nkamuo/rasta-server/utils/auth"

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

	rUser, err := auth.GetCurrentUser(c)
	if nil != err {
		message := fmt.Sprintf("Authentication Error: %s", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": message})
		return
	}

	query := model.DB.Model(&model.Vehicle{}).Preload("Model").Preload("Owner")

	if *rUser.IsAdmin {
		if ownerID := c.Query("owner_id"); ownerID != "" {
			query = query.Where("owner_id = ?", ownerID)
		}
	} else {
		query = query.Where("owner_id = ?", rUser.ID)
	}

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
	companyService := service.GetCompanyService()
	vehicleService := service.GetVehicleService()
	modelService := service.GetVehicleModelService()

	var input dto.VehicleCreationInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": err.Error()})
		return
	}

	rUser, err := auth.GetCurrentUser(c)
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{"status": "error", "message": err.Error()})
		return
	}

	var makeName, modelName *string
	var vehicleModel *model.VehicleModel
	var owner *model.User
	var company *model.Company

	if input.ModelID != nil {
		vehicleModel, err = modelService.GetById(*input.ModelID)
		if nil != err {
			message := fmt.Sprintf("Could not resolve the specified Model with [id:%s]: %s", input.ModelID, err.Error())
			c.JSON(http.StatusBadRequest, gin.H{"message": message, "status": "error"})
			return
		}
	} else {
		if input.MakeName == nil || input.ModelName == nil {
			message := fmt.Sprintf("You must provide a valid %s or both %s and %s", "ModelID", "MakeName", "ModelName")
			c.JSON(http.StatusBadRequest, gin.H{"message": message, "status": "error"})
			return
		}
		makeName, modelName = input.MakeName, input.ModelName
	}

	if input.OwnerID != nil {
		if !*rUser.IsAdmin {
			message := "You are not allowed to provide {ownerId} for this request"
			c.JSON(http.StatusForbidden, gin.H{"status": "error", "message": message})
			return
		}
		owner, err = userService.GetById(*input.OwnerID)
		if nil != err {
			message := fmt.Sprintf("Could not resolve the specified User with [id:%s]: %s", input.OwnerID, err.Error())
			c.JSON(http.StatusBadRequest, gin.H{"message": message, "status": "error"})
			return
		}
	} else {
		if input.CompanyID != nil {
			company, err = companyService.GetById(*input.CompanyID)
			if nil != err {
				message := fmt.Sprintf("Could not resolve the specified Company with [id:%s]: %s", input.OwnerID, err.Error())
				c.JSON(http.StatusBadRequest, gin.H{"message": message, "status": "error"})
				return
			}
		} else {
			owner = rUser
		}

	}

	vehicle := model.Vehicle{
		LicensePlateNumber: *input.LicensePlateNumber,
		Color:              *input.Color,
		Description:        input.Description,
		Published:          &input.Published,
	}

	if owner != nil {
		vehicle.OwnerID = &owner.ID
	}

	if company != nil {
		vehicle.CompanyID = &company.ID
	}
	if vehicleModel != nil {
		vehicle.ModelID = &vehicleModel.ID
	} else {
		vehicle.MakeName, vehicle.ModelName = makeName, modelName
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

	rUser, err := auth.GetCurrentUser(c)

	if err != nil {
		message := fmt.Sprintf("Authentication Error: %s", err.Error())
		c.JSON(http.StatusUnauthorized, gin.H{"status": "error", "message": message})
		return
	}

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

	if !*rUser.IsAdmin && (vehicle.OwnerID == nil || rUser.ID.String() != vehicle.OwnerID.String()) {
		message := fmt.Sprintf("Unathorized: You may not access this resource")
		c.JSON(http.StatusForbidden, gin.H{"status": "error", "message": message})
		return
	}

	if nil != input.Color {
		vehicle.Color = *input.Color
	}

	if nil != input.ModelName {
		vehicle.ModelName = input.ModelName
	}

	if nil != input.MakeName {
		vehicle.MakeName = input.MakeName
	}

	if nil != input.LicensePlateNumber {
		vehicle.LicensePlateNumber = *input.LicensePlateNumber
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
		vehicle.ModelID = &vehicleModel.ID
	}

	if *rUser.IsAdmin {
		if nil != input.OwnerID {
			owner, err := userService.GetById(*input.OwnerID)
			if nil != err {
				message := fmt.Sprintf("Could not resolve the specified User with [id:%s]: %s", input.OwnerID, err.Error())
				c.JSON(http.StatusBadRequest, gin.H{"message": message, "status": "error"})
				return
			}
			vehicle.OwnerID = &owner.ID
		}

		if nil != input.Published {
			vehicle.Published = input.Published
		}

	} else {
		*vehicle.Published = false
		/**
		* This is to automatically disable each vehicle after user modification so that admins
		* can review and manually re-enable them.
		 */
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
