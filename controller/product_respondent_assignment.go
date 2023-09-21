package controller

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/google/uuid"
	// "github.com/mitchellh/mapstructure"
	"github.com/nkamuo/rasta-server/data/pagination"
	"github.com/nkamuo/rasta-server/dto"
	"github.com/nkamuo/rasta-server/model"
	"github.com/nkamuo/rasta-server/repository"
	"github.com/nkamuo/rasta-server/service"
	"github.com/nkamuo/rasta-server/utils/auth"

	"github.com/gin-gonic/gin"
	// "github.com/gin-gonic/gin"
)

func FindProductRespondentAssignments(c *gin.Context) {

	respondentRepo := repository.GetRespondentRepository()
	placeRepo := repository.GetPlaceRepository()

	var assignments []model.ProductRespondentAssignment
	var page pagination.Page
	if err := c.ShouldBindQuery(&page); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": err.Error()})
		return
	}

	query := model.DB.Preload("Respondent").Preload("Product.Place") //.Preload("Place")

	if product_id := c.Query("product_id"); product_id != "" {
		productID, err := uuid.Parse(product_id)
		if nil != err {
			message := fmt.Sprintf("Error parsing product_id[%s] query: %s", product_id, err.Error())
			c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": message})
			return
		}
		if _, err := placeRepo.GetById(productID); err != nil {
			message := fmt.Sprintf("Could not find referenced place[%s]: %s", productID, err.Error())
			c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": message})
			return
		}
		query = query.Where("product_id = ?", productID)
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
		query = query.Joins("Product").Where("place_id = ?", placeID)
	}

	if respondent_id := c.Query("respondent_id"); respondent_id != "" {
		respondentID, err := uuid.Parse(respondent_id)
		if nil != err {
			message := fmt.Sprintf("Error parsing respondent_id[%s] query: %s", respondent_id, err.Error())
			c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": message})
			return
		}
		if _, err := respondentRepo.GetById(respondentID); err != nil {
			message := fmt.Sprintf("Could not find referenced Respondent[%s]: %s", respondentID, err.Error())
			c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": message})
			return
		}
		query = query.Where("respondent_id = ?", respondentID)
	}

	if err := query.Scopes(pagination.Paginate(assignments, &page, query)).Find(&assignments).Error; nil != err {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": err.Error()})
		return
	}
	page.Rows = assignments
	c.JSON(http.StatusOK, gin.H{"data": page})
}

func CreateProductRespondentAssignment(c *gin.Context) {

	userService := service.GetUserService()
	placeService := service.GetPlaceService()
	productService := service.GetProductService()
	respondentService := service.GetRespondentService()
	assignmentService := service.GetProductRespondentAssignmentService()

	rUser, err := auth.GetCurrentUser(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": err.Error()})
		return
	}

	var input dto.ProductRespondentAssignmentCreationInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": err.Error()})
		return
	}

	product, err := productService.GetById(input.ProductID)
	if nil != err {
		message := fmt.Sprintf("Could not resolve the specified Product with [id:%s]: %s", input.ProductID, err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"message": message, "status": "error"})
		return
	}

	place, err := placeService.GetById(product.PlaceID)
	if nil != err {
		message := fmt.Sprintf("Could not resolve Place [id:%s]: %s", product.PlaceID, err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"message": message, "status": "error"})
		return
	}

	respondant, err := respondentService.GetById(input.RespondentID)
	if nil != err {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": err.Error()})
		return
	}

	user, err := userService.GetById(*respondant.UserID)
	if nil != err {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": err.Error()})
		return
	}

	var crntAssignment model.ProductRespondentAssignment
	var existingCount int64
	err = model.DB.Where("product_id = ? AND respondent_id = ?", product.ID, respondant.ID).Model(&crntAssignment).Count(&existingCount).Error
	fmt.Printf("Current Assignment: %#v\n", crntAssignment)
	if nil != err {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": err.Error()})
		return

	}

	if existingCount > 0 {
		message := fmt.Sprintf(
			"There is already an assignment for \"%v\" on \"%v\" in \"%v\"",
			user.FullName(),
			product.Label,
			place.Name,
		)
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": message})
		return

	}

	// Create assignment
	assignment := model.ProductRespondentAssignment{
		ProductID:    &product.ID,
		RespondentID: &respondant.ID,
		Note:         input.Note,
		Description:  input.Description,
	}

	if *rUser.IsAdmin {
		assignment.Active = &input.Active
		assignment.AllowRespondentActivate = &input.AllowRespondentActivate
	} else {

	}

	// fmt.Printf("Input USer ID: %s\n user.ID: %s\n assignment.UserId: %s\n", input.UserId, user.ID, assignment.UserID)

	if err := assignmentService.Save(&assignment); nil != err {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "success", "data": assignment})
}

func FindProductRespondentAssignment(c *gin.Context) {
	assignmentService := service.GetProductRespondentAssignmentService()

	id, err := uuid.Parse(c.Param("id"))
	if nil != err {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "Invalid Id provided"})
		return
	}
	assignment, err := assignmentService.GetById(id)
	if nil != err {
		message := fmt.Sprintf("Could not find assignment with [id:%s]", id)
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": message})
	}
	c.JSON(http.StatusOK, gin.H{"status": "success", "data": assignment})
}

func UpdateProductRespondentAssignment(c *gin.Context) {
	assignmentService := service.GetProductRespondentAssignmentService()

	var input dto.ProductRespondentAssignmentUpdateInput
	if err := c.Copy().ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": err.Error()})
		return
	}

	// var input dto.ProductRespondentAssignmentUpdateInput
	// mapstructure.Decode(requestBody, input)

	requestingUser, err := auth.GetCurrentUser(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": err.Error()})
		return
	}

	id, err := uuid.Parse(c.Param("id"))
	if nil != err {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "Invalid Id provided"})
		return
	}
	assignment, err := assignmentService.GetById(id)
	if nil != err {
		message := fmt.Sprintf("Could not find assignment with [id:%s]", id)
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": message})
		return
	}

	if input.Note != nil {
		assignment.Note = *input.Note
	}
	if input.Description != nil {
		assignment.Description = *input.Description
	}

	if input.Active != nil {
		assignment.Active = input.Active
	}

	if input.AllowRespondentActivate != nil {
		if !*requestingUser.IsAdmin {
			message := fmt.Sprintf("Invalid Request: You may not specify the \"AllowRespondentActivate\" option")
			c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": message})
			return
		}
		assignment.Active = input.Active
	}

	if err := assignmentService.Save(assignment); nil != err {
		c.JSON(http.StatusOK, gin.H{"status": "error", "message": err.Error()})
		return
	}

	message := fmt.Sprintf("Updated assignment \"%s\"", assignment.ID)
	c.JSON(http.StatusOK, gin.H{"data": assignment, "status": "success", "message": message})
}

func DeleteProductRespondentAssignment(c *gin.Context) {
	assignmentService := service.GetProductRespondentAssignmentService()

	id, err := uuid.Parse(c.Param("id"))
	if nil != err {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "Invalid Id provided"})
		return
	}
	assignment, err := assignmentService.GetById(id)
	if nil != err {
		message := fmt.Sprintf("Could not find assignment with [id:%s]", id)
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": message})
	}

	if err := assignmentService.Delete(assignment); nil != err {
		c.JSON(http.StatusOK, gin.H{"status": "error", "message": err.Error()})
		return
	}

	message := fmt.Sprintf("Deleted assignment \"%s\"", assignment.ID)
	c.JSON(http.StatusOK, gin.H{"data": assignment, "status": "success", "message": message})
}

func ValidateProductRespondentAssignmentCategory(category model.ProductCategory) (err error) {
	switch category {
	case model.PLACE_CITY:
		return nil
	case model.PLACE_STATE:
		return nil
	case model.PLACE_COUNTRY:
		return nil
	}
	return errors.New(fmt.Sprintf("Unsupported ProductRespondentAssignment Category \"%s\"", category))
}
