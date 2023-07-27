package main

import (
	"fmt"
	"log"
	"os"

	"github.com/nkamuo/rasta-server/controller"
	"github.com/nkamuo/rasta-server/initializers"
	"github.com/nkamuo/rasta-server/middleware"
	"github.com/nkamuo/rasta-server/model"

	// "net/http"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	s3Bucket := os.Getenv("S3_BUCKET")
	secretKey := os.Getenv("SECRET_KEY")

	fmt.Println(s3Bucket, secretKey)
	// now do something with s3 or whatever

	r := gin.Default()

	config, err := initializers.LoadConfig(".")

	if err != nil {
		fmt.Println("CONFIG ERROR:", err)
	}
	model.ConnectDatabase(&config)

	api := r.Group("/api")
	api.POST("/register", controller.Register)
	api.POST("/login", controller.Login)
	api.Use(middleware.JwtAuthMiddleware())

	api.GET("/products", controller.FindProducts)
	api.GET("/products/:id", controller.FindProduct)
	api.POST("/products", controller.CreateProduct)
	api.DELETE("/products/:id", controller.DeleteProduct)

	r.Run(":8090")
}

// eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdXRob3JpemVkIjp0cnVlLCJleHAiOjE3MDM0MjY4MjIsInVzZXJfaWQiOiI4MTJjYzc3NS00NzcyLTQ4NDEtYTA5My1iNjI0ZTQ4N2ZmMmMifQ.cB74Ta0crGVPEhrfwULTI-GiCVbc4jD2tuYFr2yDWTk
