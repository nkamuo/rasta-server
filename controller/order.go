package controller

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/google/uuid"
	"github.com/nkamuo/rasta-server/dto"
	"github.com/nkamuo/rasta-server/model"
	"github.com/nkamuo/rasta-server/service"

	"github.com/gin-gonic/gin"
)

// GET /orders
// Get all orders
func FindOrders(c *gin.Context) {
	var orders []model.Order
	if err := model.DB.Find(&orders).Error; nil != err {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": orders})
}

func CreateOrder(c *gin.Context) {

	orderService := service.GetOrderService()

	var input dto.OrderCreationInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	order := model.Order{
		LicensePlaceNumber: *input.LicensePlaceNumber,
		Description:        input.Description,
		ModelID:            orderModel.ID,
		OwnerID:            owner.ID,
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

	if nil != input.ModelID {
		orderModel, err := modelService.GetById(*input.ModelID)
		if nil != err {
			message := fmt.Sprintf("Could not resolve the specified Model with [id:%s]: %s", input.ModelID, err.Error())
			c.JSON(http.StatusBadRequest, gin.H{"message": message, "status": "error"})
			return
		}
		order.ModelID = orderModel.ID
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

func ValidateOrderCategory(category model.OrderCategory) (err error) {
	switch category {
	case model.PRODUCT_FLAT_TIRE_SERVICE:
		return nil
	case model.PRODUCT_FUEL_DELIVERY_SERVICE:
		return nil
	case model.PRODUCT_TIRE_AIR_SERVICE:
		return nil
	case model.PRODUCT_TOWING_SERVICE:
		return nil
	case model.PRODUCT_JUMP_START_SERVICE:
		return nil
	case model.PRODUCT_KEY_UNLOCK_SERVICE:
		return nil
	}
	return errors.New(fmt.Sprintf("Unsupported Order Category \"%s\"", category))
}
