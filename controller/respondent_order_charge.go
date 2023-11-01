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

func FindRespondentOrderCharges(c *gin.Context) {
	// var respondentService = service.GetRespondentService()
	var charges []model.RespondentOrderCharge
	var page dto.FinancialPageRequest

	rUser, err := auth.GetCurrentUser(c)
	if err != nil {
		message := fmt.Sprintf("Authentication erro: %s", err.Error())
		c.JSON(http.StatusNotFound, gin.H{"status": "error", "message": message})
		return
	}

	query := model.DB.Preload("Request.Product.Place").Preload("Request.Order.User").Preload("Respondent.User")

	if *rUser.IsAdmin {
		if respondent_id := c.Query("respondent_id"); respondent_id != "" {
			query = query.Where("respondent_id = ?", respondent_id)
		}
	} else {
		// if respondent.User.ID.String() == rUser.ID.String() {
		// 	// query = query.Where("respondent_id = ?", respondent.ID)
		// } else {
		message := fmt.Sprintf("Permision Denied: you may not access this resource")
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": message})
		return
		// }
	}

	if status := c.Query("status"); status != "" {
		query = query.Where("status = ?", status)
	}

	if err := c.ShouldBindQuery(&page); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": err.Error()})
		return
	}

	query = query.Scopes(financial.FilterRequest(nil, &page, query))
	// query =
	if err := query.Scopes(pagination.Paginate(charges, &page.Page, query)).Find(&charges).Error; nil != err {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": err.Error()})
		return
	}
	page.Rows = charges
	c.JSON(http.StatusOK, gin.H{"data": page, "status": "success"})
}

func FindRespondentOrderChargesByRespondent(c *gin.Context) {
	var respondentService = service.GetRespondentService()
	var charges []model.RespondentOrderCharge
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
		// if respondent_id := c.Query("respondent_id"); respondent_id != "" {
		// 	query = query.Where("respondent_id = ?", respondent_id)
		// }
	} else {
		if respondent.User.ID.String() == rUser.ID.String() {
			query = query.Where("respondent_id = ?", respondent.ID).Where("respondent_id = ?", respondent.ID)
		} else {
			message := fmt.Sprintf("Permision Denied: you may not access this resource")
			c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": message})
			return
		}
	}

	if status := c.Query("status"); status != "" {
		query = query.Where("status = ?", status)
	}

	if err := c.ShouldBindQuery(&page); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": err.Error()})
		return
	}

	query = query.Scopes(financial.FilterRequest(nil, &page, query))
	// query =
	if err := query.Scopes(pagination.Paginate(charges, &page.Page, query)).Find(&charges).Error; nil != err {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": err.Error()})
		return
	}
	page.Rows = charges
	c.JSON(http.StatusOK, gin.H{"data": page, "status": "success"})
}

// func CreateRespondentOrderCharge(c *gin.Context) {

// 	respondentChargeService := service.GetRespondentOrderChargeService()

// 	var input dto.RespondentOrderChargeCreationInput
// 	if err := c.ShouldBindJSON(&input); err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
// 		return
// 	}

// 	respondentCharge := model.RespondentOrderCharge{
// 		Code:        input.Code,
// 		Title:       input.Title,
// 		ShortName:   &input.ShortName,
// 		Description: &input.Description,
// 		Rate:        input.Rate,
// 		Published:   &input.Published,
// 	}
// 	if err := respondentChargeService.Save(&respondentCharge); nil != err {
// 		c.JSON(http.StatusOK, gin.H{"status": "error", "message": err.Error()})
// 		return
// 	}

// 	c.JSON(http.StatusOK, gin.H{"data": respondentCharge, "status": "success"})
// }

func FindRespondentOrderCharge(c *gin.Context) {
	var respondentService = service.GetRespondentService()
	chargeService := service.GetRespondentOrderChargeService()

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

	charge, err := chargeService.GetById(id)
	if err != nil {
		message := fmt.Sprintf("Could not find charge with [id:%s]", id)
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": message})
		return
	}

	respondent, err := respondentService.GetById(*charge.RespondentID, "User")
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

	c.JSON(http.StatusOK, gin.H{"data": charge, "status": "success"})
}

func CommitRespondentOrderCharge(c *gin.Context) {
	// var respondentService = service.GetRespondentService()
	chargeService := service.GetRespondentOrderChargeService()

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

	charge, err := chargeService.GetById(id)
	if err != nil {
		message := fmt.Sprintf("Could not find charge with [id:%s]", id)
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": message})
		return
	}

	// respondent, err := respondentService.GetById(*charge.RespondentID, "User")
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

	if err := chargeService.Commit(charge); err != nil {
		message := fmt.Sprintf("An error occured: %s", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": message})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": charge, "status": "success"})
}

// func UpdateRespondentOrderCharge(c *gin.Context) {
// 	respondentChargeService := service.GetRespondentOrderChargeService()

// 	var input dto.RespondentOrderChargeUpdateInput
// 	if err := c.ShouldBindJSON(&input); err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
// 		return
// 	}

// 	id, err := uuid.Parse(c.Param("id"))
// 	if nil != err {
// 		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "Invalid Id provided"})
// 		return
// 	}
// 	respondentCharge, err := respondentChargeService.GetById(id)
// 	if nil != err {
// 		message := fmt.Sprintf("Could not find respondentCharge with [id:%s]", id)
// 		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": message})
// 	}

// 	if nil != input.Code {
// 		respondentCharge.Code = *input.Code
// 	}
// 	if nil != input.Rate {
// 		respondentCharge.Rate = *input.Rate
// 	}
// 	if nil != input.Title {
// 		respondentCharge.Title = *input.Title
// 	}
// 	if nil != input.ShortName {
// 		respondentCharge.ShortName = input.ShortName
// 	}
// 	if nil != input.Published {
// 		respondentCharge.Published = input.Published
// 	}
// 	if nil != input.Description {
// 		respondentCharge.Description = input.Description
// 	}

// 	if err := respondentChargeService.Save(respondentCharge); nil != err {
// 		c.JSON(http.StatusOK, gin.H{"status": "error", "message": err.Error()})
// 		return
// 	}

// 	message := fmt.Sprintf("Updated model \"%s\"", respondentCharge.ID)
// 	c.JSON(http.StatusOK, gin.H{"data": respondentCharge, "status": "success", "message": message})
// }
