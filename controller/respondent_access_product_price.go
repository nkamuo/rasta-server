package controller

import (
	"fmt"
	"net/http"

	"github.com/google/uuid"
	"github.com/stripe/stripe-go/v74"
	"github.com/stripe/stripe-go/v74/price"
	"gorm.io/gorm"

	// "github.com/mitchellh/mapstructure"
	"github.com/nkamuo/rasta-server/data/pagination"
	"github.com/nkamuo/rasta-server/dto"
	"github.com/nkamuo/rasta-server/initializers"
	"github.com/nkamuo/rasta-server/model"
	"github.com/nkamuo/rasta-server/service"
	"github.com/nkamuo/rasta-server/utils/auth"

	"github.com/gin-gonic/gin"
	// "github.com/gin-gonic/gin"
)

func FindRespondentAccessProductPrices(c *gin.Context) {

	// respondentRepo := repository.GetRespondentRepository()
	// placeRepo := repository.GetPlaceRepository()

	var prices []model.RespondentAccessProductPrice
	var page pagination.Page
	if err := c.ShouldBindQuery(&page); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": err.Error()})
		return
	}

	query := model.DB.Order("upto ASC").Order("unit_price ASC").Order("created_at ASC") //.Preload("Place")

	// TODO: change `price_product_type` to `price_id` and use it to match the Stripe price id the user is subscribed to
	if price_product_type := c.Query("price_product_type"); price_product_type != "" {

		query = query.Where("product_type = ?", price_product_type)
	}

	// if status := c.Query("status"); status != "" {
	// 	query = query.Where("active = ?", true)
	// }

	if err := query.Scopes(pagination.Paginate(prices, &page, query)).Find(&prices).Error; nil != err {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": err.Error()})
		return
	}
	page.Rows = prices
	c.JSON(http.StatusOK, gin.H{"data": page})
}

func CreateRespondentAccessProductPrice(c *gin.Context) {
	// priceService := service.GetRespondentAccessProductPriceService()

	config, err := initializers.LoadConfig()
	if err != nil {
		message := fmt.Sprintf("Error fetching prices: %s", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": message})
		return
	}

	rUser, err := auth.GetCurrentUser(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": err.Error()})
		return
	}

	if !*rUser.IsAdmin {
		message := fmt.Sprintf("An error occured: %s", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": message})
		return
	}

	var input dto.RespondentAccessProductPriceCreationInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": err.Error()})
		return
	}
	var localPrice *model.RespondentAccessProductPrice
	var productId string

	if input.Mode == model.ACCESS_PRODUCT_TYPE_PURCHASE {
		productId = config.STRIPE_RESPONDENT_PURCHASE_PRODUCT_ID
	} else if input.Mode == model.ACCESS_PRODUCT_TYPE_SUBSCIPTION {
		productId = config.STRIPE_RESPONDENT_SUBSCRIPTION_PRODUCT_ID
	} else {
		message := fmt.Sprintf("Invalid Product type provided: %s", input.Mode)
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": message})
		return
	}

	err = model.DB.Transaction(func(tx *gorm.DB) error {

		params := &stripe.PriceParams{
			Currency:   stripe.String(string(stripe.CurrencyUSD)),
			UnitAmount: stripe.Int64(int64(input.UnitPrice)),
			Nickname:   input.Label,
			Product:    stripe.String(productId),
		}

		localPrice = &model.RespondentAccessProductPrice{
			// StripePriceID: &sPrice.ID,
			Label:     input.Label,
			UnitPrice: &input.UnitPrice,
			// Upto:        input.UpTo,
			ProductType: &input.Mode,
			Description: input.Description,
			// Active: 	  input.Active,
		}

		if input.UpTo != nil {
			localPrice.Upto = input.UpTo
		} else {
			*localPrice.Upto = 1
		}

		if input.Active != nil {
			localPrice.Active = input.Active
			params.Active = stripe.Bool(*input.Active)
		} else {
			localPrice.Active = stripe.Bool(true)
		}

		sPrice, err := price.New(params)
		if err != nil {
			return err
		}
		localPrice.StripePriceID = &sPrice.ID
		if err := tx.Save(localPrice).Error; nil != err {
			sPrice, err = price.Update(sPrice.ID, &stripe.PriceParams{
				Active: stripe.Bool(false),
			})
			return err
		}
		return err
	})
	if nil != err {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "success", "data": localPrice})
}

// ////////
// // FIND PRICE
// //
// /////////////
func FindRespondentAccessProductPrice(c *gin.Context) {
	priceService := service.GetRespondentAccessProductPriceService()

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
	_price, err := priceService.GetById(id)
	if nil != err {
		message := fmt.Sprintf("Could not find price with [id:%s]", id)
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": message})
		return
	}

	if *rUser.IsAdmin {
		/// USER IS ADMIN - ADMIN CAN VIEW ANY PRICE
	} else {

		message := fmt.Sprintf("Could not find your price with [id:%s]", id)
		c.JSON(http.StatusForbidden, gin.H{"status": "error", "message": message})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "success", "data": _price})
}

func UpdateRespondentAccessProductPrice(c *gin.Context) {
	priceService := service.GetRespondentAccessProductPriceService()

	config, err := initializers.LoadConfig()
	if err != nil {
		message := fmt.Sprintf("Error fetching prices: %s", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": message})
		return
	}

	rUser, err := auth.GetCurrentUser(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": err.Error()})
		return
	}

	if !*rUser.IsAdmin {
		message := fmt.Sprintf("An error occured: %s", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": message})
		return
	}

	id, err := uuid.Parse(c.Param("id"))
	if nil != err {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "Invalid Id provided"})
		return
	}
	_price, err := priceService.GetById(id)
	if nil != err {
		message := fmt.Sprintf("Could not find price with [id:%s]", id)
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": message})
		return
	}

	var input dto.RespondentAccessProductPriceUpdateInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": err.Error()})
		return
	}
	var productId string
	params := &stripe.PriceParams{}

	if input.Mode != nil {
		if *input.Mode == model.ACCESS_PRODUCT_TYPE_PURCHASE {
			productId = config.STRIPE_RESPONDENT_PURCHASE_PRODUCT_ID
		} else if *input.Mode == model.ACCESS_PRODUCT_TYPE_SUBSCIPTION {
			productId = config.STRIPE_RESPONDENT_SUBSCRIPTION_PRODUCT_ID
		} else {
			message := fmt.Sprintf("Invalid Product type provided: %s", *input.Mode)
			c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": message})
			return
		}
	}

	if input.Active != nil {
		_price.Active = input.Active
		params.Active = stripe.Bool(*input.Active)
	}
	if input.UpTo != nil {
		_price.Upto = input.UpTo
	}
	if input.Mode != nil {
		_price.ProductType = input.Mode
		params.Product = stripe.String(productId)
	}
	if input.Label != nil {
		params.Nickname = stripe.String(*input.Label)
		_price.Label = input.Label
	}
	if input.Description != nil {
		_price.Description = input.Description
		// params.ProductData = &stripe.PriceProductDataParams{
		// 	// Description: stripe.String(*input.Description),
		// }
	}

	if input.UnitPrice != nil && input.UnitPrice != _price.UnitPrice {
		_price.UnitPrice = input.UnitPrice
		params.UnitAmount = stripe.Int64(int64(*input.UnitPrice))
		params.Currency = stripe.String(string(stripe.CurrencyUSD))
		if params.Product == nil {
			params.Product = stripe.String(*_price.ProductID())
		}
		if _price.StripePriceID != nil {
			// params.Currency = stripe.String(string(stripe.CurrencyUSD))
			_, err = price.Update(*_price.StripePriceID, &stripe.PriceParams{
				Active: stripe.Bool(false),
			})
			if err != nil {
				message := fmt.Sprintf("Could not overide Existing price: %s", *_price.StripePriceID)
				c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": message})
				return
			}
		}
		sPrice, err := price.New(params)
		if err != nil {
			message := fmt.Sprintf("Could not create a new stripe price: %s", err.Error())
			c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": message})
			return
		}
		_price.StripePriceID = &sPrice.ID

		// kfjd
	}

	err = model.DB.Transaction(func(tx *gorm.DB) error {

		if input.UnitPrice == nil {
			_, err := price.Update(*_price.StripePriceID, params)

			if err != nil {
				return err
			}
		}
		if err := tx.Save(_price).Error; nil != err {
			// sPrice, err = price.Update(sPrice.ID, &stripe.PriceParams{
			// 	Active: stripe.Bool(false),
			// })
			return err
		}
		return err
	})
	if nil != err {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "success", "data": _price})
}

func DeleteRespondentAccessProductPrice(c *gin.Context) {
	priceService := service.GetRespondentAccessProductPriceService()

	id, err := uuid.Parse(c.Param("id"))
	if nil != err {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "Invalid Id provided"})
		return
	}
	_price, err := priceService.GetById(id)
	if nil != err {
		message := fmt.Sprintf("Could not find price with [id:%s]", id)
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": message})
	}

	// if err := priceService.Delete(price); nil != err {
	// 	c.JSON(http.StatusOK, gin.H{"status": "error", "message": err.Error()})
	// 	return
	// }
	err = model.DB.Transaction(func(tx *gorm.DB) error {
		_, err = price.Update(*_price.StripePriceID, &stripe.PriceParams{
			Active: stripe.Bool(false),
		})
		if err != nil {
			return err
		}
		if err := tx.Delete(_price).Error; nil != err {
			return err
		}
		return err
	})

	if nil != err {
		c.JSON(http.StatusOK, gin.H{"status": "error", "message": err.Error()})
		return
	}

	message := fmt.Sprintf("Deleted price \"%s\"", _price.ID)
	c.JSON(http.StatusOK, gin.H{"data": _price, "status": "success", "message": message})
}
