package controller

import (
	"fmt"
	"net/http"

	"github.com/google/uuid"
	// "github.com/nkamuo/rasta-server/data/pagination"

	// "github.com/nkamuo/rasta-server/repository"
	"github.com/nkamuo/rasta-server/repository"
	"github.com/nkamuo/rasta-server/service"
	"github.com/nkamuo/rasta-server/utils/auth"

	"github.com/gin-gonic/gin"
)

// func FindRespondentWallets(c *gin.Context) {
// 	var respondentWallets []model.RespondentWallet
// 	var page pagination.Page
// 	if err := c.ShouldBindQuery(&page); err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": err.Error()})
// 		return
// 	}

// 	if err := model.DB.Scopes(pagination.Paginate(respondentWallets, &page, model.DB)).Find(&respondentWallets).Error; nil != err {
// 		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": err.Error()})
// 		return
// 	}
// 	page.Rows = respondentWallets
// 	c.JSON(http.StatusOK, gin.H{"data": page})
// }

// func CreateRespondentWallet(c *gin.Context) {

// 	respondentWalletService := service.GetRespondentWalletService()

// 	var input dto.RespondentWalletCreationInput
// 	if err := c.ShouldBindJSON(&input); err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
// 		return
// 	}

// 	respondentWallet := model.RespondentWallet{
// 		Code:        input.Code,
// 		Title:       input.Title,
// 		ShortName:   &input.ShortName,
// 		Description: &input.Description,
// 		Rate:        input.Rate,
// 		Published:   &input.Published,
// 	}
// 	if err := respondentWalletService.Save(&respondentWallet); nil != err {
// 		c.JSON(http.StatusOK, gin.H{"status": "error", "message": err.Error()})
// 		return
// 	}

// 	c.JSON(http.StatusOK, gin.H{"data": respondentWallet, "status": "success"})
// }

func FindRespondentWallet(c *gin.Context) {
	respondentService := service.GetRespondentService()
	respondentWalletRepository := repository.GetRespondentWalletRepository()

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

	if !*rUser.IsAdmin && respondent.User.ID.String() != rUser.ID.String() {
		message := fmt.Sprintf("Permision Denied: you may not access this resource")
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": message})
		return
	}

	respondentWallet, err := respondentWalletRepository.GetByRespondent(*respondent)
	if err != nil {
		message := fmt.Sprintf("Could not find Wallet for respondent \"%s\"", respondent.User.FullName())
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": message})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": respondentWallet})
}

// func UpdateRespondentWallet(c *gin.Context) {
// 	respondentWalletService := service.GetRespondentWalletService()

// 	var input dto.RespondentWalletUpdateInput
// 	if err := c.ShouldBindJSON(&input); err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
// 		return
// 	}

// 	id, err := uuid.Parse(c.Param("id"))
// 	if nil != err {
// 		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "Invalid Id provided"})
// 		return
// 	}
// 	respondentWallet, err := respondentWalletService.GetById(id)
// 	if nil != err {
// 		message := fmt.Sprintf("Could not find respondentWallet with [id:%s]", id)
// 		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": message})
// 	}

// 	if nil != input.Code {
// 		respondentWallet.Code = *input.Code
// 	}
// 	if nil != input.Rate {
// 		respondentWallet.Rate = *input.Rate
// 	}
// 	if nil != input.Title {
// 		respondentWallet.Title = *input.Title
// 	}
// 	if nil != input.ShortName {
// 		respondentWallet.ShortName = input.ShortName
// 	}
// 	if nil != input.Published {
// 		respondentWallet.Published = input.Published
// 	}
// 	if nil != input.Description {
// 		respondentWallet.Description = input.Description
// 	}

// 	if err := respondentWalletService.Save(respondentWallet); nil != err {
// 		c.JSON(http.StatusOK, gin.H{"status": "error", "message": err.Error()})
// 		return
// 	}

// 	message := fmt.Sprintf("Updated model \"%s\"", respondentWallet.ID)
// 	c.JSON(http.StatusOK, gin.H{"data": respondentWallet, "status": "success", "message": message})
// }

// func DeleteRespondentWallet(c *gin.Context) {
// 	id := c.Param("id")

// 	var respondentWallet model.RespondentWallet

// 	if err := model.DB.Where("id = ?", id).First(&respondentWallet).Error; err != nil {
// 		message := fmt.Sprintf("Could not find respondentWallet with [id:%s]", id)
// 		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": message})
// 		return
// 	}
// 	model.DB.Delete(&respondentWallet)
// 	message := fmt.Sprintf("Deleted respondentWallet \"%s\"", respondentWallet.ID)
// 	c.JSON(http.StatusOK, gin.H{"data": respondentWallet, "status": "success", "message": message})
// }
