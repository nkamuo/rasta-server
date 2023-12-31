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

func FindCompanies(c *gin.Context) {
	var companies []model.Company
	var page pagination.Page
	if err := c.ShouldBindQuery(&page); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": err.Error()})
		return
	}

	if err := model.DB.Scopes(pagination.Paginate(companies, &page, model.DB)).Find(&companies).Error; nil != err {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": err.Error()})
		return
	}
	page.Rows = companies
	c.JSON(http.StatusOK, gin.H{"data": page})
}

func CreateCompany(c *gin.Context) {

	userService := service.GetUserService()
	companyService := service.GetCompanyService()
	companyRepo := repository.GetCompanyRepository()
	// Validate input
	var input dto.CompanyCreationInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": err.Error()})
		return
	}

	user, err := userService.GetById(input.OperatorUserID)
	if nil != err {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": err.Error()})
		return
	}

	if _, err := companyRepo.GetByUser(*user); nil == err {
		message := fmt.Sprintf("There is already a respondant account \"%s\" associated with this user", user.FullName())
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": message})
		return

	}

	// Create company
	company := model.Company{
		LicenseNumber:  input.LicenseNumber,
		Title:          input.Title,
		Description:    input.Description,
		Category:       input.Category,
		OperatorUserID: user.ID,
		Active:         &input.Active,
		Published:      &input.Published,
	}

	// fmt.Printf("Input USer ID: %s\n user.ID: %s\n company.UserId: %s\n", input.UserId, user.ID, company.UserID)

	if err := companyService.Save(&company); nil != err {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "success", "data": company})
}

func FindCompany(c *gin.Context) {
	// companyService := service.GetCompanyService()

	id, err := uuid.Parse(c.Param("id"))
	if nil != err {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "Invalid Id provided"})
		return
	}
	// company, err := companyService.GetById(id)
	var company model.Company

	if err = model.DB.Where("id = ?", id).Preload("OperatorUser").First(&company).Error; nil != err {
		message := fmt.Sprintf("Could not find company with [id:%s]: %s", id, err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": message})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "success", "data": company})
}

func UpdateCompany(c *gin.Context) {
	userService := service.GetUserService()
	companyService := service.GetCompanyService()

	// Validate input
	var input dto.CompanyUpdateInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	id, err := uuid.Parse(c.Param("id"))
	if nil != err {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "Invalid Id provided"})
		return
	}
	company, err := companyService.GetById(id)
	if nil != err {
		message := fmt.Sprintf("Could not find company with [id:%s]", id)
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": message})
	}

	if nil != input.Active {
		company.Active = input.Active
	}
	if nil != input.Published {
		company.Published = input.Published
	}
	if nil != input.Title {
		company.Title = *input.Title
	}
	if nil != input.LicenseNumber {
		company.LicenseNumber = *input.LicenseNumber
	}
	if nil != input.Category {
		company.Category = *input.Category
	}
	if nil != input.OperatorUserID {
		user, err := userService.GetById(*input.OperatorUserID)
		if nil != err {
			c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": err.Error()})
			return
		}
		company.OperatorUserID = user.ID
	}

	if err := companyService.Save(company); nil != err {
		c.JSON(http.StatusOK, gin.H{"status": "error", "message": err.Error()})
		return
	}

	message := fmt.Sprintf("Updated company \"%s\"", company.ID)
	c.JSON(http.StatusOK, gin.H{"data": company, "status": "success", "message": message})
}

func DeleteCompany(c *gin.Context) {
	companyService := service.GetCompanyService()

	id, err := uuid.Parse(c.Param("id"))
	if nil != err {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "Invalid Id provided"})
		return
	}
	company, err := companyService.GetById(id)
	if nil != err {
		message := fmt.Sprintf("Could not find company with [id:%s]", id)
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": message})
	}

	if err := companyService.Delete(company); nil != err {
		c.JSON(http.StatusOK, gin.H{"status": "error", "message": err.Error()})
		return
	}

	message := fmt.Sprintf("Deleted company \"%s\"", company.ID)
	c.JSON(http.StatusOK, gin.H{"data": company, "status": "success", "message": message})
}
