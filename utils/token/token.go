package token

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	jwt "github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/nkamuo/rasta-server/initializers"
	"github.com/nkamuo/rasta-server/service"
)

func GenerateToken(user_id uuid.UUID) (string, error) {

	appSecret := initializers.CONFIG.APP_SECRET
	TOKEN_LIFESPAN := os.Getenv("TOKEN_HOUR_LIFESPAN")
	if "" == TOKEN_LIFESPAN {
		TOKEN_LIFESPAN = "3600"
	}

	token_lifespan, err := strconv.Atoi(TOKEN_LIFESPAN)

	if err != nil {
		return "", err
	}

	claims := jwt.MapClaims{}
	claims["authorized"] = true
	claims["user_id"] = user_id
	claims["exp"] = time.Now().Add(time.Hour * time.Duration(token_lifespan)).Unix()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString([]byte(appSecret))
	// return token.SignedString([]byte(os.Getenv("API_SECRET")))
}

func TokenValid(context *gin.Context) error {
	tokenString := ExtractToken(context)
	_, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(os.Getenv("API_SECRET")), nil
	})

	userID, err := ExtractUserID(context)
	if nil != err {
		return err
	}

	// userID, err := uuid.Parse(user_id)

	// Proceed to get the claims
	// claims, ok := token.Claims.(jwt.MapClaims)
	// if !ok {
	// 	return errors.New("Invalid claims format")
	// }

	// // Access individual claims
	// userID, err := uuid.Parse(claims["user_id"].(string))
	// if nil != err {
	// 	return err
	// }
	user, err := service.GetUserService().GetById(userID)
	if nil != err {
		return err
	}
	context.Set("user", user) //MAKE THE CURRENT USER AVAILABLE THROUGHOUT THE REQUEST

	if err != nil {
		return err
	}
	return nil
}

func ExtractToken(c *gin.Context) string {
	token := c.Query("token")
	if token != "" {
		return token
	}
	bearerToken := c.Request.Header.Get("Authorization")
	if len(strings.Split(bearerToken, " ")) == 2 {
		return strings.Split(bearerToken, " ")[1]
	}
	return ""
}

func ExtractUserID(c *gin.Context) (uuid.UUID, error) {

	tokenString := ExtractToken(c)
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(os.Getenv("API_SECRET")), nil
	})
	if err != nil {
		return uuid.Nil, err
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if ok && token.Valid {
		// uid, err := strconv.ParseUint(fmt.Sprintf("%.0f", claims["user_id"]), 10, 32)
		user_id, ok := claims["user_id"].(string)
		if !ok {
			return uuid.Nil, errors.New("Could not extract user ID")
		}

		uid, err := uuid.Parse(user_id)
		if nil != err {
			return uuid.Nil, err
		}

		return uid, nil
	}
	return uuid.Nil, nil
}
