package controller

import (
	"fmt"
	"net/http"

	"github.com/google/uuid"
	"github.com/nkamuo/rasta-server/data/financial"
	"github.com/nkamuo/rasta-server/data/pagination"
	"github.com/nkamuo/rasta-server/dto"
	"github.com/nkamuo/rasta-server/model"
	"github.com/nkamuo/rasta-server/service"
	"github.com/nkamuo/rasta-server/utils/auth"

	"github.com/gin-gonic/gin"
)

func FindRespondentEarnings(c *gin.Context) {
	var respondentService = service.GetRespondentService()
	var earnings []model.RespondentEarning
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

	respondent, err := respondentService.GetById(id, "User")
	if err != nil {
		var message string
		if err.Error() == "record not found" {
			message = fmt.Sprintf("Could not find respondent with [id:%s]", id)
			c.JSON(http.StatusNotFound, gin.H{"status": "error", "message": message})
		} else {
			message = fmt.Sprintf("An error occured: %s", err.Error())
			c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": message})
		}
		return
	}
	query := model.DB.Preload("Request.Product.Place").Preload("Request.Order.User")

	if *rUser.IsAdmin {

	} else {
		if respondent.User.ID.String() == rUser.ID.String() {
			query = query.Where("respondent_id = ?", respondent.ID)
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
	// query =
	if err := query.Scopes(pagination.Paginate(earnings, &page.Page, query)).Find(&earnings).Error; nil != err {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": err.Error()})
		return
	}
	page.Rows = earnings
	c.JSON(http.StatusOK, gin.H{"data": page, "status": "success"})
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
	var respondentService = service.GetRespondentService()
	earningService := service.GetRespondentEarningService()

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

	respondent, err := respondentService.GetById(*earning.RespondentID, "User")
	if err != nil {
		var message string
		if err.Error() == "record not found" {
			message = fmt.Sprintf("Could not find responder with [id:%s]", id)
			c.JSON(http.StatusNotFound, gin.H{"status": "error", "message": message})
		} else {
			message = fmt.Sprintf("An error occured: %s", err.Error())
			c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": message})
		}
		return
	}
	if *rUser.IsAdmin {

	} else {
		if respondent.User.ID.String() != rUser.ID.String() {
			message := fmt.Sprintf("Permision Denied: you may not access this resource")
			c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": message})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{"data": earning, "status": "success"})
}

func CommitRespondentEarning(c *gin.Context) {
	// var respondentService = service.GetRespondentService()
	earningService := service.GetRespondentEarningService()

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

	// respondent, err := respondentService.GetById(*earning.RespondentID, "User")
	// if err != nil {
	// 	var message string
	// 	if err.Error() == "record not found" {
	// 		message = fmt.Sprintf("Could not find responder with [id:%s]", id)
	// 		c.JSON(http.StatusNotFound, gin.H{"status": "error", "message": message})
	// 	} else {
	// 		message = fmt.Sprintf("An error occured: %s", err.Error())
	// 		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": message})
	// 	}
	// 	return
	// }
	if !*rUser.IsAdmin {
		message := fmt.Sprintf("Permision Denied: you may not access this resource")
		c.JSON(http.StatusForbidden, gin.H{"status": "error", "message": message})
		return
	}

	if err := earningService.Commit(earning); err != nil {
		message := fmt.Sprintf("An error occured: %s", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": message})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": earning, "status": "success"})
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
