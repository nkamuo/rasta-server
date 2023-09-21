package controller

import (
	// "fmt"
	"fmt"
	"net/http"

	// "os"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	// "github.com/nkamuo/rasta-server/model"
	"github.com/nkamuo/rasta-server/model"
	"github.com/nkamuo/rasta-server/service"

	// "github.com/google/uuid"
	// "github.com/nkamuo/rasta-server/service"
	// "github.com/stripe/stripe-go/payout"
	"github.com/stripe/stripe-go/v74"
	"github.com/stripe/stripe-go/v74/paymentmethod"
	"github.com/stripe/stripe-go/v74/payout"
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
	paymentService := service.GetOrderPaymentService()
	orderService := service.GetOrderService()

	id, err := uuid.Parse(c.Param("id"))
	if nil != err {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "Invalid Id provided"})
		return
	}

	order, err := orderService.GetById(id)
	if nil != err {
		message := fmt.Sprintf("Could not find product with [id:%s]", id)
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": message})
	}

	var payment *model.OrderPayment

	if order.PaymentID != nil {
		payment, err = paymentService.GetById(*order.PaymentID)
		if nil != err {
			message := fmt.Sprintf("Could not find order payment: %s", err.Error())
			c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": message})
			return
		}
	} else {
		payment, err = paymentService.InitOrderPayment(order)
		if nil != err {
			message := fmt.Sprintf("Could not init order payment: %s", err.Error())
			c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": message})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"data":   payment,
		"status": "success",
	})
}

// func CreateRespondentPayoutAccount(c *gin.Context) {

// 	var input dto.StripeExternalAccountInput;

// 	if err := c.ShouldBindJSON(&input); err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
// 		return
// 	}

// 	params := &stripe.AccountExternalAccountParams{
// 		// Object:  stripe.String("bank_account"),
// 		Country: stripe.String(input.Country),
// 		Currency: stripe.String(input.Currency),
// 		RoutingNumber: stripe.String(input.RoutingNumber),
// 		AccountNumber: stripe.String(input.AccountNumber),
// 		AccountHolderName: stripe.String(input.AccountHolderName),
// 		AccountHolderType: stripe.String(input.AccountHolderType),

// 	}

// 	payout, err := account.NewExternalAccount(params)
// 	if err != nil {
// 		message := fmt.Sprintf("Could not payout: %s", err.Error())
// 		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": message})
// 		return

// 	}

// 	c.JSON(http.StatusOK, gin.H{
// 		"data":   payout,
// 		"status": "success",
// 	})
// }

func AddPaymentMethod(c *gin.Context) {

	payoutParams := stripe.PayoutParams{
		Amount:      stripe.Int64(1000), // Amount in cents (e.g., $10.00)
		Currency:    stripe.String("usd"),
		Method:      stripe.String("instant"),
		Destination: stripe.String("your_customer_account_id"), // Customer's Stripe account ID
	}

	Payout, err := payout.New(&payoutParams)
	if err != nil {
		message := fmt.Sprintf("Could not init order payment: %s", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": message})
		return

	}

	c.JSON(http.StatusOK, gin.H{
		"data":   Payout,
		"status": "success",
	})
}
