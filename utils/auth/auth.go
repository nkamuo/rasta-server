package auth

import (
	"github.com/gin-gonic/gin"
	"github.com/nkamuo/rasta-server/model"
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
