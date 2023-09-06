package controller

import (
	"fmt"
	"net/http"

	"github.com/google/uuid"
	// "github.com/nkamuo/rasta-server/repository"
	"github.com/nkamuo/rasta-server/repository"
	"github.com/nkamuo/rasta-server/utils/auth"

	// "github.com/nkamuo/rasta-server/dto"

	"github.com/nkamuo/rasta-server/service"

	"github.com/gin-gonic/gin"
)

// func FindCompanyWallets(c *gin.Context) {
// 	var companyWallets []model.CompanyWallet
// 	var page pagination.Page
// 	if err := c.ShouldBindQuery(&page); err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": err.Error()})
// 		return
// 	}

// 	if err := model.DB.Scopes(pagination.Paginate(companyWallets, &page, model.DB)).Find(&companyWallets).Error; nil != err {
// 		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": err.Error()})
// 		return
// 	}
// 	page.Rows = companyWallets
// 	c.JSON(http.StatusOK, gin.H{"data": page})
// }

// func CreateCompanyWallet(c *gin.Context) {

// 	companyWalletService := service.GetCompanyWalletService()

// 	var input dto.CompanyWalletCreationInput
// 	if err := c.ShouldBindJSON(&input); err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
// 		return
// 	}

// 	companyWallet := model.CompanyWallet{
// 		Code:        input.Code,
// 		Title:       input.Title,
// 		ShortName:   &input.ShortName,
// 		Description: &input.Description,
// 		Rate:        input.Rate,
// 		Published:   &input.Published,
// 	}
// 	if err := companyWalletService.Save(&companyWallet); nil != err {
// 		c.JSON(http.StatusOK, gin.H{"status": "error", "message": err.Error()})
// 		return
// 	}

// 	c.JSON(http.StatusOK, gin.H{"data": companyWallet, "status": "success"})
// }

func FindCompanyWallet(c *gin.Context) {
	companyService := service.GetCompanyService()
	companyWalletRepository := repository.GetCompanyWalletRepository()

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

	if !*rUser.IsAdmin && company.OperatorUser.ID.String() != rUser.ID.String() {
		message := fmt.Sprintf("Permision Denied: you may not access this resource")
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": message})
		return
	}

	companyWallet, err := companyWalletRepository.GetByCompany(*company)
	if err != nil {
		message := fmt.Sprintf("Could not find Wallet for company \"%s\"", company.Title)
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": message})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": companyWallet})
}

// func UpdateCompanyWallet(c *gin.Context) {
// 	companyWalletService := service.GetCompanyWalletService()

// 	var input dto.CompanyWalletUpdateInput
// 	if err := c.ShouldBindJSON(&input); err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
// 		return
// 	}

// 	id, err := uuid.Parse(c.Param("id"))
// 	if nil != err {
// 		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "Invalid Id provided"})
// 		return
// 	}
// 	companyWallet, err := companyWalletService.GetById(id)
// 	if nil != err {
// 		message := fmt.Sprintf("Could not find companyWallet with [id:%s]", id)
// 		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": message})
// 	}

// 	if nil != input.Code {
// 		companyWallet.Code = *input.Code
// 	}
// 	if nil != input.Rate {
// 		companyWallet.Rate = *input.Rate
// 	}
// 	if nil != input.Title {
// 		companyWallet.Title = *input.Title
// 	}
// 	if nil != input.ShortName {
// 		companyWallet.ShortName = input.ShortName
// 	}
// 	if nil != input.Published {
// 		companyWallet.Published = input.Published
// 	}
// 	if nil != input.Description {
// 		companyWallet.Description = input.Description
// 	}

// 	if err := companyWalletService.Save(companyWallet); nil != err {
// 		c.JSON(http.StatusOK, gin.H{"status": "error", "message": err.Error()})
// 		return
// 	}

// 	message := fmt.Sprintf("Updated model \"%s\"", companyWallet.ID)
// 	c.JSON(http.StatusOK, gin.H{"data": companyWallet, "status": "success", "message": message})
// }

// func DeleteCompanyWallet(c *gin.Context) {
// 	id := c.Param("id")

// 	var companyWallet model.CompanyWallet

// 	if err := model.DB.Where("id = ?", id).First(&companyWallet).Error; err != nil {
// 		message := fmt.Sprintf("Could not find companyWallet with [id:%s]", id)
// 		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": message})
// 		return
// 	}
// 	model.DB.Delete(&companyWallet)
// 	message := fmt.Sprintf("Deleted companyWallet \"%s\"", companyWallet.ID)
// 	c.JSON(http.StatusOK, gin.H{"data": companyWallet, "status": "success", "message": message})
// }
