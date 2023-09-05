package controller

import (
	"fmt"
	"net/http"

	"github.com/google/uuid"
	"github.com/nkamuo/rasta-server/data/pagination"
	"github.com/nkamuo/rasta-server/model"
	"github.com/nkamuo/rasta-server/service"

	"github.com/gin-gonic/gin"
)

func FindRespondentEarnings(c *gin.Context) {
	var respondentEarnings []model.RespondentEarning
	var page pagination.Page
	if err := c.ShouldBindQuery(&page); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": err.Error()})
		return
	}

	if err := model.DB.Scopes(pagination.Paginate(respondentEarnings, &page, model.DB)).Find(&respondentEarnings).Error; nil != err {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": err.Error()})
		return
	}
	page.Rows = respondentEarnings
	c.JSON(http.StatusOK, gin.H{"data": page})
}

// func CreateRespondentEarning(c *gin.Context) {

// 	respondentEarningService := service.GetRespondentEarningService()

// 	var input dto.RespondentEarningCreationInput
// 	if err := c.ShouldBindJSON(&input); err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
// 		return
// 	}

// 	respondentEarning := model.RespondentEarning{
// 		Code:        input.Code,
// 		Title:       input.Title,
// 		ShortName:   &input.ShortName,
// 		Description: &input.Description,
// 		Rate:        input.Rate,
// 		Published:   &input.Published,
// 	}
// 	if err := respondentEarningService.Save(&respondentEarning); nil != err {
// 		c.JSON(http.StatusOK, gin.H{"status": "error", "message": err.Error()})
// 		return
// 	}

// 	c.JSON(http.StatusOK, gin.H{"data": respondentEarning, "status": "success"})
// }

func FindRespondentEarning(c *gin.Context) {
	respondentEarningService := service.GetRespondentEarningService()

	id, err := uuid.Parse(c.Param("id"))
	if nil != err {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "Invalid Id provided"})
		return
	}

	respondentEarning, err := respondentEarningService.GetById(id)
	if err != nil {
		message := fmt.Sprintf("Could not find respondentEarning with [id:%s]", id)
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": message})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": respondentEarning})
}

// func UpdateRespondentEarning(c *gin.Context) {
// 	respondentEarningService := service.GetRespondentEarningService()

// 	var input dto.RespondentEarningUpdateInput
// 	if err := c.ShouldBindJSON(&input); err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
// 		return
// 	}

// 	id, err := uuid.Parse(c.Param("id"))
// 	if nil != err {
// 		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "Invalid Id provided"})
// 		return
// 	}
// 	respondentEarning, err := respondentEarningService.GetById(id)
// 	if nil != err {
// 		message := fmt.Sprintf("Could not find respondentEarning with [id:%s]", id)
// 		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": message})
// 	}

// 	if nil != input.Code {
// 		respondentEarning.Code = *input.Code
// 	}
// 	if nil != input.Rate {
// 		respondentEarning.Rate = *input.Rate
// 	}
// 	if nil != input.Title {
// 		respondentEarning.Title = *input.Title
// 	}
// 	if nil != input.ShortName {
// 		respondentEarning.ShortName = input.ShortName
// 	}
// 	if nil != input.Published {
// 		respondentEarning.Published = input.Published
// 	}
// 	if nil != input.Description {
// 		respondentEarning.Description = input.Description
// 	}

// 	if err := respondentEarningService.Save(respondentEarning); nil != err {
// 		c.JSON(http.StatusOK, gin.H{"status": "error", "message": err.Error()})
// 		return
// 	}

// 	message := fmt.Sprintf("Updated model \"%s\"", respondentEarning.ID)
// 	c.JSON(http.StatusOK, gin.H{"data": respondentEarning, "status": "success", "message": message})
// }

func DeleteRespondentEarning(c *gin.Context) {
	id := c.Param("id")

	var respondentEarning model.RespondentEarning

	if err := model.DB.Where("id = ?", id).First(&respondentEarning).Error; err != nil {
		message := fmt.Sprintf("Could not find respondentEarning with [id:%s]", id)
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": message})
		return
	}
	model.DB.Delete(&respondentEarning)
	message := fmt.Sprintf("Deleted respondentEarning \"%s\"", respondentEarning.ID)
	c.JSON(http.StatusOK, gin.H{"data": respondentEarning, "status": "success", "message": message})
}
