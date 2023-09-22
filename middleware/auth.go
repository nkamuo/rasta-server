package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/nkamuo/rasta-server/utils/token"
)

func JwtAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		err := token.TokenValid(c)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"status": "error", "message": "Unauthorized"})
			c.Abort()
			return
		}
		c.Next()
	}
}
