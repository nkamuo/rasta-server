package auth

import (
	"github.com/gin-gonic/gin"
	"github.com/nkamuo/rasta-server/model"
	"github.com/nkamuo/rasta-server/repository"
	"github.com/nkamuo/rasta-server/service"
	"github.com/nkamuo/rasta-server/utils/token"
)

func GetCurrentUser(c *gin.Context) (user *model.User, err error) {
	userService := service.GetUserService()
	userID, err := token.ExtractUserID(c)
	if nil != err {
		return nil, err
	}
	user, err = userService.GetById(userID)
	if nil != err {
		return nil, err
	}
	return user, err
}

func GetCurrentRespondent(c *gin.Context, preload ...string) (user *model.Respondent, err error) {
	respondentRepo := repository.GetRespondentRepository()

	requestingUser, err := GetCurrentUser(c)
	if err != nil {
		return nil, err
	}

	respondant, err := respondentRepo.GetByUser(*requestingUser, preload...)
	if err != nil {
		return nil, err
	}

	return respondant, nil
}
