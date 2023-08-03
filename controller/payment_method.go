package controller

import (
	"fmt"
	"net/http"

	"github.com/google/uuid"
	"github.com/nkamuo/rasta-server/data/pagination"
	"github.com/nkamuo/rasta-server/dto"
	"github.com/nkamuo/rasta-server/model"
	"github.com/nkamuo/rasta-server/service"

	"github.com/gin-gonic/gin"
)

func FindPaymentMethods(c *gin.Context) {
	var paymentMethodes []model.PaymentMethod
	var page pagination.Page
	if err := c.ShouldBindQuery(&page); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": err.Error()})
		return
	}

	if err := model.DB.Scopes(pagination.Paginate(paymentMethodes, &page, model.DB)).Find(&paymentMethodes).Error; nil != err {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": err.Error()})
		return
	}
	page.Rows = paymentMethodes
	c.JSON(http.StatusOK, gin.H{"data": page})
}

func FindUserPaymentMethods(c *gin.Context) {
	userService := service.GetUserService()

	userID, err := uuid.Parse(c.Param("id"))
	if nil != err {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "Invalid user id provided"})
		return
	}
	if _, err := userService.GetById(userID); nil != err {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": err.Error()})
		return
	}

	var paymentMethodes []model.PaymentMethod
	if err = model.DB.
		Joins("JOIN users ON users.id = payment_methods.user_id").
		Where("users.id = ?", userID).
		// Preload("User").
		Find(&paymentMethodes).Error; nil != err {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": paymentMethodes})
}

func CreatePaymentMethod(c *gin.Context) {

	userService := service.GetUserService()
	paymentMethodeService := service.GetPaymentMethodService()

	var input dto.PaymentMethodCreationInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := userService.GetById(input.UserID)
	if nil != err {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": err.Error()})
		return
	}

	if err := model.ValidatePaymentMethodCategory(input.Category); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error(), "status": "error"})
		return
	}

	paymentMethode := model.PaymentMethod{
		Category:    input.Category,
		Description: input.Description,
		Details:     input.Details,
		Active:      input.Active,
		UserID:      &user.ID,
	}
	if err := paymentMethodeService.Save(&paymentMethode); nil != err {
		c.JSON(http.StatusOK, gin.H{"status": "error", "message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": paymentMethode, "status": "success"})
}

func FindPaymentMethod(c *gin.Context) {
	paymentMethodeService := service.GetPaymentMethodService()

	id, err := uuid.Parse(c.Param("id"))
	if nil != err {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "Invalid Id provided"})
		return
	}

	paymentMethode, err := paymentMethodeService.GetById(id)
	if err != nil {
		message := fmt.Sprintf("Could not find paymentMethode with [id:%s]", id)
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": message})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": paymentMethode})
}

func UpdatePaymentMethod(c *gin.Context) {
	paymentMethodeService := service.GetPaymentMethodService()

	var input dto.PaymentMethodUpdateInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	id, err := uuid.Parse(c.Param("id"))
	if nil != err {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "Invalid Id provided"})
		return
	}
	paymentMethode, err := paymentMethodeService.GetById(id)
	if nil != err {
		message := fmt.Sprintf("Could not find paymentMethode with [id:%s]", id)
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": message})
	}

	if nil != input.Category {
		paymentMethode.Category = *input.Category
		if err := model.ValidatePaymentMethodCategory(*input.Category); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": err.Error(), "status": "error"})
			return
		}
	}
	if nil != input.Details {
		paymentMethode.Details = *input.Details
	}
	if nil != input.Active {
		paymentMethode.Active = *input.Active
	}
	if nil != input.Description {
		paymentMethode.Description = *input.Description
	}

	if err := paymentMethodeService.Save(paymentMethode); nil != err {
		c.JSON(http.StatusOK, gin.H{"status": "error", "message": err.Error()})
		return
	}

	message := fmt.Sprintf("Updated payment method \"%s\"", paymentMethode.ID)
	c.JSON(http.StatusOK, gin.H{"data": paymentMethode, "status": "success", "message": message})
}

func DeletePaymentMethod(c *gin.Context) {
	id := c.Param("id")

	var paymentMethode model.PaymentMethod

	if err := model.DB.Where("id = ?", id).First(&paymentMethode).Error; err != nil {
		message := fmt.Sprintf("Could not find paymentMethode with [id:%s]", id)
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": message})
		return
	}
	model.DB.Delete(&paymentMethode)
	message := fmt.Sprintf("Deleted paymentMethode \"%s\"", paymentMethode.ID)
	c.JSON(http.StatusOK, gin.H{"data": paymentMethode, "status": "success", "message": message})
}
