package controller

import (
	"fmt"
	"net/http"

	"github.com/google/uuid"
	"github.com/nkamuo/rasta-server/data/financial"
	"github.com/nkamuo/rasta-server/data/pagination"
	"github.com/nkamuo/rasta-server/dto"
	"github.com/nkamuo/rasta-server/utils/auth"

	// "github.com/nkamuo/rasta-server/dto"
	"github.com/nkamuo/rasta-server/model"
	"github.com/nkamuo/rasta-server/service"

	"github.com/gin-gonic/gin"
)

func FindCompanyEarnings(c *gin.Context) {
	var companyService = service.GetCompanyService()
	var earnings []model.CompanyEarning
	var page dto.FinancialPageRequest

	id, err := uuid.Parse(c.Param("id"))
	if nil != err {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "Invalid Id provided"})
		return
	}

	rUser, err := auth.GetCurrentUser(c)
	if err != nil {
		message := fmt.Sprintf("Authentication erro: %s", err.Error())
		c.JSON(http.StatusNotFound, gin.H{"status": "error", "message": message})
		return
	}

	company, err := companyService.GetById(id, "OperatorUser")
	if err != nil {
		var message string
		if err.Error() == "record not found" {
			message = fmt.Sprintf("Could not find company with [id:%s]", id)
			c.JSON(http.StatusNotFound, gin.H{"status": "error", "message": message})
		} else {
			message = fmt.Sprintf("An error occured: %s", err.Error())
			c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": message})
		}
		return
	}
	query := model.DB

	if *rUser.IsAdmin {
		if company_id := c.Query("company_id"); company_id != "" {
			query = query.Where("company_id = ?", company_id)
		}
	} else {
		if company.OperatorUser.ID.String() == rUser.ID.String() {
			query = query.Where("company_id = ?", company.ID)
		} else {
			message := fmt.Sprintf("Permision Denied: you may not access this resource")
			c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": message})
			return
		}
	}

	if err := c.ShouldBindQuery(&page); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": err.Error()})
		return
	}

	query = query.Scopes(financial.FilterRequest(nil, &page, query))
	query = query.Scopes(pagination.Paginate(earnings, &page.Page, query))
	if err := query.Find(&earnings).Error; nil != err {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": err.Error()})
		return
	}
	page.Rows = earnings
	c.JSON(http.StatusOK, gin.H{"data": page})
}

// func CreateCompanyEarning(c *gin.Context) {

// 	earningService := service.GetCompanyEarningService()

// 	var input dto.CompanyEarningCreationInput
// 	if err := c.ShouldBindJSON(&input); err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
// 		return
// 	}

// 	earning := model.CompanyEarning{
// 		Code:        input.Code,
// 		Title:       input.Title,
// 		ShortName:   &input.ShortName,
// 		Description: &input.Description,
// 		Rate:        input.Rate,
// 		Published:   &input.Published,
// 	}
// 	if err := earningService.Save(&earning); nil != err {
// 		c.JSON(http.StatusOK, gin.H{"status": "error", "message": err.Error()})
// 		return
// 	}

// 	c.JSON(http.StatusOK, gin.H{"data": earning, "status": "success"})
// }

func FindCompanyEarning(c *gin.Context) {
	var companyService = service.GetCompanyService()
	var earningService = service.GetCompanyEarningService()

	id, err := uuid.Parse(c.Param("id"))
	if nil != err {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "Invalid Id provided"})
		return
	}

	rUser, err := auth.GetCurrentUser(c)
	if err != nil {
		message := fmt.Sprintf("Authentication erro: %s", err.Error())
		c.JSON(http.StatusNotFound, gin.H{"status": "error", "message": message})
		return
	}

	earning, err := earningService.GetById(id)
	if err != nil {
		message := fmt.Sprintf("Could not find earning with [id:%s]", id)
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": message})
		return
	}

	company, err := companyService.GetById(*earning.CompanyID, "OperatorUser")
	if err != nil {
		var message string
		if err.Error() == "record not found" {
			message = fmt.Sprintf("Could not find company with [id:%s]", id)
			c.JSON(http.StatusNotFound, gin.H{"status": "error", "message": message})
		} else {
			message = fmt.Sprintf("An error occured: %s", err.Error())
			c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": message})
		}
		return
	}
	if *rUser.IsAdmin {

	} else {
		if company.OperatorUser.ID.String() != rUser.ID.String() {
			message := fmt.Sprintf("Permision Denied: you may not access this resource")
			c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": message})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{"data": earning})
}

// func UpdateCompanyEarning(c *gin.Context) {
// 	earningService := service.GetCompanyEarningService()

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
// 	earning, err := earningService.GetById(id)
// 	if nil != err {
// 		message := fmt.Sprintf("Could not find earning with [id:%s]", id)
// 		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": message})
// 	}

// 	if nil != input.Code {
// 		earning.Code = *input.Code
// 	}
// 	if nil != input.Rate {
// 		earning.Rate = *input.Rate
// 	}
// 	if nil != input.Title {
// 		earning.Title = *input.Title
// 	}
// 	if nil != input.ShortName {
// 		earning.ShortName = input.ShortName
// 	}
// 	if nil != input.Published {
// 		earning.Published = input.Published
// 	}
// 	if nil != input.Description {
// 		earning.Description = input.Description
// 	}

// 	if err := earningService.Save(earning); nil != err {
// 		c.JSON(http.StatusOK, gin.H{"status": "error", "message": err.Error()})
// 		return
// 	}

// 	message := fmt.Sprintf("Updated model \"%s\"", earning.ID)
// 	c.JSON(http.StatusOK, gin.H{"data": earning, "status": "success", "message": message})
// }

func DeleteCompanyEarning(c *gin.Context) {
	id := c.Param("id")

	var earning model.CompanyEarning

	if err := model.DB.Where("id = ?", id).First(&earning).Error; err != nil {
		message := fmt.Sprintf("Could not find earning with [id:%s]", id)
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": message})
		return
	}
	model.DB.Delete(&earning)
	message := fmt.Sprintf("Deleted earning \"%s\"", earning.ID)
	c.JSON(http.StatusOK, gin.H{"data": earning, "status": "success", "message": message})
}
