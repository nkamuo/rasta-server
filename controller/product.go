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

// GET /products
// Get all products
func FindProducts(c *gin.Context) {
	var products []model.Product
	var page pagination.Page
	if err := c.ShouldBindQuery(&page); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": err.Error()})
		return
	}

	if err := model.DB.Scopes(pagination.Paginate(products, &page, model.DB)).Preload("Place").Find(&products).Error; nil != err {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": err.Error()})
		return
	}
	page.Rows = products
	c.JSON(http.StatusOK, gin.H{"data": page})
}

func CreateProduct(c *gin.Context) {

	productRepository := repository.GetProductRepository()
	productService := service.GetProductService()
	placeService := service.GetPlaceService()

	var input dto.ProductCreationInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := model.ValidateProductCategory(input.Category); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error(), "status": "error"})
		return
	}

	place, err := placeService.GetById(input.PlaceID)
	if nil != err {
		message := fmt.Sprintf("Could not resolve the specified Place with [id:%s]: %s", input.PlaceID, err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"message": message, "status": "error"})
		return
	}

	existingProduct, err := productRepository.GetByPlaceIdAndCategory(input.PlaceID, input.Category)
	if nil != err {
		if err.Error() == "record not found" {
		} else {
			c.JSON(http.StatusBadRequest, gin.H{"message": err.Error(), "status": "error"})
			return
		}
	}
	if existingProduct != nil {
		message := fmt.Sprintf("Service %s is already provided at %s", input.Category, place.Name)
		c.JSON(http.StatusBadRequest, gin.H{"message": message, "status": "error"})
		return
	}

	product := model.Product{
		Published:   input.Published,
		Rate:        input.Rate,
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
	productService := service.GetProductService()

	// var requestBody map[string]interface{}
	// if err := c.Copy().BindJSON(&requestBody); err != nil {
	// 	c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "Invalid JSON format"})
	// 	return
	// }

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
	product, err := productService.GetById(id)
	if nil != err {
		message := fmt.Sprintf("Could not find product with [id:%s]", id)
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": message})
	}

	if nil != input.Published {
		product.Published = *input.Published
	}

	if nil != input.Rate {
		product.Rate = *input.Rate
	}

	if nil != input.Label {
		product.Label = *input.Label
	}
	if nil != input.IconImage {
		product.IconImage = *input.IconImage
	}
	if nil != input.IconImage {
		product.CoverImage = *input.CoverImage
	}
	if nil != input.Title {
		product.Title = *input.Title
	}
	if nil != input.Description {
		product.Description = *input.Description
	}

	if err := productService.Save(product); nil != err {
		c.JSON(http.StatusOK, gin.H{"status": "error", "message": err.Error()})
		return
	}

	message := fmt.Sprintf("Updated Product \"%s\"", product.ID)
	c.JSON(http.StatusOK, gin.H{"data": product, "status": "success", "message": message})
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
