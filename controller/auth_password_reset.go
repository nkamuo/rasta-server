package controller

import (
	"fmt"
	"net/http"
	"time"

	// "github.com/d-vignesh/go-jwt-auth/utils"
	"github.com/gin-gonic/gin"
	"github.com/nkamuo/rasta-server/dto"
	"github.com/nkamuo/rasta-server/model"
	"github.com/nkamuo/rasta-server/repository"
	"github.com/nkamuo/rasta-server/service"
	mailing "github.com/nkamuo/rasta-server/utils"
)

func AuthPasswordResetGenerateCode(c *gin.Context) {

	userVerificationRequestService := service.GetUserVerificationRequestService()
	userRepo := repository.GetUserRepository()

	var input dto.UserFormResetPasswordCodeRequestInput

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": err.Error()})
		return
	}

	user, err := userRepo.GetByEmail(input.Email)
	if nil != err {
		message := fmt.Sprintf("Could not resolve user with email[%s]: %s", input.Email, err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": message})
		return
	}
	Now := time.Now()
	ExpiresAt := Now.Add(time.Minute * 15)
	Code := mailing.GenerateRandomNumbers(6)
	Token := mailing.GenerateRandomString(32)

	userVerificationRequest := &model.UserVerificationRequest{
		Email:       user.Email,
		Code:        Code,
		Token:       Token,
		ExpiresAt:   ExpiresAt,
		RequestType: "PASSWORD_RESET",
	}

	// mailReq := &mail.Mail{
	// 	from:    "",
	// 	to:      []string{user.Email},
	// 	subject: "Password Reset",
	// 	mtype:   mail.PassReset,
	// 	data: &mail.MailData{
	// 		Username: user.Username,
	// 		Code:     code,
	// 	},
	// }

	Name := user.FullName()

	err = mailing.SendPasswordResetEmail(Name, user.Email, Code)
	if err != nil {
		message := fmt.Sprintf("Could not send password reset code to user with email[%s]: %s", input.Email, err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": message})
		return
	}

	err = userVerificationRequestService.Save(userVerificationRequest)
	if err != nil {
		message := fmt.Sprintf("Error creating password reset request: %s", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": message})
		return
	}

	data := map[string]interface{}{
		"expiresAt": ExpiresAt,
		"sentAt":    Now,
	}

	c.JSON(http.StatusOK, gin.H{"status": "success", "data": data, "message": "code sent successfully"})
}

func AuthPasswordResetVerifyCode(c *gin.Context) {
	userVerificationRequestService := service.GetUserVerificationRequestService()
	userRepo := repository.GetUserRepository()

	var input dto.UserFormResetPasswordCodeValidationInput

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": err.Error()})
		return
	}

	_, err := userRepo.GetByEmail(input.Email)
	if nil != err {
		message := fmt.Sprintf("Could not resolve user with email[%s]: %s", input.Email, err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": message})
		return
	}

	userVerificationRequest, err := userVerificationRequestService.GetByCodeAndEmail(input.Code, input.Email)
	if nil != err {
		message := fmt.Sprintf("Invalid code provided")
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": message})
		return
	}

	Now := time.Now()

	if userVerificationRequest.ExpiresAt.Before(Now) {
		message := fmt.Sprintf("Code has expired")
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": message})
		return
	}

	data := map[string]interface{}{
		"Token": userVerificationRequest.Token,
	}

	c.JSON(http.StatusOK, gin.H{"status": "success", "data": data, "message": "code verified successfully"})
}

func AuthPasswordResetCommit(c *gin.Context) {

	userVerificationRequestService := service.GetUserVerificationRequestService()
	userService := service.GetUserService()
	userRepo := repository.GetUserRepository()

	var input dto.UserFormResetPasswordInput

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": err.Error()})
		return
	}

	_, err := userRepo.GetByEmail(input.Email)
	if nil != err {
		message := fmt.Sprintf("Could not resolve user with email[%s]: %s", input.Email, err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": message})
		return
	}

	userVerificationRequest, err := userVerificationRequestService.GetByCodeAndEmail(input.Code, input.Email)
	if nil != err {
		message := fmt.Sprintf("Invalid code provided")
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": message})
		return
	}

	Now := time.Now()

	if userVerificationRequest.ExpiresAt.Before(Now) {
		message := fmt.Sprintf("Code has expired")
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": message})
		return
	}

	if userVerificationRequest.Token != input.ResetToken {
		message := fmt.Sprintf("Invalid reset token provided")
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": message})
		return
	}

	user, err := userRepo.GetByEmail(input.Email)

	if nil != err {
		message := fmt.Sprintf("Could not resolve user with email[%s]: %s", input.Email, err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": message})
		return
	}

	// >> RESET PASSWORD
	userService.HashUserPassword(user, input.Password)
	err = userService.Save(user)
	// << RESET PASSWORD

	if nil != err {
		message := fmt.Sprintf("Error saving user: %s", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": message})
		return
	}

	err = userVerificationRequestService.Delete(userVerificationRequest)
	if nil != err {
		message := fmt.Sprintf("Error deleting user verification request: %s", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": message})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "success", "message": "password reset success"})
}
