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
	placeRepo := repository.GetPlaceRepository()

	var products []model.Product
	var page pagination.Page
	if err := c.ShouldBindQuery(&page); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": err.Error()})
		return
	}

	query := model.DB.Preload("Place")

	if category := c.Query("category"); category != "" {
		err := model.ValidateProductCategory(category)
		if nil != err {
			message := fmt.Sprintf("%s", err.Error())
			c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": message})
			return
		}
		query = query.Where("category = ?", category)
	}

	if place_id := c.Query("place_id"); place_id != "" {
		placeID, err := uuid.Parse(place_id)
		if nil != err {
			message := fmt.Sprintf("Error parsing place_id[%s] query: %s", place_id, err.Error())
			c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": message})
			return
		}
		if _, err := placeRepo.GetById(placeID); err != nil {
			message := fmt.Sprintf("Could not find referenced place[%s]: %s", placeID, err.Error())
			c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": message})
			return
		}
		query = query.Where("place_id = ?", placeID)
	}

	if location_ref := c.Query("location"); location_ref != "" {
		location, err := service.GetLocationService().Search(location_ref)
		if nil != err {
			message := fmt.Sprintf("Error parsing loation[%s] query: %s", location_ref, err.Error())
			c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": message})
			return
		}
		place, err := placeRepo.GetByLocation(location)
		if err != nil {
			message := fmt.Sprintf("There was an error identifying your province[%s]. It might not be supported yet: %s", location.Address, err.Error())
			c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": message})
			return
		}
		query = query.Where("place_id = ?", place.ID)
	}

	if err := query.Scopes(pagination.Paginate(products, &page, query)).Preload("Place").Find(&products).Error; nil != err {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": err.Error()})
		return
	}
	page.Rows = products
	c.JSON(http.StatusOK, gin.H{"data": page})
}

func FindProductByCategoryAndLocation(c *gin.Context) {
	placeRepo := repository.GetPlaceRepository()

	query := model.DB.Preload("Place")

	if category := c.Query("category"); category != "" {
		err := model.ValidateProductCategory(category)
		if nil != err {
			message := fmt.Sprintf("%s", err.Error())
			c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": message})
			return
		}
		query = query.Where("category = ?", category)
	} else {
		message := fmt.Sprintf("\"category\" query parameter is required")
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": message})
		return
	}

	if location_ref := c.Query("location"); location_ref != "" {
		location, err := service.GetLocationService().Search(location_ref)
		if nil != err {
			message := fmt.Sprintf("Error parsing loation[%s] query: %s", location_ref, err.Error())
			c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": message})
			return
		}
		place, err := placeRepo.GetByLocation(location)
		if err != nil {
			message := fmt.Sprintf("There was an error identifying your province[%s]. It might not be supported yet: %s", location.Address, err.Error())
			c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": message})
			return
		}
		query = query.Where("place_id = ?", place.ID)
	} else {
		message := fmt.Sprintf("\"location\" query parameter is required")
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": message})
		return
	}

	var product model.Product

	err := query.First(&product).Error
	if err != nil {
		message := fmt.Sprintf("Could not find product for the location and category")
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": message})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "success", "data": product})
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
		Published:   &input.Published,
		Bundled:     &input.Bundled,
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

	product, err := productService.GetById(id, "Place")
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
		product.Published = input.Published
	}
	if nil != input.Bundled {
		product.Bundled = input.Bundled
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
