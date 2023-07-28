package controller

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/google/uuid"
	// "github.com/mitchellh/mapstructure"
	"github.com/nkamuo/rasta-server/dto"
	"github.com/nkamuo/rasta-server/model"
	"github.com/nkamuo/rasta-server/service"
	"github.com/nkamuo/rasta-server/utils/auth"

	"github.com/gin-gonic/gin"
	// "github.com/gin-gonic/gin"
)

func FindProductRespondentAssignments(c *gin.Context) {
	var assignments []model.ProductRespondentAssignment

	if err := model.DB.Find(&assignments).Error; nil != err {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "success", "data": assignments})
}

func CreateProductRespondentAssignment(c *gin.Context) {

	productService := service.GetProductService()
	respondentService := service.GetRespondentService()
	assignmentService := service.GetProductRespondentAssignmentService()

	user, err := auth.GetCurrentUser(c)
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

	respondant, err := respondentService.GetById(input.RespondentID)
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
			"There is already an assignment for respondent \"%v\" on Product \"%v\"",
			respondant.User.FullName(),
			product.Label,
		)
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": message})
		return

	}

	// Create assignment
	assignment := model.ProductRespondentAssignment{
		ProductID:    product.ID,
		RespondentID: respondant.ID,
		Note:         input.Note,
		Description:  input.Description,
	}

	if user.IsAdmin {
		assignment.Active = input.Active
		assignment.AllowRespondentActivate = input.AllowRespondentActivate
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

	var requestBody map[string]interface{}
	if err := c.Copy().BindJSON(&requestBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": err.Error()})
		return
	}

	// var input dto.ProductRespondentAssignmentUpdateInput
	// mapstructure.Decode(requestBody, input)

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

	if note, ok := requestBody["note"]; ok {
		assignment.Note = note.(string)
	}
	if description, ok := requestBody["description"]; ok {
		assignment.Description = description.(string)
	}

	if active, ok := requestBody["active"]; ok {
		assignment.Active = active.(bool)
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
