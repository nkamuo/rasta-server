package controller

import (
	"fmt"
	"net/http"

	"github.com/google/uuid"
	"github.com/nkamuo/rasta-server/data/pagination"

	// "github.com/nkamuo/rasta-server/dto"
	"github.com/nkamuo/rasta-server/model"
	"github.com/nkamuo/rasta-server/service"

	"github.com/gin-gonic/gin"
)

func FindCompanyEarnings(c *gin.Context) {
	var companyEarnings []model.CompanyEarning
	var page pagination.Page
	if err := c.ShouldBindQuery(&page); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": err.Error()})
		return
	}

	if err := model.DB.Scopes(pagination.Paginate(companyEarnings, &page, model.DB)).Find(&companyEarnings).Error; nil != err {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": err.Error()})
		return
	}
	page.Rows = companyEarnings
	c.JSON(http.StatusOK, gin.H{"data": page})
}

// func CreateCompanyEarning(c *gin.Context) {

// 	companyEarningService := service.GetCompanyEarningService()

// 	var input dto.CompanyEarningCreationInput
// 	if err := c.ShouldBindJSON(&input); err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
// 		return
// 	}

// 	companyEarning := model.CompanyEarning{
// 		Code:        input.Code,
// 		Title:       input.Title,
// 		ShortName:   &input.ShortName,
// 		Description: &input.Description,
// 		Rate:        input.Rate,
// 		Published:   &input.Published,
// 	}
// 	if err := companyEarningService.Save(&companyEarning); nil != err {
// 		c.JSON(http.StatusOK, gin.H{"status": "error", "message": err.Error()})
// 		return
// 	}

// 	c.JSON(http.StatusOK, gin.H{"data": companyEarning, "status": "success"})
// }

func FindCompanyEarning(c *gin.Context) {
	companyEarningService := service.GetCompanyEarningService()

	id, err := uuid.Parse(c.Param("id"))
	if nil != err {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "Invalid Id provided"})
		return
	}

	companyEarning, err := companyEarningService.GetById(id)
	if err != nil {
		message := fmt.Sprintf("Could not find companyEarning with [id:%s]", id)
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": message})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": companyEarning})
}

// func UpdateCompanyEarning(c *gin.Context) {
// 	companyEarningService := service.GetCompanyEarningService()

// 	var input dto.CompanyEarningUpdateInput
// 	if err := c.ShouldBindJSON(&input); err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
// 		return
// 	}

// 	id, err := uuid.Parse(c.Param("id"))
// 	if nil != err {
// 		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "Invalid Id provided"})
// 		return
// 	}
// 	companyEarning, err := companyEarningService.GetById(id)
// 	if nil != err {
// 		message := fmt.Sprintf("Could not find companyEarning with [id:%s]", id)
// 		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": message})
// 	}

// 	if nil != input.Code {
// 		companyEarning.Code = *input.Code
// 	}
// 	if nil != input.Rate {
// 		companyEarning.Rate = *input.Rate
// 	}
// 	if nil != input.Title {
// 		companyEarning.Title = *input.Title
// 	}
// 	if nil != input.ShortName {
// 		companyEarning.ShortName = input.ShortName
// 	}
// 	if nil != input.Published {
// 		companyEarning.Published = input.Published
// 	}
// 	if nil != input.Description {
// 		companyEarning.Description = input.Description
// 	}

// 	if err := companyEarningService.Save(companyEarning); nil != err {
// 		c.JSON(http.StatusOK, gin.H{"status": "error", "message": err.Error()})
// 		return
// 	}

// 	message := fmt.Sprintf("Updated model \"%s\"", companyEarning.ID)
// 	c.JSON(http.StatusOK, gin.H{"data": companyEarning, "status": "success", "message": message})
// }

func DeleteCompanyEarning(c *gin.Context) {
	id := c.Param("id")

	var companyEarning model.CompanyEarning

	if err := model.DB.Where("id = ?", id).First(&companyEarning).Error; err != nil {
		message := fmt.Sprintf("Could not find companyEarning with [id:%s]", id)
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": message})
		return
	}
	model.DB.Delete(&companyEarning)
	message := fmt.Sprintf("Deleted companyEarning \"%s\"", companyEarning.ID)
	c.JSON(http.StatusOK, gin.H{"data": companyEarning, "status": "success", "message": message})
}
