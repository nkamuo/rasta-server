package controller

import (
	"net/http"
	"strings"

	"github.com/nkamuo/rasta-server/auth"
	"github.com/nkamuo/rasta-server/dto"
	"github.com/nkamuo/rasta-server/model"
	"github.com/nkamuo/rasta-server/repository"
	"github.com/nkamuo/rasta-server/service"

	"github.com/gin-gonic/gin"
)

func Register(c *gin.Context) {

	userSerice := service.GetUserService()

	var input dto.UserRegistrationInput

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user := model.User{
		Email:     strings.TrimSpace(input.Email),
		FirstName: input.FirstName,
		LastName:  input.LastName,
		Phone:     input.Phone,
	}

	// u.HashedPassword = input.Password

	userSerice.HashUserPassword(&user, input.Password)
	err := userSerice.Save(&user)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "registration success"})

}

func Login(c *gin.Context) {

	userRepo := repository.GetUserRepository()
	// userService := service.GetUserService();

	var input dto.UserFormLoginInput

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := userRepo.GetByEmail(input.Username)
	if nil != err {
		c.JSON(http.StatusBadRequest, gin.H{"error": "username or password is incorrect.", "user": user})
		return
	}

	token, err := auth.LoginCheck(*user, input.Password)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "username or password is incorrect.", "message": err.Error(), "user": user})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": token})

}
