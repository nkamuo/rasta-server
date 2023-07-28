package controller

import (
	"fmt"
	"net/http"

	"github.com/google/uuid"
	"github.com/nkamuo/rasta-server/dto"
	"github.com/nkamuo/rasta-server/model"
	"github.com/nkamuo/rasta-server/repository"
	"github.com/nkamuo/rasta-server/service"

	"github.com/gin-gonic/gin"
)

func FindRespondents(c *gin.Context) {
	var respondents []model.Respondent
	model.DB.Find(&respondents)

	c.JSON(http.StatusOK, gin.H{"status": "success", "data": respondents})
}

func FindRespondentsByCompany(c *gin.Context) {

	companyService := service.GetCompanyService()

	companyID, err := uuid.Parse(c.Param("id"))
	if nil != err {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "Invalid companyId provided"})
		return
	}

	if _, err := companyService.GetById(companyID); nil != err {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": err.Error()})
		return

	}

	var respondents []model.Respondent
	model.DB.
		Joins("JOIN companies ON companies.id = respondents.company_id").
		Where("companies.id = ?", companyID).Find(&respondents)

	if respondents == nil {
		respondents = make([]model.Respondent, 0)
	}

	c.JSON(http.StatusOK, gin.H{"status": "success", "data": respondents})
}

func AddRespondentToCompany(c *gin.Context) {

	companyService := service.GetCompanyService()
	respondentService := service.GetRespondentService()

	companyID, err := uuid.Parse(c.Param("id"))
	if nil != err {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "Invalid companyId provided"})
		return
	}

	company, err := companyService.GetById(companyID)
	if nil != err {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": err.Error()})
		return
	}

	var input dto.RespondentCompanyAssignmentInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": err.Error()})
		return
	}

	respondant, err := respondentService.GetById(input.RespondentID)
	if nil != err {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": err.Error()})
		return
	}

	if err := respondentService.AssignToCompany(respondant, company); nil != err {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": err.Error()})
		return
	}
	c.JSON(http.StatusBadRequest, gin.H{"status": "success", "message": "Respondent assinged to company successfully"})
}

func RemoveRespondentFromCompany(c *gin.Context) {

	companyService := service.GetCompanyService()
	respondentService := service.GetRespondentService()

	companyID, err := uuid.Parse(c.Param("id"))
	if nil != err {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "Invalid companyId provided"})
		return
	}

	respondentID, err := uuid.Parse(c.Param("respondent_id"))
	if nil != err {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "Invalid companyId provided"})
		return
	}

	company, err := companyService.GetById(companyID)
	if nil != err {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": err.Error()})
		return
	}

	respondant, err := respondentService.GetById(respondentID)
	if nil != err {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": err.Error()})
		return
	}

	if err := respondentService.RemoveFromCompany(respondant, company); nil != err {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": err.Error()})
		return
	}
	c.JSON(http.StatusBadRequest, gin.H{"status": "success", "message": "Respondent removed from company successfully"})
}

func CreateRespondent(c *gin.Context) {

	userService := service.GetUserService()
	respondentService := service.GetRespondentService()
	respondentRepo := repository.GetRespondentRepository()
	// Validate input
	var input dto.RespondentCreationInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": err.Error()})
		return
	}

	user, err := userService.GetById(input.UserId)
	if nil != err {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": err.Error()})
		return
	}

	if _, err := respondentRepo.GetByUser(*user); nil == err {
		message := fmt.Sprintf("There is already a respondant account \"%s\" associated with this user", user.FullName())
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": message})
		return

	}

	// Create respondent
	respondent := model.Respondent{
		UserID: &user.ID,
	}

	// fmt.Printf("Input USer ID: %s\n user.ID: %s\n respondent.UserId: %s\n", input.UserId, user.ID, respondent.UserID)

	if err := respondentService.Save(&respondent); nil != err {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "success", "data": respondent})
}

func FindRespondent(c *gin.Context) {
	respondentService := service.GetRespondentService()

	id, err := uuid.Parse(c.Param("id"))
	if nil != err {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "Invalid Id provided"})
		return
	}
	respondent, err := respondentService.GetById(id)
	if nil != err {
		message := fmt.Sprintf("Could not find respondent with [id:%s]", id)
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": message})
	}
	c.JSON(http.StatusOK, gin.H{"status": "success", "data": respondent})
}

func UpdateRespondent(c *gin.Context) {
	respondentService := service.GetRespondentService()

	// Validate input
	var input dto.RespondentUpdateInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	id, err := uuid.Parse(c.Param("id"))
	if nil != err {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "Invalid Id provided"})
		return
	}
	respondent, err := respondentService.GetById(id)
	if nil != err {
		message := fmt.Sprintf("Could not find respondent with [id:%s]", id)
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": message})
	}

	if err := respondentService.Save(respondent); nil != err {
		c.JSON(http.StatusOK, gin.H{"status": "error", "message": err.Error()})
		return
	}

	message := fmt.Sprintf("Updated respondent \"%s\"", respondent.ID)
	c.JSON(http.StatusOK, gin.H{"data": respondent, "status": "success", "message": message})
}

func DeleteRespondent(c *gin.Context) {
	respondentService := service.GetRespondentService()

	id, err := uuid.Parse(c.Param("id"))
	if nil != err {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "Invalid Id provided"})
		return
	}
	respondent, err := respondentService.GetById(id)
	if nil != err {
		message := fmt.Sprintf("Could not find respondent with [id:%s]", id)
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": message})
	}

	if err := respondentService.Delete(respondent); nil != err {
		c.JSON(http.StatusOK, gin.H{"status": "error", "message": err.Error()})
		return
	}

	message := fmt.Sprintf("Deleted respondent \"%s\"", respondent.ID)
	c.JSON(http.StatusOK, gin.H{"data": respondent, "status": "success", "message": message})
}
