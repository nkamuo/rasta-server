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

	"github.com/gin-gonic/gin"
)

func FindOrders(c *gin.Context) {
	var orders []model.Order
	var page pagination.Page
	if err := c.ShouldBindQuery(&page); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": err.Error()})
		return
	}

	if err := model.DB.Preload("User").Scopes(pagination.Paginate(orders, &page, model.DB)).Find(&orders).Error; nil != err {
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
		if !*requestingUser.IsAdmin {
			c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "Invalid request"})
		} else {

			if ruser, err := userService.GetById(*input.UserID); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": fmt.Sprintf("Could not resolve user: %s", err.Error())})
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

	order := model.Order{
		UserID: &user.ID,
		Items:  &Requests,
	}

	if nil != paymentMethod {
		order.PaymentMethodID = &paymentMethod.ID
	}

	if err := orderService.Save(&order); nil != err {
		c.JSON(http.StatusOK, gin.H{"status": "error", "message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": order, "status": "success"})
}

func FindOrder(c *gin.Context) {

	id, err := uuid.Parse(c.Param("id"))
	if nil != err {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "Invalid Id provided"})
		return
	}

	var order model.Order
	if err := model.DB.Where("id = ?", id).Preload("Items").First(&order).Error; nil != err {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": order})
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

	if err := orderService.Save(order); nil != err {
		c.JSON(http.StatusOK, gin.H{"status": "error", "message": err.Error()})
		return
	}

	message := fmt.Sprintf("Updated model \"%s\"", order.ID)
	c.JSON(http.StatusOK, gin.H{"data": order, "status": "success", "message": message})
}

func DeleteOrder(c *gin.Context) {
	id := c.Param("id")

	var order model.Order

	if err := model.DB.Where("id = ?", id).First(&order).Error; err != nil {
		message := fmt.Sprintf("Could not find order with [id:%s]", id)
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": message})
		return
	}
	model.DB.Delete(&order)
	message := fmt.Sprintf("Deleted order \"%s\"", order.ID)
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

	product, err := productService.GetById(*input.ProductID)
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

	switch product.Category {
	case model.PRODUCT_FLAT_TIRE_SERVICE:

		if nil == input.VehicleInfo {
			return nil, errors.New("You have to provide vehicle Information")
		}
		if nil == input.VehicleInfo.VehicleDescription {
			return nil, errors.New("You have to provide some description about flat tire service")
		}
		if dLength := len(*input.VehicleInfo.VehicleDescription); 20 > dLength || dLength > 500 {
			return nil, errors.New("Description must be upto 20 and less than 500 charachers")
		}
		break

	case model.PRODUCT_FUEL_DELIVERY_SERVICE:
		if nil == input.FuelInfo {
			return nil, errors.New("You have to provide vehicle Information")
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

func validateProductRequestInput(product *model.Product, iItem dto.RequestInput) (err error) {

	return
}
