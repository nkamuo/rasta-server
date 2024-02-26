package utils

import (
	"fmt"

	"github.com/nkamuo/rasta-server/initializers"
	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

func SendPasswordResetEmail(
	name string,
	email string,
	code string,
) (err error) {

	config, err := initializers.LoadConfig()
	if err != nil {
		return err
	}

	ApiKey := config.SENDGRID_API_KEY
	FromName := config.SENDGRID_FROM_NAME
	FromEmail := config.SENDGRID_FROM_EMAIL

	from := mail.NewEmail(FromName, FromEmail)
	subject := "Huqt Password Reset"
	to := mail.NewEmail(name, email)
	plainTextContent := fmt.Sprintf("Hello %s!. Your password reset code is %s", name, code)
	htmlContent := fmt.Sprintf("<strong>Hello %s!</strong><br>Your password reset code is <strong>%s</strong>", name, code)
	message := mail.NewSingleEmail(from, subject, to, plainTextContent, htmlContent)
	client := sendgrid.NewSendClient(ApiKey)
	response, err := client.Send(message)
	fmt.Println(response.StatusCode)
	return err
	// if err != nil {
	// 	return err
	// } else {
	// 	fmt.Println(response.StatusCode)
	// 	fmt.Println(response.Body)
	// 	fmt.Println(response.Headers)
	// }
}
