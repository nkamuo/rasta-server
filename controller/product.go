package controller

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/google/uuid"
	"github.com/nkamuo/rasta-server/dto"
	"github.com/nkamuo/rasta-server/model"
	"github.com/nkamuo/rasta-server/service"

	"github.com/gin-gonic/gin"
)

// GET /products
// Get all products
func FindProducts(c *gin.Context) {
	var products []model.Product
	if err := model.DB.Find(&products).Error; nil != err {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": products})
}

func CreateProduct(c *gin.Context) {

	productService := service.GetProductService()
	placeService := service.GetPlaceService()

	var input dto.ProductCreationInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := ValidateProductCategory(input.Category); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error(), "status": "error"})
		return
	}

	place, err := placeService.GetById(input.PlaceID)
	if nil != err {
		message := fmt.Sprintf("Could not resolve the specified Place with [id:%s]: %s", input.PlaceID, err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"message": message, "status": "error"})
		return
	}

	product := model.Product{
		Label:       input.Label,
		Title:       input.Title,
		Category:    input.Category,
		Description: input.Description,
		IconImage:   input.IconImage,
		CoverImage:  input.CoverImage,
		PlaceID:     place.ID,
	}
	if err := productService.Save(&product); nil != err {
		c.JSON(http.StatusOK, gin.H{"status": "error", "message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": product, "status": "success"})
}

func FindProduct(c *gin.Context) {
	productService := service.GetProductService()

	id, err := uuid.Parse(c.Param("id"))
	if nil != err {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "Invalid Id provided"})
		return
	}

	product, err := productService.GetById(id)
	if err != nil {
		message := fmt.Sprintf("Could not find product with [id:%s]", id)
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": message})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": product})
}

func UpdateProduct(c *gin.Context) {
	placeService := service.GetProductService()

	var requestBody map[string]interface{}
	if err := c.Copy().BindJSON(&requestBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "Invalid JSON format"})
		return
	}

	var input dto.ProductUpdateInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	id, err := uuid.Parse(c.Param("id"))
	if nil != err {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "Invalid Id provided"})
		return
	}
	place, err := placeService.GetById(id)
	if nil != err {
		message := fmt.Sprintf("Could not find place with [id:%s]", id)
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": message})
	}

	if _, ok := requestBody["label"]; ok {
		place.Label = input.Label
	}
	if _, ok := requestBody["iconImage"]; ok {
		place.IconImage = input.IconImage
	}
	if _, ok := requestBody["coverImage"]; ok {
		place.CoverImage = input.CoverImage
	}
	if _, ok := requestBody["title"]; ok {
		place.Title = input.Title
	}
	if _, ok := requestBody["description"]; ok {
		place.Description = input.Description
	}
	if _, ok := requestBody["description"]; ok {
		place.Description = input.Description
	}

	if err := placeService.Save(place); nil != err {
		c.JSON(http.StatusOK, gin.H{"status": "error", "message": err.Error()})
		return
	}

	message := fmt.Sprintf("Updated place \"%s\"", place.ID)
	c.JSON(http.StatusOK, gin.H{"data": place, "status": "success", "message": message})
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

func ValidateProductCategory(category model.ProductCategory) (err error) {
	switch category {
	case model.PRODUCT_FLAT_TIRE_SERVICE:
		return nil
	case model.PRODUCT_FUEL_DELIVERY_SERVICE:
		return nil
	case model.PRODUCT_TIRE_AIR_SERVICE:
		return nil
	case model.PRODUCT_TOWING_SERVICE:
		return nil
	case model.PRODUCT_JUMP_START_SERVICE:
		return nil
	case model.PRODUCT_KEY_UNLOCK_SERVICE:
		return nil
	}
	return errors.New(fmt.Sprintf("Unsupported Product Category \"%s\"", category))
}
