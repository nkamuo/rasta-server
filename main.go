package main

import (
	"fmt"
	"os"

	"github.com/nkamuo/rasta-server/controller"
	"github.com/nkamuo/rasta-server/initializers"
	"github.com/nkamuo/rasta-server/middleware"
	"github.com/nkamuo/rasta-server/model"

	// "net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	// "github.com/joho/godotenv"
)

func main() {
	// err := godotenv.Load()
	// if err != nil {
	// 	log.Fatal("Error loading .env file")
	// }

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

	corsConfig := cors.DefaultConfig()
	corsConfig.AllowOrigins = []string{"http://localhost:8080", "http://localhost:3000"}
	corsConfig.AddAllowHeaders("Authorization")
	// config.AllowOrigins = []string{"http://google.com", "http://facebook.com"}
	// config.AllowAllOrigins = true

	r.Use(cors.New(corsConfig))

	api := r.Group("/api")

	// PUBLIC ENDPOINTS
	api.POST("/register", controller.Register)
	api.POST("/login", controller.Login)

	api.GET("/motorist_request_situations", controller.FindMotoristRequestSituations)
	api.GET("/motorist_request_situations/:id", controller.FindMotoristRequestSituation)
	api.POST("/motorist_request_situations", controller.CreateMotoristRequestSituation)
	api.PATCH("/motorist_request_situations/:id", controller.UpdateMotoristRequestSituation)
	api.DELETE("/motorist_request_situations/:id", controller.DeleteMotoristRequestSituation)

	// AUTHENTICATION MIDDLEWARE
	api.Use(middleware.JwtAuthMiddleware())
	//AUTHENICATED ENDPOINTS
	api.GET("/me", controller.GetCurrentUser)
	api.GET("/me/respondent", controller.GetCurrentRespondent)
	api.GET("/me/respondent/session", controller.FindCurrentRespondentSession)
	api.GET("/me/respondent/session/requests", controller.FindAvailableOrdersForRespondent)
	api.GET("/me/respondent/session/requests/:id", controller.FindOrderForRespondent)
	api.POST("/me/respondent/session/requests/:id/claim", controller.RespondentClaimOrder)

	api.GET("/respondent_sessions", controller.FindRespondentSessions)
	api.GET("/respondent_sessions/:id", controller.FindRespondentSession)
	api.POST("/respondent_sessions", controller.CreateRespondentSession)
	api.PATCH("/respondent_sessions/:id", controller.UpdateRespondentSession)
	api.DELETE("/respondent_sessions/:id", controller.DeleteRespondentSession)

	//PROTECTED ENDPOINTS
	api.GET("/products", controller.FindProducts)
	api.GET("/products/find_by_category_and_location", controller.FindProductByCategoryAndLocation)
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
	//
	api.GET("/orders/:id/payment", controller.CreatePaymentIntent)
	//REQUESTS -> ORDER-ITEMS
	api.GET("/requests", controller.FindRequests)

	api.GET("/respondents", controller.FindRespondents)
	api.GET("/respondents/:id", controller.FindRespondent)
	api.POST("/respondents", controller.CreateRespondent)
	api.PATCH("/respondents/:id", controller.UpdateRespondent)
	api.DELETE("/respondents/:id", controller.DeleteRespondent)

	//RESPONDENT REVIEWS
	api.GET("/respondent_service_reviews", controller.FindRespondentServiceReviews)
	api.GET("/respondent_service_reviews/:id", controller.FindRespondentServiceReview)
	api.POST("/respondent_service_reviews", controller.CreateRespondentServiceReview)
	api.PATCH("/respondent_service_reviews/:id", controller.UpdateRespondentServiceReview)
	api.DELETE("/respondent_service_reviews/:id", controller.DeleteRespondentServiceReview)

	// COMPANY RESPONDANTS
	api.GET("/companies/:id/respondents", controller.FindRespondentsByCompany)
	api.POST("/companies/:id/respondents", controller.AddRespondentToCompany)
	api.DELETE("/companies/:id/respondents/:respondent_id", controller.RemoveRespondentFromCompany)

	api.GET("/companies", controller.FindCompanies)
	api.GET("/companies/:id", controller.FindCompany)
	api.POST("/companies", controller.CreateCompany)
	api.PATCH("/companies/:id", controller.UpdateCompany)
	api.DELETE("/companies/:id", controller.DeleteCompany)

	api.GET("/vehicles", controller.FindVehicles)
	api.GET("/vehicles/:id", controller.FindVehicle)
	api.POST("/vehicles", controller.CreateVehicle)
	api.PATCH("/vehicles/:id", controller.UpdateVehicle)
	api.DELETE("/vehicles/:id", controller.DeleteVehicle)

	api.GET("/vehicle_models", controller.FindVehicleModels)
	api.GET("/vehicle_models/:id", controller.FindVehicleModel)
	api.POST("/vehicle_models", controller.CreateVehicleModel)
	api.PATCH("/vehicle_models/:id", controller.UpdateVehicleModel)
	api.DELETE("/vehicle_models/:id", controller.DeleteVehicleModel)

	api.GET("/places", controller.FindPlaces)
	api.GET("/places/:id", controller.FindPlace)
	api.GET("/places/find-by-location", controller.FindPlaceByLocation)
	api.POST("/places", controller.CreatePlace)
	api.PATCH("/places/:id", controller.UpdatePlace)
	api.DELETE("/places/:id", controller.DeletePlace)

	//LOCATION
	api.GET("/locations", controller.FindLocations)
	api.GET("/locations/:id", controller.FindLocation)
	api.GET("/locations/distance", controller.ResolveDistanceMatrix)
	api.GET("/locations/resolve", controller.ResolveLocation)
	api.GET("/locations/resolve_by_ip", controller.ResolveLocationByIPAddress)
	api.DELETE("/locations/:id", controller.DeleteLocation)

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

	api.GET("/fuel_type_place_rates", controller.FindFuelTypePlaceRates)
	api.GET("/fuel_type_place_rates/:id", controller.FindFuelTypePlaceRate)
	api.POST("/fuel_type_place_rates", controller.CreateFuelTypePlaceRate)
	api.PATCH("/fuel_type_place_rates/find_by_type_and_location", controller.FindFuelTypePlaceRateByTypeAndLocation)
	api.PATCH("/fuel_type_place_rates/:id", controller.UpdateFuelTypePlaceRate)
	api.DELETE("/fuel_type_place_rates/:id", controller.DeleteFuelTypePlaceRate)

	api.GET("/towing_place_rates", controller.FindTowingPlaceRates)
	api.GET("/towing_place_rates/:id", controller.FindTowingPlaceRate)
	api.GET("/towing_place_rates/find_by_place_and_distance", controller.FindTowingRateByPlaceAndDistance)
	api.GET("/towing_place_rates/find_by_origin_and_destination", controller.FindTowingRateByOriginAndDestination)
	api.POST("/towing_place_rates", controller.CreateTowingPlaceRate)
	api.PATCH("/towing_place_rates/:id", controller.UpdateTowingPlaceRate)
	api.DELETE("/towing_place_rates/:id", controller.DeleteTowingPlaceRate)

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
