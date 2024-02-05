package controller

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/nkamuo/rasta-server/data/pagination"
	"github.com/nkamuo/rasta-server/dto"
	"github.com/nkamuo/rasta-server/initializers"
	"github.com/nkamuo/rasta-server/utils/auth"
	"github.com/stripe/stripe-go/v74"
	"github.com/stripe/stripe-go/v74/price"

	// "github.com/stripe/stripe-go/v74/subscription"
	"github.com/stripe/stripe-go/v74/checkout/session"
)

func FindRespondentPurchasePrices(c *gin.Context) {
	config, err := initializers.LoadConfig()
	if err != nil {
		message := fmt.Sprintf("Error fetching prices: %s", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": message})
		return
	}
	var page pagination.Page
	if err := c.ShouldBindQuery(&page); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": err.Error()})
		return
	}

	var Limit *int64
	if page.Limit != nil {
		Limit = stripe.Int64(int64(*page.Limit))
	}

	Query := fmt.Sprintf("active:'true' AND product: \"%s\"", config.STRIPE_RESPONDENT_PURCHASE_PRODUCT_ID)
	params := &stripe.PriceSearchParams{
		SearchParams: stripe.SearchParams{
			Query: Query,
			Limit: Limit,
			// Page:  Page,
			//   Query: "active:'true' AND metadata['order_id']:'6735'",
		},
	}

	iter := price.Search(params)
	var prices []*stripe.Price

	for iter.Next() {
		prices = append(prices, iter.Price())
	}

	if iter.Err() != nil {
		message := fmt.Sprintf("Error fetching prices: %s", iter.Err().Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": message})
		return
	}

	var TotalRows int64 = 0
	var TotalPages int64 = 0
	if iter.PriceSearchResult().TotalCount != nil {
		TotalRows = int64(*iter.PriceSearchResult().TotalCount)
	} else {
		TotalRows = int64(len(prices))
	}
	if Limit != nil {
		TotalPages = TotalRows / (*Limit)
	} else {
		TotalPages = TotalRows / (10)
	}
	page.TotalPages = int(TotalPages)
	page.TotalRows = TotalRows
	page.Rows = prices

	c.JSON(http.StatusOK, gin.H{"staus": "success", "data": page})
}

func CreateRespondentPurchaseCheckoutSession(c *gin.Context) {
	var input dto.RespondentPurchaseSessionCheckoutInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": err.Error()})
		return
	}

	rUser, err := auth.GetCurrentUser(c)
	if err != nil {
		message := fmt.Sprintf("You may not access this resource ")
		c.JSON(http.StatusUnauthorized, gin.H{"status": "error", "message": message})
		return
	}

	pParams := &stripe.PriceParams{}

	rPrice, err := price.Get(input.PriceID, pParams)
	if err != nil {
		message := fmt.Sprintf("Error fetching prices: %s", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": message})
		return
	}

	mode := stripe.CheckoutSessionModePayment
	if rPrice.Recurring != nil {
		mode = stripe.CheckoutSessionModePayment
	}

	var Quantity int64 = 1
	if input.Quantity != nil {
		Quantity = *input.Quantity
	}
	params := &stripe.CheckoutSessionParams{
		// Customer: rUser.StripeCustomerID,
		// PaymentMethodTypes: stripe.StringSlice([]string{"card"}),
		LineItems: []*stripe.CheckoutSessionLineItemParams{
			{
				Price:    stripe.String(input.PriceID), // Replace with your actual Price ID
				Quantity: stripe.Int64(Quantity),
			},
		},
		Mode:     stripe.String(string(mode)),
		Customer: rUser.StripeCustomerID,

		// SuccessURL: stripe.String("https://yourwebsite.com/success"),
		// CancelURL:  stripe.String("https://yourwebsite.com/cancel"),
	}

	session, err := session.New(params)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"session": map[string]interface{}{
			"id":          session.ID,
			"paymentLink": session.PaymentLink,
		},
	})
}
