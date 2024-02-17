package middleware

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/nkamuo/rasta-server/service"
	"github.com/nkamuo/rasta-server/utils/auth"
)

func CanHandleMotoristRequestMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		respondentService := service.GetRespondentService()
		responent, err := auth.GetCurrentRespondent(c)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"status": "error", "message": "Unauthorized"})
			c.Abort()
			return
		}

		canHandle, err := respondentService.CanHandleMotoristRequest(responent)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": err.Error()})
			c.Abort()
			return
		}

		if !canHandle {
			message := fmt.Sprintf("You are not allowed to handle this request: %s", err.Error())
			c.JSON(http.StatusForbidden, gin.H{"status": "error", "message": message})
			c.Abort()
			return
		}
		c.Next()
	}
}
