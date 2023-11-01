package controller

import (
	"fmt"
	"net/http"

	// "github.com/mitchellh/mapstructure"

	"github.com/nkamuo/rasta-server/dto"
	"github.com/nkamuo/rasta-server/service"
	"github.com/nkamuo/rasta-server/utils/auth"

	"github.com/gin-gonic/gin"
	// "github.com/gin-gonic/gin"
)

func UpdateRespondentVehicle(c *gin.Context) {

	// placeService := service.GetPlaceService()
	// productService := service.GetProductService()
	respondentService := service.GetRespondentService()
	vehicleService := service.GetVehicleService()

	var input dto.RespondentVehicleSelectionInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": err.Error()})
		return
	}

	respondant, err := auth.GetCurrentRespondent(c, "User", "Vehicle")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": err.Error()})
		return
	}
	rUser := respondant.User

	vehicle, err := vehicleService.GetById(input.VehicleID)
	if nil != err {
		message := fmt.Sprintf("Could not resolve the specified Vehicle with [id:%s]: %s", input.VehicleID, err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"message": message, "status": "error"})
		return
	}

	if vehicle.OwnerID == nil || vehicle.OwnerID.String() != rUser.ID.String() {
		message := fmt.Sprintf("Don't have permission to assign this vehicle [id:%s]", input.VehicleID)
		c.JSON(http.StatusBadRequest, gin.H{"message": message, "status": "error"})
		return
	}

	rRespondant, err := respondentService.GetById(respondant.ID)
	if nil != err {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": err.Error()})
		return
	}
	rRespondant.VehicleID = &vehicle.ID

	if err := respondentService.Save(rRespondant); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "success", "data": vehicle})
	return

}

// ////////
// // FIND SESSION
// //
// /////////////
func FindRespondentVehicle(c *gin.Context) {
	respondant, err := auth.GetCurrentRespondent(c, "User", "Vehicle")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "success", "data": respondant.Vehicle})
}
