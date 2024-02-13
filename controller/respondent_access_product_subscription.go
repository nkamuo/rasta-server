package controller

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/google/uuid"
	"github.com/stripe/stripe-go/v74"
	"github.com/stripe/stripe-go/v74/subscription"

	// "github.com/mitchellh/mapstructure"
	"github.com/nkamuo/rasta-server/data/pagination"
	"github.com/nkamuo/rasta-server/initializers"
	"github.com/nkamuo/rasta-server/model"
	"github.com/nkamuo/rasta-server/repository"
	"github.com/nkamuo/rasta-server/service"
	"github.com/nkamuo/rasta-server/utils/auth"

	"github.com/gin-gonic/gin"
	// "github.com/gin-gonic/gin"
)

func FindUserSubscriptions(c *gin.Context) {
	userService := service.GetUserService()

	config, err := initializers.LoadConfig()
	if err != nil {
		message := fmt.Sprintf("You may not access this resource ")
		c.JSON(http.StatusUnauthorized, gin.H{"status": "error", "message": message})
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
	params := &stripe.SubscriptionListParams{
		Customer: user.StripeCustomerID,
	}

	iter := subscription.List(params)
	var subscriptions []*stripe.Subscription

	for iter.Next() {
		sub := iter.Subscription()
		for _, item := range sub.Items.Data {
			if item.Price.Product.ID == config.STRIPE_RESPONDENT_SUBSCRIPTION_PRODUCT_ID {
				subscriptions = append(subscriptions, sub)
				break
			}
			// fmt.Println(item.Price.ID)
		}
	}

	if iter.Err() != nil {
		message := fmt.Sprintf("Error fetching prices: %s", iter.Err().Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": message})
		return
	}

	page.Rows = subscriptions

	c.JSON(http.StatusOK, gin.H{"staus": "success", "data": page})
}

// ////////
// // FIND SESSION
// //
// /////////////
func FindRespondentAccessProductSubscription(c *gin.Context) {
	subscriptionService := service.GetRespondentAccessProductSubscriptionService()

	rUser, err := auth.GetCurrentUser(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": err.Error()})
		return
	}

	id, err := uuid.Parse(c.Param("id"))
	if nil != err {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "Invalid Id provided"})
		return
	}
	// Preload("Respondent.User").Preload("Assignments").Preload("Assignments.Assignment.Product")
	subscription, err := subscriptionService.GetById(id, "Respondent.User", "Assignments.Assignment.Product")
	if nil != err {
		message := fmt.Sprintf("Could not find subscription with [id:%s]", id)
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": message})
		return
	}

	if *rUser.IsAdmin {
		/// USER IS ADMIN - ADMIN CAN VIEW ANY SESSION
	} else {

		respondant, err := auth.GetCurrentRespondent(c)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": err.Error()})
			return
		}

		if respondant.ID.String() != (*subscription.RespondentID).String() {
			message := fmt.Sprintf("Could not find your subscription with [id:%s]", id)
			c.JSON(http.StatusForbidden, gin.H{"status": "error", "message": message})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{"status": "success", "data": subscription})
}

// ////////
// // CLOSE SESSION
// //
// /////////////
func CloseRespondentAccessProductSubscription(c *gin.Context) {
	subscriptionService := service.GetRespondentAccessProductSubscriptionService()
	respondant, err := auth.GetCurrentRespondent(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": err.Error()})
		return
	}
	id, err := uuid.Parse(c.Param("id"))
	if nil != err {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "Invalid Id provided"})
		return
	}
	subscription, err := subscriptionService.GetById(id)
	if nil != err {
		message := fmt.Sprintf("Could not find subscription with [id:%s]", id)
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": message})
		return
	}
	if respondant.ID.String() != (*subscription.RespondentID).String() {
		message := fmt.Sprintf("Could not find your subscription with [id:%s]", id)
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": message})
		return
	}
	//  TODO: Close the subscription - The admin should be able to close a  given subscription for whatever reason
	if err := subscriptionService.Close(subscription); nil != err {
		message := fmt.Sprintf("Could not close subscription [id:%s]: %s", id, err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": message})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "success", "data": subscription})
}

func FindCurrentRespondentAccessProductSubscription(c *gin.Context) {
	subscriptionRepo := repository.GetRespondentAccessProductSubscriptionRepository()

	respondant, err := auth.GetCurrentRespondent(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": err.Error()})
		return
	}

	subscription, err := subscriptionRepo.GetActiveByRespondent(*respondant, "Assignments", "Assignments.Assignment", "Assignments.Assignment.Product")
	if nil != err {
		message := fmt.Sprintf("Could not find active subscription for respondent[id:%s]", respondant.ID)
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": message})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "success", "data": subscription})
}

func DeleteRespondentAccessProductSubscription(c *gin.Context) {
	subscriptionService := service.GetRespondentAccessProductSubscriptionService()

	id, err := uuid.Parse(c.Param("id"))
	if nil != err {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "Invalid Id provided"})
		return
	}
	subscription, err := subscriptionService.GetById(id)
	if nil != err {
		message := fmt.Sprintf("Could not find subscription with [id:%s]", id)
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": message})
	}

	if err := subscriptionService.Delete(subscription); nil != err {
		c.JSON(http.StatusOK, gin.H{"status": "error", "message": err.Error()})
		return
	}

	message := fmt.Sprintf("Deleted subscription \"%s\"", subscription.ID)
	c.JSON(http.StatusOK, gin.H{"data": subscription, "status": "success", "message": message})
}

func ValidateRespondentAccessProductSubscriptionCategory(category model.ProductCategory) (err error) {
	switch category {
	case model.PLACE_CITY:
		return nil
	case model.PLACE_STATE:
		return nil
	case model.PLACE_COUNTRY:
		return nil
	}
	return errors.New(fmt.Sprintf("Unsupported RespondentAccessProductSubscription Category \"%s\"", category))
}
