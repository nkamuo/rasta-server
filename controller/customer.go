package controller

import (
	"fmt"
	"net/http"

	"github.com/nkamuo/rasta-server/dto"
	"github.com/nkamuo/rasta-server/model"

	"github.com/gin-gonic/gin"
)

// GET /Customers
// Get all Customers
func FindCustomers(c *gin.Context) {
	var Customers []model.Customer
	model.DB.Find(&Customers)

	c.JSON(http.StatusOK, gin.H{"data": Customers})
}

func CreateCustomer(c *gin.Context) {
	// Validate input
	var input dto.CreateCustomerInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Create book
	Customer := model.Customer{UserID: input.UserId}
	model.DB.Create(&Customer)

	c.JSON(http.StatusOK, gin.H{"data": Customer})
}

func FindCustomer(c *gin.Context) {
	id := c.Param("id")

	var Customer model.Customer

	if err := model.DB.Where("id = ?", id).First(&Customer).Error; err != nil {
		message := fmt.Sprintf("Could not find Customer with [id:%s]", id)
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": message})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": Customer})
}

func DeleteCustomer(c *gin.Context) {
	id := c.Param("id")

	var Customer model.Customer

	if err := model.DB.Where("id = ?", id).First(&Customer).Error; err != nil {
		message := fmt.Sprintf("Could not find Customer with [id:%s]", id)
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": message})
		return
	}
	model.DB.Model(&Customer).Delete(Customer)
	message := fmt.Sprintf("Deleted Customer \"%s\"", Customer.User.Email)
	c.JSON(http.StatusOK, gin.H{"data": Customer, "status": "success", "message": message})
}
