package controller

import (
	"fmt"
	"log"
	"net/http"

	"github.com/google/uuid"
	"github.com/nkamuo/rasta-server/data/pagination"
	"github.com/nkamuo/rasta-server/dto"
	"github.com/nkamuo/rasta-server/model"
	"github.com/nkamuo/rasta-server/service"
	"github.com/nkamuo/rasta-server/utils/auth"
	"github.com/stripe/stripe-go/v74"
	"github.com/stripe/stripe-go/v74/customer"
	"github.com/stripe/stripe-go/v74/paymentmethod"

	"github.com/gin-gonic/gin"
)

func FindPaymentMethods(c *gin.Context) {
	var paymentMethodes []model.PaymentMethod
	var page pagination.Page
	if err := c.ShouldBindQuery(&page); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": err.Error()})
		return
	}

	if err := model.DB.Scopes(pagination.Paginate(paymentMethodes, &page, model.DB)).Find(&paymentMethodes).Error; nil != err {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": err.Error()})
		return
	}
	page.Rows = paymentMethodes
	c.JSON(http.StatusOK, gin.H{"data": page})
}

func FindUserPaymentMethods(c *gin.Context) {
	userService := service.GetUserService()

	rUser, err := auth.GetCurrentUser(c)
	if err != nil {
		message := fmt.Sprintf("You may not access this resource ")
		c.JSON(http.StatusUnauthorized, gin.H{"status": "error", "message": message})
		return
	}

	var user *model.User
	if c.Param("id") != "" {
		userID, err := uuid.Parse(c.Param("id"))
		if nil != err {
			c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "Invalid user id provided"})
			return
		}
		if user, err = userService.GetById(userID); nil != err {
			c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": err.Error()})
			return
		}

		if user.ID.String() != rUser.ID.String() && !*rUser.IsAdmin {
			message := fmt.Sprintf("You may not access this resource ")
			c.JSON(http.StatusForbidden, gin.H{"status": "error", "message": message})
			return
		}
	} else {
		user = rUser
	}

	// HERE >>
	var page pagination.Page
	if err := c.ShouldBindQuery(&page); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": err.Error()})
		return
	}

	// params := &stripe.PriceSearchParams{
	// 	// Limit: 10, // Adjust as needed
	// 	Product: stripe.String(config.STRIPE_RESPONDENT_SUBSCRIPTION_PRODUCT_ID),
	// 	Active:  stripe.Bool(true),
	// }

	// var Page *int64
	// if page.Page != 0 {
	// 	Page = stripe.Int64(int64(page.Page))
	// }
	// var Limit *int64
	// if page.Limit != nil {
	// 	Limit = stripe.Int64(int64(*page.Limit))
	// }

	// Query := fmt.Sprintf("active:'true' AND product: \"%s\"", "")
	params := &stripe.PaymentMethodListParams{
		Customer: user.StripeCustomerID,
	}

	iter := paymentmethod.List(params)
	var paymentMethods []*stripe.PaymentMethod

	for iter.Next() {
		paymentMethods = append(paymentMethods, iter.PaymentMethod())
	}

	if iter.Err() != nil {
		message := fmt.Sprintf("Error fetching prices: %s", iter.Err().Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": message})
		return
	}

	page.Rows = paymentMethods

	c.JSON(http.StatusOK, gin.H{"staus": "success", "data": page})
}

func CreatePaymentMethod(c *gin.Context) {

	userService := service.GetUserService()

	var input dto.PaymentMethodCreationInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	rUser, err := auth.GetCurrentUser(c)
	if err != nil {
		message := fmt.Sprintf("You may not access this resource ")
		c.JSON(http.StatusUnauthorized, gin.H{"status": "error", "message": message})
		return
	}

	var user *model.User
	if c.Param("id") != "" {
		userID, err := uuid.Parse(c.Param("id"))
		if nil != err {
			c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "Invalid user id provided"})
			return
		}
		if user, err = userService.GetById(userID); nil != err {
			c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": err.Error()})
			return
		}

		if user.ID.String() != rUser.ID.String() && !*rUser.IsAdmin {
			message := fmt.Sprintf("You may not access this resource ")
			c.JSON(http.StatusForbidden, gin.H{"status": "error", "message": message})
			return
		}
	} else {
		user = rUser
	}

	// Create a payment method for the customer
	paymentMethodParams := &stripe.PaymentMethodParams{
		Type: stripe.String(string(stripe.PaymentMethodTypeCard)),
		Card: &stripe.PaymentMethodCardParams{
			Token: stripe.String(input.Token), // Replace with an actual card token or source ID
		},
	}

	paymentMethod, err := paymentmethod.New(paymentMethodParams)
	if err != nil {
		message := fmt.Sprintf("Could not create payment method: %s", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": message})
		return
	}

	attachParams := &stripe.PaymentMethodAttachParams{
		Customer: stripe.String(*user.StripeCustomerID),
	}

	attachedPaymentMethod, err := paymentmethod.Attach(paymentMethod.ID, attachParams)
	if err != nil {
		message := fmt.Sprintf("Could not attach payment method to customer: %s", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": message})
		return
	}

	// Set the attached payment method as the default for the customer
	params := &stripe.CustomerParams{
		InvoiceSettings: &stripe.CustomerInvoiceSettingsParams{
			DefaultPaymentMethod: stripe.String(attachedPaymentMethod.ID),
		},
	}
	_, err = customer.Update(*user.StripeCustomerID, params)
	if err != nil {
		message := fmt.Sprintf("Could not set default payment method: %s", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": message})
		return
	}

	log.Default().Printf("Attached Payment method: %s", attachedPaymentMethod.Type)

	// attachedPaymentMethod.

	c.JSON(http.StatusOK, gin.H{"data": paymentMethod, "status": "success"})
}

func FindPaymentMethod(c *gin.Context) {
	paymentMethodeService := service.GetPaymentMethodService()

	id, err := uuid.Parse(c.Param("id"))
	if nil != err {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "Invalid Id provided"})
		return
	}

	paymentMethode, err := paymentMethodeService.GetById(id)
	if err != nil {
		message := fmt.Sprintf("Could not find paymentMethode with [id:%s]", id)
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": message})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": paymentMethode})
}

func SelectDefaultPaymentMethod(c *gin.Context) {
	userService := service.GetUserService()

	var input dto.SelectDefaultPaymentMethodInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	rUser, err := auth.GetCurrentUser(c)
	if err != nil {
		message := fmt.Sprintf("You may not access this resource ")
		c.JSON(http.StatusUnauthorized, gin.H{"status": "error", "message": message})
		return
	}

	var user *model.User
	if c.Param("id") != "" {
		userID, err := uuid.Parse(c.Param("id"))
		if nil != err {
			c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "Invalid user id provided"})
			return
		}
		if user, err = userService.GetById(userID); nil != err {
			c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": err.Error()})
			return
		}

		if user.ID.String() != rUser.ID.String() && !*rUser.IsAdmin {
			message := fmt.Sprintf("You may not access this resource ")
			c.JSON(http.StatusForbidden, gin.H{"status": "error", "message": message})
			return
		}
	} else {
		user = rUser
	}

	if _, err := userService.UpdateStripeCustomer(rUser, true); err != nil {
		message := fmt.Sprintf("Could not update user stripe customer: %s", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": message})
		return
	}

	pmParams := &stripe.PaymentMethodParams{Customer: stripe.String(*rUser.StripeCustomerID)}

	paymentmethod.Get(input.PaymentMethodID, pmParams)

	// Set the attached payment method as the default for the customer
	params := &stripe.CustomerParams{
		InvoiceSettings: &stripe.CustomerInvoiceSettingsParams{
			DefaultPaymentMethod: stripe.String(input.PaymentMethodID),
		},
	}

	_, err = customer.Update(*user.StripeCustomerID, params)
	if err != nil {
		message := fmt.Sprintf("Could not set default payment method: %s", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": message})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "success"})

}

func UpdatePaymentMethod(c *gin.Context) {
	paymentMethodeService := service.GetPaymentMethodService()

	var input dto.PaymentMethodUpdateInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	id, err := uuid.Parse(c.Param("id"))
	if nil != err {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "Invalid Id provided"})
		return
	}
	paymentMethode, err := paymentMethodeService.GetById(id)
	if nil != err {
		message := fmt.Sprintf("Could not find paymentMethode with [id:%s]", id)
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": message})
	}

	if nil != input.Category {
		paymentMethode.Category = *input.Category
		if err := model.ValidatePaymentMethodCategory(*input.Category); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": err.Error(), "status": "error"})
			return
		}
	}
	if nil != input.Details {
		paymentMethode.Details = *input.Details
	}
	if nil != input.Active {
		paymentMethode.Active = *input.Active
	}
	if nil != input.Description {
		paymentMethode.Description = *input.Description
	}

	if err := paymentMethodeService.Save(paymentMethode); nil != err {
		c.JSON(http.StatusOK, gin.H{"status": "error", "message": err.Error()})
		return
	}

	message := fmt.Sprintf("Updated payment method \"%s\"", paymentMethode.ID)
	c.JSON(http.StatusOK, gin.H{"data": paymentMethode, "status": "success", "message": message})
}

func DeletePaymentMethod(c *gin.Context) {
	id := c.Param("id")

	var paymentMethode model.PaymentMethod

	if err := model.DB.Where("id = ?", id).First(&paymentMethode).Error; err != nil {
		message := fmt.Sprintf("Could not find paymentMethode with [id:%s]", id)
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": message})
		return
	}
	model.DB.Delete(&paymentMethode)
	message := fmt.Sprintf("Deleted paymentMethode \"%s\"", paymentMethode.ID)
	c.JSON(http.StatusOK, gin.H{"data": paymentMethode, "status": "success", "message": message})
}
