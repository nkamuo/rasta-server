package controller

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/google/uuid"
	"github.com/nkamuo/rasta-server/data/pagination"
	"github.com/nkamuo/rasta-server/dto"
	"github.com/nkamuo/rasta-server/model"
	"github.com/nkamuo/rasta-server/service"

	"github.com/gin-gonic/gin"
)

func FindVehicleModels(c *gin.Context) {
	var vehicleModels []model.VehicleModel
	var page pagination.Page
	if err := c.ShouldBindQuery(&page); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": err.Error()})
		return
	}
	query := model.DB.Model(&model.VehicleModel{})

	if page.Search != "" {
		nameSearchQuery := strings.Join([]string{"%", page.Search, "%"}, "")
		query = query.Where("label LIKE ? OR title LIKE ?", nameSearchQuery, nameSearchQuery)
	}
	query = query.Scopes(pagination.Paginate(vehicleModels, &page, query))
	if err := query.Find(&vehicleModels).Error; nil != err {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": err.Error()})
		return
	}
	page.Rows = vehicleModels
	c.JSON(http.StatusOK, gin.H{"data": page})
}

func CreateVehicleModel(c *gin.Context) {

	vehicleModelService := service.GetVehicleModelService()

	var input dto.VehicleModelCreationInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := model.ValidateVehicleCategory(input.Category); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error(), "status": "error"})
		return
	}

	vehicleModel := model.VehicleModel{
		Label:       input.Label,
		Title:       input.Title,
		Category:    input.Category,
		Description: input.Description,
		IconImage:   input.IconImage,
		CoverImage:  input.CoverImage,
	}
	if err := vehicleModelService.Save(&vehicleModel); nil != err {
		message := fmt.Sprintf("An error occurred while saving entry: %s", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": message})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": vehicleModel, "status": "success"})
}

func FindVehicleModel(c *gin.Context) {
	vehicleModelService := service.GetVehicleModelService()

	id, err := uuid.Parse(c.Param("id"))
	if nil != err {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "Invalid Id provided"})
		return
	}

	vehicleModel, err := vehicleModelService.GetById(id)
	if err != nil {
		message := fmt.Sprintf("Could not find vehicleModel with [id:%s]", id)
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": message})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": vehicleModel})
}

func UpdateVehicleModel(c *gin.Context) {
	vehicleModelService := service.GetVehicleModelService()

	// var requestBody map[string]interface{}
	// if err := c.Copy().BindJSON(&requestBody); err != nil {
	// 	c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "Invalid JSON format"})
	// 	return
	// }

	var input dto.VehicleModelUpdateInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	id, err := uuid.Parse(c.Param("id"))
	if nil != err {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "Invalid Id provided"})
		return
	}
	vehicleModel, err := vehicleModelService.GetById(id)
	if nil != err {
		message := fmt.Sprintf("Could not find vehicleModel with [id:%s]", id)
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": message})
	}

	if nil != input.Label {
		vehicleModel.Label = *input.Label
	}
	if nil != input.IconImage {
		vehicleModel.IconImage = *input.IconImage
	}
	if nil != input.IconImage {
		vehicleModel.CoverImage = *input.CoverImage
	}
	if nil != input.Title {
		vehicleModel.Title = *input.Title
	}
	if nil != input.Description {
		vehicleModel.Description = *input.Description
	}

	if err := vehicleModelService.Save(vehicleModel); nil != err {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": err.Error()})
		return
	}

	message := fmt.Sprintf("Updated place \"%s\"", vehicleModel.ID)
	c.JSON(http.StatusOK, gin.H{"data": vehicleModel, "status": "success", "message": message})
}

func DeleteVehicleModel(c *gin.Context) {
	vehicleModelService := service.GetVehicleModelService()

	id := c.Param("id")

	var vehicleModel model.VehicleModel

	if err := model.DB.Where("id = ?", id).First(&vehicleModel).Error; err != nil {
		message := fmt.Sprintf("Could not find vehicleModel with [id:%s]", id)
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": message})
		return
	}

	if err := vehicleModelService.Delete(&vehicleModel); nil != err {
		message := fmt.Sprintf("An error occurred while deleting entry:\"%s\"", err.Error())
		c.JSON(http.StatusOK, gin.H{"data": vehicleModel, "status": "success", "message": message})
		return
	}

	message := fmt.Sprintf("Deleted vehicleModel \"%s\"", vehicleModel.Title)
	c.JSON(http.StatusOK, gin.H{"data": vehicleModel, "status": "success", "message": message})
}
