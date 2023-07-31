package controller

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/google/uuid"
	"github.com/nkamuo/rasta-server/dto"
	"github.com/nkamuo/rasta-server/model"
	"github.com/nkamuo/rasta-server/service"
	"github.com/nkamuo/rasta-server/utils/auth"

	"github.com/gin-gonic/gin"
)

func FindOrders(c *gin.Context) {
	var orders []model.Order
	if err := model.DB.Find(&orders).Error; nil != err {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": orders})
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
		if !requestingUser.IsAdmin {
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

	for _, iItem := range input.Items {

	}

	order := model.Order{
		UserID: &user.ID,
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
	orderService := service.GetOrderService()

	id, err := uuid.Parse(c.Param("id"))
	if nil != err {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "Invalid Id provided"})
		return
	}

	order, err := orderService.GetById(id)
	if err != nil {
		message := fmt.Sprintf("Could not find order with [id:%s]", id)
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": message})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": order})
}

func UpdateOrder(c *gin.Context) {
	userService := service.GetUserService()
	orderService := service.GetOrderService()
	modelService := service.GetOrderModelService()

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

	if nil != input.LicensePlaceNumber {
		order.LicensePlaceNumber = *input.LicensePlaceNumber
	}
	if nil != input.Description {
		order.Description = *input.Description
	}

	if nil != input.OwnerID {
		owner, err := userService.GetById(*input.OwnerID)
		if nil != err {
			message := fmt.Sprintf("Could not resolve the specified User with [id:%s]: %s", input.OwnerID, err.Error())
			c.JSON(http.StatusBadRequest, gin.H{"message": message, "status": "error"})
			return
		}
		order.OwnerID = owner.ID
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

func buildOrderItem(input dto.OrderItemInput, requestingUser *model.User) (orderItem *model.OrderItem, err error) {
	productService := service.GetProductService()

	product, err := productService.GetById(*input.ProductID)
	if nil != err {
		message := fmt.Sprintf("Could not resolve the specified product with [id:%s]", input.ProductID)
		return nil, errors.New(message)
	}
	if !product.Published {
		message := fmt.Sprintf("Product \"%s\" is currently not active", product.Title)
		return nil, errors.New(message)
	}

	if !requestingUser.IsAdmin {
		if input.UnitPrice != nil {
			return nil, errors.New("Invalid request. You can't specify unit price")
		}
	}

}

func validateProductOrderItemInput(product *model.Product, iItem dto.OrderItemInput) (err error) {

	return
}
