package controller

import (
	"fmt"
	"net/http"
	"strings"

	//
	"github.com/nkamuo/rasta-server/auth"
	"github.com/nkamuo/rasta-server/dto"
	"github.com/nkamuo/rasta-server/model"
	"github.com/nkamuo/rasta-server/repository"
	"github.com/nkamuo/rasta-server/service"

	// "github.com/nkamuo/rasta-server/utils/token"
	"github.com/gin-gonic/gin"
	utils "github.com/nkamuo/rasta-server/utils/auth"
)

func Register(c *gin.Context) {

	userRepo := repository.GetUserRepository()
	userService := service.GetUserService()
	placeService := service.GetPlaceService()
	respondentService := service.GetRespondentService()

	var referrer *model.User

	var input dto.UserRegistrationInput

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": err.Error()})
		return
	}

	if input.ReferrerCode != nil {
		if _referrer, err := userRepo.GetByReferralCode(*input.ReferrerCode); nil != err {
			message := fmt.Sprintf("Could not resolve referrer from referralCode[%s]: %s", *input.ReferrerCode, err.Error())
			c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": message})
			return
		} else {
			referrer = _referrer
		}
	}

	var respondent *model.Respondent
	if input.IsRespondent {
		if nil == input.RespondentPlaceId {
			message := fmt.Sprintf("Can't create a respondent account without specifying a base place of operation")
			c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": message})
			return
		} else {
			place, err := placeService.GetById(*input.RespondentPlaceId)
			if err != nil {
				message := fmt.Sprintf("Could not resolve place with [id:%s]", *input.RespondentPlaceId)
				c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": message})
				return
			}

			respondent = &model.Respondent{
				PlaceID: &place.ID,
			}
		}
	}

	user := model.User{
		Email:     strings.TrimSpace(input.Email),
		FirstName: input.FirstName,
		LastName:  input.LastName,
		Phone:     input.Phone,
	}

	if nil != referrer {
		user.ReferrerID = &referrer.ID
	}

	// u.HashedPassword = input.Password

	userService.UpdateStripeCustomer(&user, false)
	userService.HashUserPassword(&user, input.Password)
	err := userService.Save(&user)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": err.Error()})
		return
	}

	if respondent != nil {
		respondent.UserID = &user.ID
		err := respondentService.Save(respondent)
		if err != nil {
			message := fmt.Sprintf("Error creating responder profile: %s", err.Error())
			c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": message})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{"status": "success", "message": "registration success"})

}

func Login(c *gin.Context) {

	userRepo := repository.GetUserRepository()
	// userService := service.GetUserService();

	var input dto.UserFormLoginInput

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": err.Error()})
		return
	}

	user, err := userRepo.GetByEmail(input.Username)
	if nil != err {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "username or password is incorrect."})
		return
	}

	token, err := auth.LoginCheck(*user, input.Password)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "success", "token": token})

}

func GetCurrentUser(c *gin.Context) {

	user, err := utils.GetCurrentUser(c)

	if nil != err {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "success", "data": user})

}

func DeleteCurrentUser(c *gin.Context) {

	userService := service.GetUserService()
	user, err := utils.GetCurrentUser(c)

	if nil != err {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": err.Error()})
		return
	}

	if err = userService.Delete(user); nil != err {
		message := fmt.Sprintf("Error deleting user: %s", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "data": user, "message": message})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "success", "data": user, "message": "user deleted successfully"})

}
