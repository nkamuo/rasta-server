package controller

import (
	"fmt"
	"net/http"

	"github.com/google/uuid"
	"github.com/nkamuo/rasta-server/dto"
	"github.com/nkamuo/rasta-server/model"
	"github.com/nkamuo/rasta-server/service"

	"github.com/gin-gonic/gin"
)

func FindFuelTypes(c *gin.Context) {
	var fuelTypes []model.FuelType
	if err := model.DB.Find(&fuelTypes).Error; nil != err {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": fuelTypes})
}

func CreateFuelType(c *gin.Context) {

	fuelTypeService := service.GetFuelTypeService()

	var input dto.FuelTypeCreationInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	fuelType := model.FuelType{
		Code:        input.Code,
		Title:       input.Title,
		ShortName:   &input.ShortName,
		Description: &input.Description,
		Rate:        input.Rate,
		Published:   input.Published,
	}
	if err := fuelTypeService.Save(&fuelType); nil != err {
		c.JSON(http.StatusOK, gin.H{"status": "error", "message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": fuelType, "status": "success"})
}

func FindFuelType(c *gin.Context) {
	fuelTypeService := service.GetFuelTypeService()

	id, err := uuid.Parse(c.Param("id"))
	if nil != err {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "Invalid Id provided"})
		return
	}

	fuelType, err := fuelTypeService.GetById(id)
	if err != nil {
		message := fmt.Sprintf("Could not find fuelType with [id:%s]", id)
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": message})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": fuelType})
}

func UpdateFuelType(c *gin.Context) {
	fuelTypeService := service.GetFuelTypeService()

	var input dto.FuelTypeUpdateInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	id, err := uuid.Parse(c.Param("id"))
	if nil != err {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "Invalid Id provided"})
		return
	}
	fuelType, err := fuelTypeService.GetById(id)
	if nil != err {
		message := fmt.Sprintf("Could not find fuelType with [id:%s]", id)
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": message})
	}

	if nil != input.Code {
		fuelType.Code = *input.Code
	}
	if nil != input.Rate {
		fuelType.Rate = *input.Rate
	}
	if nil != input.Title {
		fuelType.Title = *input.Title
	}
	if nil != input.ShortName {
		fuelType.ShortName = input.ShortName
	}
	if nil != input.Published {
		fuelType.Published = *input.Published
	}
	if nil != input.Description {
		fuelType.Description = input.Description
	}

	if err := fuelTypeService.Save(fuelType); nil != err {
		c.JSON(http.StatusOK, gin.H{"status": "error", "message": err.Error()})
		return
	}

	message := fmt.Sprintf("Updated model \"%s\"", fuelType.ID)
	c.JSON(http.StatusOK, gin.H{"data": fuelType, "status": "success", "message": message})
}

func DeleteFuelType(c *gin.Context) {
	id := c.Param("id")

	var fuelType model.FuelType

	if err := model.DB.Where("id = ?", id).First(&fuelType).Error; err != nil {
		message := fmt.Sprintf("Could not find fuelType with [id:%s]", id)
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": message})
		return
	}
	model.DB.Delete(&fuelType)
	message := fmt.Sprintf("Deleted fuelType \"%s\"", fuelType.ID)
	c.JSON(http.StatusOK, gin.H{"data": fuelType, "status": "success", "message": message})
}
