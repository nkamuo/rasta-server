package controller

import (
	// "fmt"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	// "github.com/google/uuid"
	// "github.com/nkamuo/rasta-server/service"
	"github.com/stripe/stripe-go/v74"
	"github.com/stripe/stripe-go/v74/paymentintent"
	"github.com/stripe/stripe-go/v74/paymentmethod"
)

func createPaymentMethod(c *gin.Context) {
	// Get the user's submitted card details from the request.
	// For simplicity, let's assume the card information is submitted as JSON.
	var cardDetails struct {
		CardNumber  string `json:"card_number"`
		ExpiryMonth int64  `json:"expiry_month"`
		ExpiryYear  int64  `json:"expiry_year"`
		CVC         string `json:"cvc"`
	}
	if err := c.BindJSON(&cardDetails); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request data"})
		return
	}

	// Create the PaymentMethod with the card details.
	pmParams := &stripe.PaymentMethodParams{
		Type: stripe.String("card"),
		Card: &stripe.PaymentMethodCardParams{
			Number:   stripe.String(cardDetails.CardNumber),
			ExpMonth: stripe.Int64(cardDetails.ExpiryMonth),
			ExpYear:  stripe.Int64(cardDetails.ExpiryYear),
			CVC:      stripe.String(cardDetails.CVC),
		},
	}
	paymentMethod, err := paymentmethod.New(pmParams)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error creating PaymentMethod"})
		return
	}

	// paymentMethod.Card.

	c.JSON(http.StatusOK, gin.H{"payment_method_id": paymentMethod.ID})
}

func CreatePaymentIntent(c *gin.Context) {

	stripe.Key = os.Getenv("STRIPE_SECRET_KEY")
	// productService := service.GetProductService()

	// id, err := uuid.Parse(c.Param("id"))
	// if nil != err {
	// 	c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "Invalid Id provided"})
	// 	return
	// }
	// _, err = productService.GetById(id)
	// if nil != err {
	// 	message := fmt.Sprintf("Could not find product with [id:%s]", id)
	// 	c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": message})
	// }

	// Create a PaymentIntent with amount and currency
	params := &stripe.PaymentIntentParams{
		Amount:   stripe.Int64(int64(409)),
		Currency: stripe.String(string(stripe.CurrencyUSD)),
		AutomaticPaymentMethods: &stripe.PaymentIntentAutomaticPaymentMethodsParams{
			Enabled: stripe.Bool(true),
		},
	}

	pi, err := paymentintent.New(params)
	log.Printf("pi.New: %v", pi.ClientSecret)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		log.Printf("pi.New: %v", err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"clientSecret": pi.ClientSecret,
	})
}
