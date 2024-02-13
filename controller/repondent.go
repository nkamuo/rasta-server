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
	"github.com/nkamuo/rasta-server/utils/auth"

	"github.com/gin-gonic/gin"
)

func FindRespondents(c *gin.Context) {
	var respondents []model.Respondent
	var page pagination.Page
	if err := c.ShouldBindQuery(&page); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": err.Error()})
		return
	}

	if err := model.DB.Preload("User").Preload("Place").Scopes(pagination.Paginate(respondents, &page, model.DB)).Find(&respondents).Error; nil != err {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": err.Error()})
		return
	}
	page.Rows = respondents
	c.JSON(http.StatusOK, gin.H{"data": page})
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
	if err = model.DB.
		Joins("JOIN companies ON companies.id = respondents.company_id").
		Where("companies.id = ?", companyID).
		Preload("User").
		// Preload("Company").
		Find(&respondents).Error; nil != err {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": err.Error()})
		return
	}

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
	c.JSON(http.StatusOK, gin.H{"status": "success", "message": "Respondent assinged to company successfully"})
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
	placeService := service.GetPlaceService()
	companyService := service.GetCompanyService()
	vehicleService := service.GetVehicleService()
	respondentService := service.GetRespondentService()
	respondentRepo := repository.GetRespondentRepository()
	// Validate input
	var input dto.RespondentCreationInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": err.Error()})
		return
	}

	user, err := userService.GetById(input.UserID)
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

	if nil != input.VehicleID {
		if _, err := vehicleService.GetById(*input.VehicleID); nil != err {
			message := fmt.Sprintf("Could not find vehicle with [id:%s]: %s", input.VehicleID, err.Error())
			c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": message})
			return
		}
		respondent.VehicleID = input.VehicleID
	}

	if nil != input.CompanyID {
		if _, err := companyService.GetById(*input.CompanyID); nil != err {
			message := fmt.Sprintf("Could not find company with [id:%s]: %s", input.CompanyID, err.Error())
			c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": message})
			return
		}
		respondent.CompanyID = input.CompanyID
	}

	if nil != input.PlaceID {
		if _, err := placeService.GetById(*input.PlaceID); nil != err {
			message := fmt.Sprintf("Could not find place with [id:%s]: %s", input.PlaceID, err.Error())
			c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": message})
			return
		}
		respondent.PlaceID = input.PlaceID
	}

	if nil != input.Active {
		respondent.Active = *input.Active
	}

	// fmt.Printf("Input USer ID: %s\n user.ID: %s\n respondent.UserId: %s\n", input.UserId, user.ID, respondent.UserID)

	if err := respondentService.Save(&respondent); nil != err {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "success", "data": respondent})
}

func FindRespondent(c *gin.Context) {
	// respondentService := service.GetRespondentService()

	id, err := uuid.Parse(c.Param("id"))
	if nil != err {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "Invalid Id provided"})
		return
	}
	// respondent, err := respondentService.GetById(id)
	var respondent model.Respondent
	err = model.DB.Where("id = ?", id).Preload("User").Preload("Place", "AccessBalance", "AccessSubscription").Preload("Company").First(&respondent).Error
	if nil != err {
		message := fmt.Sprintf("Could not find respondent with [id:%s]: %s", id, err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": message})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "success", "data": respondent})
}

func GetCurrentRespondent(c *gin.Context) {
	respondentRepo := repository.GetRespondentRepository()

	user, err := auth.GetCurrentUser(c)
	if nil != err {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": err.Error()})
		return
	}
	respondent, err := respondentRepo.GetByUser(*user, "User", "Place", "Vehicle", "Company", "AccessBalance", "AccessSubscription")
	if nil != err {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "success", "data": respondent})

}

func UpdateRespondent(c *gin.Context) {
	placeService := service.GetPlaceService()
	companyService := service.GetCompanyService()
	vehicleService := service.GetVehicleService()
	respondentService := service.GetRespondentService()
	respondentRepo := repository.GetRespondentRepository()

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
		return
	}

	if nil != input.VehicleID {
		if _, err := vehicleService.GetById(*input.VehicleID); nil != err {
			message := fmt.Sprintf("Could not find vehicle with [id:%s]: %s", input.VehicleID, err.Error())
			c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": message})
			return
		}
		respondent.VehicleID = input.VehicleID
	}

	if nil != input.CompanyID {
		if _, err := companyService.GetById(*input.CompanyID); nil != err {
			message := fmt.Sprintf("Could not find company with [id:%s]: %s", input.CompanyID, err.Error())
			c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": message})
			return
		}
		respondent.CompanyID = input.CompanyID
	}

	if nil != input.PlaceID {
		if _, err := placeService.GetById(*input.PlaceID); nil != err {
			message := fmt.Sprintf("Could not find place with [id:%s]: %s", input.PlaceID, err.Error())
			c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": message})
			return
		}
		respondent.PlaceID = input.PlaceID
	}

	if nil != input.Active {
		respondent.Active = *input.Active
	}

	if err := respondentService.Save(respondent); nil != err {
		c.JSON(http.StatusOK, gin.H{"status": "error", "message": err.Error()})
		return
	}

	nRespondent, err := respondentRepo.GetById(respondent.ID, "Place", "User", "Vehicle")
	if nil != err {
		c.JSON(http.StatusOK, gin.H{"status": "error", "message": err.Error()})
		return
	}

	message := fmt.Sprintf("Updated respondent \"%s\"", respondent.ID)
	c.JSON(http.StatusOK, gin.H{"data": nRespondent, "status": "success", "message": message})
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
