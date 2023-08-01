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
	api.GET("/me", controller.GetCurrentUser)
	api.GET("/me/respondent", controller.GetCurrentRespondent)

	api.GET("/products", controller.FindProducts)
	api.GET("/products/:id", controller.FindProduct)
	api.POST("/products", controller.CreateProduct)
	api.PATCH("/products/:id", controller.UpdateProduct)
	api.DELETE("/products/:id", controller.DeleteProduct)

	api.GET("/users", controller.FindUsers)
	api.GET("/users/:id", controller.FindUser)
	api.POST("/users", controller.CreateUser)
	api.PATCH("/users/:id", controller.UpdateUser)
	api.DELETE("/users/:id", controller.DeleteUser)

	api.GET("/orders", controller.FindOrders)
	api.GET("/orders/:id", controller.FindOrder)
	api.POST("/orders", controller.CreateOrder)
	api.PATCH("/orders/:id", controller.UpdateOrder)
	api.DELETE("/orders/:id", controller.DeleteOrder)

	api.GET("/respondents", controller.FindRespondents)
	api.GET("/respondents/:id", controller.FindRespondent)
	api.POST("/respondents", controller.CreateRespondent)
	api.PATCH("/respondents/:id", controller.UpdateRespondent)
	api.DELETE("/respondents/:id", controller.DeleteRespondent)
	// COMPANY RESPONDANTS
	api.GET("/companies/:id/respondents", controller.FindRespondentsByCompany)
	api.POST("/companies/:id/respondents", controller.AddRespondentToCompany)
	api.DELETE("/companies/:id/respondents/:respondent_id", controller.RemoveRespondentFromCompany)

	api.GET("/companies", controller.FindCompanies)
	api.GET("/companies/:id", controller.FindCompany)
	api.POST("/companies", controller.CreateCompany)
	api.PATCH("/companies/:id", controller.UpdateCompany)
	api.DELETE("/companies/:id", controller.DeleteCompany)

	api.GET("/places", controller.FindPlaces)
	api.GET("/places/:id", controller.FindPlace)
	api.POST("/places", controller.CreatePlace)
	api.PATCH("/places/:id", controller.UpdatePlace)
	api.DELETE("/places/:id", controller.DeletePlace)

	api.GET("/product_respondent_assignments", controller.FindProductRespondentAssignments)
	api.GET("/product_respondent_assignments/:id", controller.FindProductRespondentAssignment)
	api.POST("/product_respondent_assignments", controller.CreateProductRespondentAssignment)
	api.PATCH("/product_respondent_assignments/:id", controller.UpdateProductRespondentAssignment)
	api.DELETE("/product_respondent_assignments/:id", controller.DeleteProductRespondentAssignment)

	api.GET("/fuel_types", controller.FindFuelTypes)
	api.GET("/fuel_types/:id", controller.FindFuelType)
	api.POST("/fuel_types", controller.CreateFuelType)
	api.PATCH("/fuel_types/:id", controller.UpdateFuelType)
	api.DELETE("/fuel_types/:id", controller.DeleteFuelType)

	api.GET("/payment_methods", controller.FindPaymentMethods)
	api.GET("/payment_methods/:id", controller.FindPaymentMethod)
	api.POST("/payment_methods", controller.CreatePaymentMethod)
	api.PATCH("/payment_methods/:id", controller.UpdatePaymentMethod)
	api.DELETE("/payment_methods/:id", controller.DeletePaymentMethod)

	// USER PAYMENT METHODS
	api.GET("/users/:id/payment_methods", controller.FindUserPaymentMethods)

	// product_respondent_assignment

	r.Run(":8090")
}

// eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdXRob3JpemVkIjp0cnVlLCJleHAiOjE3MDM0MjY4MjIsInVzZXJfaWQiOiI4MTJjYzc3NS00NzcyLTQ4NDEtYTA5My1iNjI0ZTQ4N2ZmMmMifQ.cB74Ta0crGVPEhrfwULTI-GiCVbc4jD2tuYFr2yDWTk
