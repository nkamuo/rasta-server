package controller

import (
	"fmt"
	"net/http"

	"github.com/google/uuid"
	"github.com/nkamuo/rasta-server/data/pagination"
	"github.com/nkamuo/rasta-server/dto"
	"github.com/nkamuo/rasta-server/model"
	"github.com/nkamuo/rasta-server/repository"
	"github.com/nkamuo/rasta-server/service"

	"github.com/gin-gonic/gin"
)

func FindMotoristRequestSituations(c *gin.Context) {

	SituationRepo := repository.GetMotoristRequestSituationRepository()
	// var motoristRequestSituations []model.MotoristRequestSituation
	var page pagination.Page
	if err := c.ShouldBindQuery(&page); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": err.Error()})
		return
	}
	// query := model.DB.Model(&model.MotoristRequestSituation{})

	// if page.Search != "" {
	// 	nameSearchQuery := strings.Join([]string{"%", page.Search, "%"}, "")
	// 	query = query.Where("label LIKE ? OR title LIKE ?", nameSearchQuery, nameSearchQuery)
	// }
	// query = query.Scopes(pagination.Paginate(motoristRequestSituations, &page, query))
	// if err := query.Find(&motoristRequestSituations).Error; nil != err {
	// 	c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": err.Error()})
	// 	return
	// }
	if Situations, total, err := SituationRepo.FindAllDefault(); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": err.Error()})
		return
	} else {
		page.Rows = Situations
		page.TotalRows = total
	}
	c.JSON(http.StatusOK, gin.H{"data": page})
}

func CreateMotoristRequestSituation(c *gin.Context) {

	motoristRequestSituationService := service.GetMotoristRequestSituationService()

	var input dto.MotoristRequestSituationCreationInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	motoristRequestSituation := model.MotoristRequestSituation{
		Code:        input.Code,
		Title:       input.Title,
		SubTitlte:   input.SubTitlte,
		Note:        input.Note,
		Description: input.Description,
	}
	if err := motoristRequestSituationService.Save(&motoristRequestSituation); nil != err {
		message := fmt.Sprintf("An error occurred while saving entry: %s", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": message})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": motoristRequestSituation, "status": "success"})
}

func FindMotoristRequestSituation(c *gin.Context) {
	motoristRequestSituationService := service.GetMotoristRequestSituationService()

	id, err := uuid.Parse(c.Param("id"))
	if nil != err {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "Invalid Id provided"})
		return
	}

	motoristRequestSituation, err := motoristRequestSituationService.GetById(id)
	if err != nil {
		message := fmt.Sprintf("Could not find motoristRequestSituation with [id:%s]", id)
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": message})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": motoristRequestSituation})
}

func UpdateMotoristRequestSituation(c *gin.Context) {
	motoristRequestSituationService := service.GetMotoristRequestSituationService()

	// var requestBody map[string]interface{}
	// if err := c.Copy().BindJSON(&requestBody); err != nil {
	// 	c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "Invalid JSON format"})
	// 	return
	// }

	var input dto.MotoristRequestSituationUpdateInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	id, err := uuid.Parse(c.Param("id"))
	if nil != err {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "Invalid Id provided"})
		return
	}
	motoristRequestSituation, err := motoristRequestSituationService.GetById(id)
	if nil != err {
		message := fmt.Sprintf("Could not find motoristRequestSituation with [id:%s]", id)
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": message})
	}

	if nil != input.Title {
		motoristRequestSituation.Title = *input.Title
	}
	if nil != input.SubTitlte {
		motoristRequestSituation.SubTitlte = *input.SubTitlte
	}
	if nil != input.Note {
		motoristRequestSituation.Note = *input.Note
	}
	if nil != input.Description {
		motoristRequestSituation.Description = *input.Description
	}

	if err := motoristRequestSituationService.Save(motoristRequestSituation); nil != err {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": err.Error()})
		return
	}

	message := fmt.Sprintf("Updated place \"%s\"", motoristRequestSituation.ID)
	c.JSON(http.StatusOK, gin.H{"data": motoristRequestSituation, "status": "success", "message": message})
}

func DeleteMotoristRequestSituation(c *gin.Context) {
	motoristRequestSituationService := service.GetMotoristRequestSituationService()

	id := c.Param("id")

	var motoristRequestSituation model.MotoristRequestSituation

	if err := model.DB.Where("id = ?", id).First(&motoristRequestSituation).Error; err != nil {
		message := fmt.Sprintf("Could not find motoristRequestSituation with [id:%s]", id)
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": message})
		return
	}

	if err := motoristRequestSituationService.Delete(&motoristRequestSituation); nil != err {
		message := fmt.Sprintf("An error occurred while deleting entry:\"%s\"", err.Error())
		c.JSON(http.StatusOK, gin.H{"data": motoristRequestSituation, "status": "success", "message": message})
		return
	}

	message := fmt.Sprintf("Deleted motoristRequestSituation \"%s\"", motoristRequestSituation.Title)
	c.JSON(http.StatusOK, gin.H{"data": motoristRequestSituation, "status": "success", "message": message})
}
