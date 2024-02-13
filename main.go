package main

import (
	"fmt"
	"os"

	"github.com/nkamuo/rasta-server/command"
	"github.com/nkamuo/rasta-server/initializers"
	"github.com/nkamuo/rasta-server/model"
	"github.com/nkamuo/rasta-server/startup"
	"github.com/nkamuo/rasta-server/web"
	// "go-cli-for-git/cmd"
	// "net/http"
	// "github.com/joho/godotenv"
)

func main() {

	s3Bucket := os.Getenv("S3_BUCKET")
	secretKey := os.Getenv("SECRET_KEY")

	fmt.Println(s3Bucket, secretKey)
	// now do something with s3 or whatever

	config, err := initializers.LoadConfig()
	if err != nil {
		fmt.Println("CONFIG ERROR:", err)
	}
	model.ConnectDatabase(&config)

	if err := startup.Boot(); err != nil {
		fmt.Println("BOOT ERROR:", err)
	}

	// command.Execute()
	command.StartWebServer(web.WebServerConfig{
		Port: "8090",
	})
}

// func main() {
// 	// err := godotenv.Load()
// 	// if err != nil {
// 	// 	log.Fatal("Error loading .env file")
// 	// }

// 	s3Bucket := os.Getenv("S3_BUCKET")
// 	secretKey := os.Getenv("SECRET_KEY")

// 	fmt.Println(s3Bucket, secretKey)
// 	// now do something with s3 or whatever

// 	config, err := initializers.LoadConfig(".")

// 	if err != nil {
// 		fmt.Println("CONFIG ERROR:", err)
// 	}
// 	model.ConnectDatabase(&config)
// 	command.StartWebServer()
// }

// // eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdXRob3JpemVkIjp0cnVlLCJleHAiOjE3MDM0MjY4MjIsInVzZXJfaWQiOiI4MTJjYzc3NS00NzcyLTQ4NDEtYTA5My1iNjI0ZTQ4N2ZmMmMifQ.cB74Ta0crGVPEhrfwULTI-GiCVbc4jD2tuYFr2yDWTk
