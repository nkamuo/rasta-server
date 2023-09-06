package controller

import (
	"fmt"
	"net/http"

	"github.com/google/uuid"
	"github.com/nkamuo/rasta-server/data/financial"
	"github.com/nkamuo/rasta-server/data/pagination"
	"github.com/nkamuo/rasta-server/dto"
	"github.com/nkamuo/rasta-server/repository"
	"github.com/nkamuo/rasta-server/utils/auth"

	// "github.com/nkamuo/rasta-server/dto"
	"github.com/nkamuo/rasta-server/model"
	"github.com/nkamuo/rasta-server/service"

	"github.com/gin-gonic/gin"
)

func FindCompanyWithdrawals(c *gin.Context) {
	var companyService = service.GetCompanyService()
	var withdrawals []model.CompanyWithdrawal
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
	query = query.Scopes(pagination.Paginate(withdrawals, &page.Page, query))
	if err := query.Find(&withdrawals).Error; nil != err {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": err.Error()})
		return
	}
	page.Rows = withdrawals
	c.JSON(http.StatusOK, gin.H{"data": page})
}

// CREATE A WITHDRAWAL
func CreateCompanyWithdrawal(c *gin.Context) {

	var companyService = service.GetCompanyService()
	var walletRepo = repository.GetCompanyWalletRepository()
	var withdrawalService = service.GetCompanyWithdrawalService()

	var input dto.WithdrawalRequest

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

	if *rUser.IsAdmin {

	} else {
		if company.OperatorUser.ID.String() == rUser.ID.String() {

		} else {
			message := fmt.Sprintf("Permision Denied: you may not access this resource")
			c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": message})
			return
		}
	}

	wallet, err := walletRepo.GetByCompany(*company)
	if err != nil {
		message := fmt.Sprintf("Error Resolving Wallet: %s", err.Error())
		c.JSON(http.StatusNotFound, gin.H{"status": "error", "message": message})
		return
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	description := input.Description
	if description == nil {
		var Desc = ""
		description = &Desc
	}
	withdrawal, err := withdrawalService.Init(*wallet, input.Amount, *description)
	if nil != err {
		c.JSON(http.StatusOK, gin.H{"status": "error", "message": err.Error()})
		return
	}
	if err := withdrawalService.Save(withdrawal); nil != err {
		c.JSON(http.StatusOK, gin.H{"status": "error", "message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": withdrawal, "status": "success"})
}

func FindCompanyWithdrawal(c *gin.Context) {
	var companyService = service.GetCompanyService()
	var walletService = service.GetCompanyWalletService()
	var withdrawalService = service.GetCompanyWithdrawalService()

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

	withdrawal, err := withdrawalService.GetById(id)
	if err != nil {
		message := fmt.Sprintf("Could not find withdrawal with [id:%s]", id)
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": message})
		return
	}

	wallet, err := walletService.GetById(*withdrawal.WalletID)
	if err != nil {
		message := fmt.Sprintf("Error Processing Request")
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": message})
		return
	}

	company, err := companyService.GetById(*wallet.CompanyID, "OperatorUser")
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

	c.JSON(http.StatusOK, gin.H{"data": withdrawal})
}

// func UpdateCompanyWithdrawal(c *gin.Context) {
// 	withdrawalService := service.GetCompanyWithdrawalService()

// 	var input dto.CompanyWithdrawalUpdateInput
// 	if err := c.ShouldBindJSON(&input); err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
// 		return
// 	}

// 	id, err := uuid.Parse(c.Param("id"))
// 	if nil != err {
// 		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "Invalid Id provided"})
// 		return
// 	}
// 	withdrawal, err := withdrawalService.GetById(id)
// 	if nil != err {
// 		message := fmt.Sprintf("Could not find withdrawal with [id:%s]", id)
// 		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": message})
// 	}

// 	if nil != input.Code {
// 		withdrawal.Code = *input.Code
// 	}
// 	if nil != input.Rate {
// 		withdrawal.Rate = *input.Rate
// 	}
// 	if nil != input.Title {
// 		withdrawal.Title = *input.Title
// 	}
// 	if nil != input.ShortName {
// 		withdrawal.ShortName = input.ShortName
// 	}
// 	if nil != input.Published {
// 		withdrawal.Published = input.Published
// 	}
// 	if nil != input.Description {
// 		withdrawal.Description = input.Description
// 	}

// 	if err := withdrawalService.Save(withdrawal); nil != err {
// 		c.JSON(http.StatusOK, gin.H{"status": "error", "message": err.Error()})
// 		return
// 	}

// 	message := fmt.Sprintf("Updated model \"%s\"", withdrawal.ID)
// 	c.JSON(http.StatusOK, gin.H{"data": withdrawal, "status": "success", "message": message})
// }

func DeleteCompanyWithdrawal(c *gin.Context) {
	id := c.Param("id")

	var withdrawal model.CompanyWithdrawal

	if err := model.DB.Where("id = ?", id).First(&withdrawal).Error; err != nil {
		message := fmt.Sprintf("Could not find withdrawal with [id:%s]", id)
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": message})
		return
	}
	model.DB.Delete(&withdrawal)
	message := fmt.Sprintf("Deleted withdrawal \"%s\"", withdrawal.ID)
	c.JSON(http.StatusOK, gin.H{"data": withdrawal, "status": "success", "message": message})
}
