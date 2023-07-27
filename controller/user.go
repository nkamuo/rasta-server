package controller

import (
	"fmt"
	"net/http"

	"github.com/google/uuid"
	"github.com/nkamuo/rasta-server/dto"
	"github.com/nkamuo/rasta-server/model"
	"github.com/nkamuo/rasta-server/service"

	"github.com/gin-gonic/gin"
)

// GET /users
// Get all users
func FindUsers(c *gin.Context) {
	var users []model.User
	model.DB.Find(&users)

	c.JSON(http.StatusOK, gin.H{"data": users})
}

func CreateUser(c *gin.Context) {

	userService := service.GetUserService()
	// Validate input
	var input dto.UserCreationInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	// Create user
	user := model.User{
		Email: input.Email,
		// PERSONAL INFOR
		FirstName: input.FirstName,
		LastName:  input.LastName,
		Phone:     input.Phone,
	}

	userService.HashUserPassword(&user, input.Password)
	userService.Save(&user)

	c.JSON(http.StatusOK, gin.H{"data": user})
}

func FindUser(c *gin.Context) {
	userService := service.GetUserService()

	id, err := uuid.Parse(c.Param("id"))
	if nil != err {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "Invalid Id provided"})
		return
	}
	user, err := userService.GetById(id)
	if nil != err {
		message := fmt.Sprintf("Could not find user with [id:%s]", id)
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": message})
	}
	c.JSON(http.StatusOK, gin.H{"data": user})
}

func UpdateUser(c *gin.Context) {
	userService := service.GetUserService()

	// Validate input
	var input dto.UserUpdateInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	id, err := uuid.Parse(c.Param("id"))
	if nil != err {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "Invalid Id provided"})
		return
	}
	user, err := userService.GetById(id)
	if nil != err {
		message := fmt.Sprintf("Could not find user with [id:%s]", id)
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": message})
	}

	user.FirstName = input.FirstName
	user.LastName = input.LastName

	if err := userService.Save(user); nil != err {
		c.JSON(http.StatusOK, gin.H{"status": "error", "message": err.Error()})
		return
	}

	message := fmt.Sprintf("Updated user \"%s\"", user.FullName())
	c.JSON(http.StatusOK, gin.H{"data": user, "status": "success", "message": message})
}

func DeleteUser(c *gin.Context) {
	userService := service.GetUserService()

	id, err := uuid.Parse(c.Param("id"))
	if nil != err {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "Invalid Id provided"})
		return
	}
	user, err := userService.GetById(id)
	if nil != err {
		message := fmt.Sprintf("Could not find user with [id:%s]", id)
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": message})
	}

	if err := userService.Delete(user); nil != err {
		c.JSON(http.StatusOK, gin.H{"status": "error", "message": err.Error()})
		return
	}

	message := fmt.Sprintf("Deleted user \"%s\"", user.FullName())
	c.JSON(http.StatusOK, gin.H{"data": user, "status": "success", "message": message})
}
