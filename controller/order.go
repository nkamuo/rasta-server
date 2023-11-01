package controller

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/google/uuid"
	"github.com/nkamuo/rasta-server/data/pagination"
	"github.com/nkamuo/rasta-server/dto"
	"github.com/nkamuo/rasta-server/model"
	"github.com/nkamuo/rasta-server/repository"
	"github.com/nkamuo/rasta-server/service"
	"github.com/nkamuo/rasta-server/utils/auth"
	"gorm.io/gorm"

	"github.com/gin-gonic/gin"
)

func FindOrders(c *gin.Context) {

	userService := service.GetUserService()

	var orders []model.Order
	var page pagination.Page
	if err := c.ShouldBindQuery(&page); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": err.Error()})
		return
	}

	requestingUser, err := auth.GetCurrentUser(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": err.Error()})
		return
	}

	query := model.DB.Preload("User").
		Preload("Items").Preload("Items.Product").
		Preload("Items.Origin").Preload("Items.Destination")

	if *requestingUser.IsAdmin {
		if user_id := c.Query("user_id"); user_id != "" {
			userID, err := uuid.Parse(user_id)
			if err != nil {
				message := fmt.Sprintf("Could not parse parameter user_id[%s] into a valid UUID: %s", user_id, err.Error())
				c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": message})
				return
			}
			if _, err := userService.GetById(userID); err != nil {
				message := fmt.Sprintf("Could not find User with [id:%s]: %s", userID, err.Error())
				c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": message})
				return
			} else {
				query = query.Where("orders.user_id = ?", userID)
			}

		}
	} else {
		query = query.Where("orders.user_id = ?", requestingUser.ID)
	}

	if status := c.Query("status"); status != "" {
		query = query.Where("orders.status = ?", status)
	}

	if err := query.Scopes(pagination.Paginate(orders, &page, query)).Find(&orders).Error; nil != err {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": err.Error()})
		return
	}
	page.Rows = orders
	c.JSON(http.StatusOK, gin.H{"data": page})
}

func CreateOrder(c *gin.Context) {

	userService := service.GetUserService()
	orderService := service.GetOrderService()
	paymentMethodService := service.GetPaymentMethodService()
	situationService := service.GetMotoristRequestSituationService()

	var input dto.OrderCreationInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	requestingUser, err := auth.GetCurrentUser(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": err.Error()})
		return
	}

	var user *model.User
	var paymentMethod *model.PaymentMethod

	if nil != input.UserID {
		if *input.UserID != requestingUser.ID && !*requestingUser.IsAdmin {
			c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "Invalid request: You may not specify userId other than yours"})
			return
		} else {

			if ruser, err := userService.GetById(*input.UserID); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": fmt.Sprintf("Could not resolve user: %s", err.Error())})
				return
			} else {
				user = ruser
			}
		}
	} else {
		user = requestingUser
	}

	if nil != input.PaymentMethodID {
		if FpaymentMethod, err := paymentMethodService.GetById(*input.PaymentMethodID); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": fmt.Sprintf("Could not resolve payment method: %s", err.Error())})
			return
		} else {
			if *FpaymentMethod.UserID != user.ID {
				c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": fmt.Sprintf("Could not resolve user payment method: %s", err.Error())})
				return
			}
			paymentMethod = FpaymentMethod
		}
	}

	if input.Items == nil || len(input.Items) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": fmt.Sprintf("You must provide at least one order item")})
		return
	}

	var Requests []model.Request

	for _, iItem := range input.Items {
		if Request, err := buildRequest(iItem, requestingUser); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": fmt.Sprintf("%s", err.Error())})
			return
		} else {
			Requests = append(Requests, *Request)
		}
	}

	var Situations []*model.MotoristRequestSituation
	for _, iSituationID := range input.Situations {
		situation, err := situationService.GetById(iSituationID)
		if err != nil {
			message := fmt.Sprintf("Error resolving situatuin with [%s]: %s", iSituationID, err.Error())
			c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": message})
			return
		}
		Situations = append(Situations, situation)
	}

	order := model.Order{
		UserID: &user.ID,
		Items:  &Requests,
		// Situations: &Situations,
	}

	if nil != paymentMethod {
		order.PaymentMethodID = &paymentMethod.ID
	}

	if err := orderService.Process(&order); nil != err {
		c.JSON(http.StatusOK, gin.H{"status": "error", "message": err.Error()})
		return
	}

	err = model.DB.Transaction(func(tx *gorm.DB) error {

		// var Sit []*model.MotoristRequestSituation
		var Sits []*model.OrderMotoristRequestSituation

		for _, sit := range Situations {
			Sits = append(Sits, &model.OrderMotoristRequestSituation{SituationID: &sit.ID})
		}
		// order.Situations = &Sit // Temporalily remove the Situations to prevent GORM from trying to insert them before the order

		order.OrderMotoristRequestSituations = Sits
		if err = tx.Save(&order).Error; err != nil {
			return err
		}
		// order.Situations = &Situations

		// if err = tx.Save(&order).Error; err != nil {
		// 	return err
		// }

		// if err := tx.Model(&order).Association("Situations").Unscoped().Append(Sits); err != nil {
		// 	return err
		// }
		// order.Situations = &Situations
		return nil
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": err.Error()})
		return
	}

	// for _, iSituationID := range input.Situations {
	// 	situation, err := situationService.GetById(iSituationID)
	// 	if err != nil {
	// 		message := fmt.Sprintf("Error resolving situatuin with [%s]: %s", iSituationID, err.Error())
	// 		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": message})
	// 		return
	// 	}
	// 	model.DB.Model(&order).Association("Situations").Append(situation)
	// 	// Situations = append(Situations, situation)
	// }

	c.JSON(http.StatusOK, gin.H{"data": order, "status": "success"})
}

func FindOrder(c *gin.Context) {

	respondentRepo := repository.GetRespondentRepository()
	sessionRepo := repository.GetRespondentSessionRepository()
	// orderService := service.GetOrderService()
	paymentService := service.GetOrderPaymentService()

	id, err := uuid.Parse(c.Param("id"))
	if nil != err {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "Invalid Id provided"})
		return
	}

	query := model.DB.Where("id = ?", id).
		Preload("User").
		Preload("Fulfilment.Responder.User").
		Preload("Fulfilment.Responder.Company").
		Preload("Fulfilment.Responder.Vehicle").
		Preload("Payment").Preload("Items").
		Preload("Adjustments").Preload("Items.Product").
		Preload("Items.Origin").Preload("Items.Destination").
		Preload("Items.FuelTypeInfo").Preload("Items.VehicleInfo")

	var order model.Order
	if err := query.First(&order).Error; nil != err {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": err.Error()})
		return
	}

	requestingUser, err := auth.GetCurrentUser(c)
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{"status": "error", "message": err.Error()})
		return
	} else {
		if !*requestingUser.IsAdmin && order.UserID.String() != requestingUser.ID.String() {
			c.JSON(http.StatusForbidden, gin.H{"status": "error", "message": err.Error()})
			return
		}
	}

	if order.Fulfilment != nil && order.Fulfilment.Coordinates == nil {

		respondent, err := respondentRepo.GetById(*order.Fulfilment.ResponderID)
		if err != nil {
			message := fmt.Sprintf("Authentication error")
			c.JSON(http.StatusForbidden, gin.H{"status": "error", "message": message})
			return
		}

		session, err := sessionRepo.GetById(*order.Fulfilment.SessionID, "Assignments.Assignment.Product")
		if nil != err {
			message := fmt.Sprintf("Could not find active session for respondent[id:%s]", respondent.ID)
			c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": message})
			return
		}
		currentCoords := session.CurrentCoordinates
		if currentCoords != nil {
			order.Fulfilment.Coordinates = currentCoords
		}
	}

	if order.Payment != nil {
		if err := paymentService.UpdatePaymentStatus(order.Payment); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": err.Error()})
			return
		}
	} else {

		// if rOrder, err := orderService.GetById(order.ID); err != nil {
		// 	c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": err.Error()})
		// 	return
		// } else {
		// 	if _, err := paymentService.InitOrderPayment(rOrder); err != nil {
		// 		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": err.Error()})
		// 		return
		// 	}
		// }

	}

	c.JSON(http.StatusOK, gin.H{"data": order, "status": "success"})
}

func UpdateOrder(c *gin.Context) {
	// userService := service.GetUserService()
	orderService := service.GetOrderService()
	paymentMethodService := service.GetPaymentMethodService()

	var input dto.OrderUpdateInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	id, err := uuid.Parse(c.Param("id"))
	if nil != err {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "Invalid Id provided"})
		return
	}
	order, err := orderService.GetById(id)
	if nil != err {
		message := fmt.Sprintf("Could not find order with [id:%s]", id)
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": message})
	}

	// requestingUser, err := auth.GetCurrentUser(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": err.Error()})
		return
	}

	var user *model.User
	var paymentMethod *model.PaymentMethod

	if nil != input.PaymentMethodID {
		if FpaymentMethod, err := paymentMethodService.GetById(*input.PaymentMethodID); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": fmt.Sprintf("Could not resolve payment method: %s", err.Error())})
			return
		} else {
			if *FpaymentMethod.UserID != user.ID {
				c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": fmt.Sprintf("Could not resolve user payment method: %s", err.Error())})
				return
			}
			paymentMethod = FpaymentMethod
		}
	}

	if nil != paymentMethod {
		order.PaymentMethodID = &paymentMethod.ID
	}

	if err := orderService.Process(order); nil != err {
		c.JSON(http.StatusOK, gin.H{"status": "error", "message": err.Error()})
		return
	}

	if err := orderService.Save(order); nil != err {
		c.JSON(http.StatusOK, gin.H{"status": "error", "message": err.Error()})
		return
	}

	message := fmt.Sprintf("Updated model \"%s\"", order.ID)
	c.JSON(http.StatusOK, gin.H{"data": order, "status": "success", "message": message})
}

func DeleteOrder(c *gin.Context) {
	orderService := service.GetOrderService()

	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		message := fmt.Sprintf("Error Passing\"%s\" into a valid UUID", c.Param("id"))
		c.JSON(http.StatusOK, gin.H{"status": "error", "message": message})

	}

	order, err := orderService.GetById(id)
	if err != nil {
		message := fmt.Sprintf("Could not find order with [id:%s]", id)
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": message})
		return
	}

	if rUser, err := auth.GetCurrentUser(c); err != nil {
		message := fmt.Sprintf("Authentication Problem")
		c.JSON(http.StatusForbidden, gin.H{"status": "error", "message": message})
		return

	} else {
		if !*rUser.IsAdmin && rUser.ID.String() != order.UserID.String() {
			message := fmt.Sprintf("Restricated Resource")
			c.JSON(http.StatusForbidden, gin.H{"status": "error", "message": message})
		}
	}

	if err := orderService.Delete(order); err != nil {
		message := fmt.Sprintf("An error occured: %s", err.Error())
		c.JSON(http.StatusForbidden, gin.H{"status": "error", "message": message})
	}

	message := fmt.Sprintf("Deleted order \"%s\"", order.ID)
	c.JSON(http.StatusOK, gin.H{"data": order, "status": "success", "message": message})
}

func CompleteOrder(c *gin.Context) {
	orderService := service.GetOrderService()

	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		message := fmt.Sprintf("Error Passing\"%s\" into a valid UUID", c.Param("id"))
		c.JSON(http.StatusOK, gin.H{"status": "error", "message": message})

	}

	order, err := orderService.GetById(id)
	if err != nil {
		message := fmt.Sprintf("Could not find order with [id:%s]", id)
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": message})
		return
	}

	if rUser, err := auth.GetCurrentUser(c); err != nil {
		message := fmt.Sprintf("Authentication Problem")
		c.JSON(http.StatusForbidden, gin.H{"status": "error", "message": message})
		return

	} else {
		if !*rUser.IsAdmin && rUser.ID.String() != order.UserID.String() {
			message := fmt.Sprintf("Restricated Resource")
			c.JSON(http.StatusForbidden, gin.H{"status": "error", "message": message})
		}
	}

	if err := orderService.CompleteOrder(order, false, nil); err != nil {
		message := fmt.Sprintf("An error occured: %s", err.Error())
		c.JSON(http.StatusForbidden, gin.H{"status": "error", "message": message})
	}

	query := model.DB.Where("id = ?", id).
		Preload("User").
		Preload("Fulfilment.Responder.User").
		Preload("Fulfilment.Responder.Vehicle").
		Preload("Payment").Preload("Items").
		Preload("Adjustments").Preload("Items.Product").
		Preload("Items.Origin").Preload("Items.Destination").
		Preload("Items.FuelTypeInfo").Preload("Items.VehicleInfo")

	var rOrder model.Order
	if err := query.First(&rOrder).Error; nil != err {
		// REFETCHINING ORDER FAILED - DO NOTHING
	} else {
		order = &rOrder
	}

	message := fmt.Sprintf("Order completed order \"%s\"", order.ID)
	c.JSON(http.StatusOK, gin.H{"data": order, "status": "success", "message": message})
}

func buildRequest(input dto.RequestInput, requestingUser *model.User) (Request *model.Request, err error) {
	placeService := service.GetPlaceService()
	productService := service.GetProductService()
	locationService := service.GetLocationService()
	// fuelTypeService := service.GetFuelTypeService()
	fuelTypeRepository := repository.GetFuelTypeRepository()

	var vehicleInfo *model.RequestVehicleInfo
	var fuelTypeInfo *model.RequestFuelTypeInfo
	var rate uint64
	var quantity uint64 = 1

	if nil != input.Quantity {
		quantity = *input.Quantity
	}

	var origin, destination *model.Location

	product, err := productService.GetById(input.ProductID)
	if nil != err {
		message := fmt.Sprintf("Could not resolve the specified product with [id:%s]", input.ProductID)
		return nil, errors.New(message)
	}
	place, err := placeService.GetById(product.PlaceID)
	if nil != err {
		message := fmt.Sprintf("Could not resolve the specified Place with [id:%s]", product.PlaceID)
		return nil, errors.New(message)
	}

	if !*place.Active {
		message := fmt.Sprintf("Place \"%s\" is currently not active", place.Name)
		return nil, errors.New(message)
	}

	if !*product.Published {
		message := fmt.Sprintf("Product \"%s\" is currently not active", product.Title)
		return nil, errors.New(message)
	}

	if input.Rate != nil {
		if !*requestingUser.IsAdmin {
			return nil, errors.New("Invalid request. You can't specify unit price")
		} else {
			rate = *input.Rate
		}
	} else {
		rate = product.Rate
	}

	if nil != input.Origin {
		origin, err = locationService.Resolve(*input.Origin)
		if nil != err {
			return nil, errors.New(
				fmt.Sprintf(
					"Could not resolve origin location \"%s\". Failed with error :%s",
					*input.Origin,
					err.Error(),
				),
			)
		}
	}

	if nil == input.Destination {
		return nil, errors.New(fmt.Sprintf("Provide the destination location for \"%v\"", product.Title))
	} else {
		destination, err = locationService.Resolve(*input.Destination)
		if nil != err {
			return nil, errors.New(
				fmt.Sprintf(
					"Could not resolve destination location \"%s\". Failed with error :%s",
					*input.Destination,
					err.Error(),
				),
			)
		}
	}

	if nil == input.VehicleInfo {
		return nil, errors.New("You have to provide vehicle Information")
	}
	vehicleInfo, err = buildVehicleInfo(*input.VehicleInfo)
	if err != nil {
		return nil, err
	}

	switch product.Category {
	case model.PRODUCT_FLAT_TIRE_SERVICE:

		break

	case model.PRODUCT_TOWING_SERVICE:

		if origin == nil || destination == nil {
			return nil, errors.New(fmt.Sprintf("Origin and Destination are requried for %v", model.PRODUCT_TOWING_SERVICE))
		}

		distanceInfo, err := locationService.GetDistance(origin, destination)
		if err != nil {
			return nil, err
		}

		towRate, err := service.GetTowingPlaceRateService().GetByPlaceAndDistance(*place, int64(distanceInfo.Distance.Value))
		if err != nil {
			if err.Error() != "record not found" {
				return nil, err
			} else {
				rate = product.Rate
			}
		} else {
			rate = *towRate.Rate
		}

		break

	case model.PRODUCT_FUEL_DELIVERY_SERVICE:
		if nil == input.FuelInfo {
			return nil, errors.New("You have to provide fuel Information")
		}

		fuelType, err := fuelTypeRepository.GetByCode(input.FuelInfo.FuelTypeCode)
		if nil != err {
			return nil, errors.New(fmt.Sprintf("Unsupported Fuel type[%s]: %v", input.FuelInfo.FuelTypeCode, err.Error()))
		}
		if placeRate, err := fuelTypeRepository.GetRateForTypeInPlace(*fuelType, *place); nil == err && nil != placeRate {
			rate = placeRate.Rate
		} else {
			rate = fuelType.Rate
		}

		fuelTypeInfo = &model.RequestFuelTypeInfo{
			FuelTypeCode: fuelType.Code,
			FuelTypeID:   &fuelType.ID,
		}

		break
	}

	Request = &model.Request{
		ProductID: &product.ID,
		Rate:      rate,
		Quantity:  quantity,
		//
		VehicleInfo: vehicleInfo,
		//
		FuelTypeInfo: fuelTypeInfo,
	}

	if nil != origin {
		Request.OriginID = &origin.ID
	}
	if nil != destination {
		Request.DestinationID = &destination.ID
	}

	return Request, nil

}

func buildVehicleInfo(input dto.RequestVehicleInformationInput) (vehicleInfo *model.RequestVehicleInfo, err error) {

	// if nil == input.VehicleDescription {
	// 	return nil, errors.New("You have to provide some description about flat tire service")
	// }
	// if dLength := len(*input.VehicleDescription); 20 > dLength || dLength > 500 {
	// 	return nil, errors.New("Description must be upto 20 and less than 500 charachers")
	// }

	if nil == input.Make {
		return nil, errors.New("Vehicle Make not Specified")
	}
	if nil == input.Model {
		return nil, errors.New("Vehicle Model not Specified")
	}
	if nil == input.Color {
		return nil, errors.New("Vehicle Color not Specified")
	}
	if nil == input.LicensePlateNumber {
		return nil, errors.New("Vehicle License Plate Number not Specified")
	}
	// if dLength := len(*input.VehicleDescription); 20 > dLength || dLength > 500 {
	// 	return nil, errors.New("Description must be upto 20 and less than 500 charachers")
	// }

	return &model.RequestVehicleInfo{
		MakeName:           input.Make,
		ModelName:          input.Model,
		BodyColor:          input.Color,
		LicensePlateNumber: input.LicensePlateNumber,
	}, nil

}

func validateProductRequestInput(product *model.Product, iItem dto.RequestInput) (err error) {

	return
}
