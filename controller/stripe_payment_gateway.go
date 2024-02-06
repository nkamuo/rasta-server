package controller

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/nkamuo/rasta-server/initializers"
	"github.com/nkamuo/rasta-server/service"
	"github.com/stripe/stripe-go/webhook"
)

func StripePaymentGatewayWebhook(c *gin.Context) {

	config, err := initializers.LoadConfig()
	if err != nil {
		message := fmt.Sprintf("Error fetching prices: %s", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": message})
		return
	}

	// Retrieve the raw request body
	payload, err := c.GetRawData()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	// Retrieve the webhook signature from the request headers
	sig := c.GetHeader("Stripe-Signature")

	// Verify the webhook signature
	event, err := webhook.ConstructEvent(payload, sig, config.STRIPE_WEBHOOK_SIGNING_SECRET_KEY)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid signature"})
		return
	}

	// Handle the event based on its type
	switch event.Type {
	case "checkout.session.completed":
		session := event.Data.Object
		sessionID := session["id"].(string)
		if err = processCheckoutSession(sessionID); err != nil {
			message := fmt.Sprintf("Error processing checkout session: %s", err.Error())
			c.JSON(http.StatusInternalServerError, gin.H{"error": message})
			return
		}

		// Handle the completed checkout session, update your database, etc.
		fmt.Println(fmt.Sprintf("%s:%s", event.Type, "\n\n\n"), session)
		break
	case "payment_intent.created":
		intent := event.Data.Object
		// Handle the successful payment intent, update your database, etc.
		fmt.Println(fmt.Sprintf("%s:%s", event.Type, "\n\n\n"), intent)
		break

	case "payment_intent.succeeded":
		intent := event.Data.Object
		// Handle the successful payment intent, update your database, etc.
		fmt.Println(fmt.Sprintf("%s:%s", event.Type, "\n\n\n"), intent)
		break
	// Add more cases for other event types if needed
	default:
		// Unexpected event type
		fmt.Println("Unhandled event type:", event.Type)
	}

	c.Status(http.StatusOK)
}

func processCheckoutSession(sessionID string) (err error) {
	// Fetch the checkout session from the Stripe API
	purchaseService := service.GetRespondentAccessProductPurchaseService()
	session, err := purchaseService.GetByStripeCheckoutID(sessionID)
	if err != nil {
		return nil
	}

	if err = purchaseService.Commit(session); err != nil {
		return err
	}

	return nil
	// session, err := session.Get(sessionID, nil)
	// if err != nil {
	// 	// Handle the error
	// 	return
	// }

	// Handle the completed checkout session, update your database, etc.
	// fmt.Println(session)
}
