package controller

import (
	"fmt"
	"net/http"

	"github.com/nkamuo/rasta-server/dto"
	"github.com/nkamuo/rasta-server/model"

	"github.com/gin-gonic/gin"
)

// GET /products
// Get all products
func FindProducts(c *gin.Context) {
	var products []model.Product
	model.DB.Find(&products)

	c.JSON(http.StatusOK, gin.H{"data": products})
}

func CreateProduct(c *gin.Context) {
	// Validate input
	var input dto.CreateProductInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Create book
	product := model.Product{Title: input.Title, Description: input.Description}
	model.DB.Create(&product)

	c.JSON(http.StatusOK, gin.H{"data": product})
}

func FindProduct(c *gin.Context) {
	id := c.Param("id")

	var product model.Product

	if err := model.DB.Where("id = ?", id).First(&product).Error; err != nil {
		message := fmt.Sprintf("Could not find product with [id:%s]", id)
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": message})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": product})
}

func DeleteProduct(c *gin.Context) {
	id := c.Param("id")

	var product model.Product

	if err := model.DB.Where("id = ?", id).First(&product).Error; err != nil {
		message := fmt.Sprintf("Could not find product with [id:%s]", id)
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": message})
		return
	}
	model.DB.Delete(&product)
	message := fmt.Sprintf("Deleted product \"%s\"", product.Title)
	c.JSON(http.StatusOK, gin.H{"data": product, "status": "success", "message": message})
}
