package controller

import (
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/nkamuo/rasta-server/dto"
	"github.com/nkamuo/rasta-server/repository"
	"github.com/nkamuo/rasta-server/service"
	"github.com/nkamuo/rasta-server/utils/auth"

	"github.com/gin-gonic/gin"
)

func ClientVerifyOrderRespondentDetails(c *gin.Context) {

	orderService := service.GetOrderService()
	fulfilmentService := service.GetOrderFulfilmentService()
	// respondentService := service.GetClientService()

	rUser, err := auth.GetCurrentUser(c)
	if err != nil {
		message := fmt.Sprintf("Authentication error")
		c.JSON(http.StatusForbidden, gin.H{"status": "error", "message": message})
		return
	}

	id, err := uuid.Parse(c.Param("id"))
	if nil != err {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "Invalid Id provided"})
		return
	}

	order, err := orderService.GetById(id)
	if err != nil {
		message := fmt.Sprintf("Could not find Order with [id:%s]", id)
		c.JSON(http.StatusNotFound, gin.H{"status": "error", "message": message})
		return
	}

	if order.UserID.String() != rUser.ID.String() {
		message := fmt.Sprintf("You may not access this resource")
		c.JSON(http.StatusForbidden, gin.H{"status": "error", "message": message})
		return
	}

	if order.FulfilmentID == nil {
		message := fmt.Sprintf("Order [id:%s] is not assigned to a responder yet", id)
		c.JSON(http.StatusExpectationFailed, gin.H{"status": "error", "message": message})
		return
	}
	fulfilment, err := fulfilmentService.GetById(*order.FulfilmentID, "Order.User")
	if err != nil {
		message := fmt.Sprintf("Error Fetching order[id:%s] fulfilment details: %s", id, err.Error())
		c.JSON(http.StatusNotFound, gin.H{"status": "error", "message": message})
		return
	}

	now := time.Now()
	fulfilment.VerifiedResponderAt = &now

	if err := fulfilmentService.Save(fulfilment); err != nil {
		message := fmt.Sprintf("Task failed: %s", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": message})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": fulfilment, "status": "success"})
}

func ClientConfirmCompleteOrder(c *gin.Context) {

	orderService := service.GetOrderService()
	fulfilmentService := service.GetOrderFulfilmentService()

	id, err := uuid.Parse(c.Param("id"))
	if nil != err {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "Invalid Id provided"})
		return
	}

	var input dto.ClientOrderConfirmationRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": err.Error()})
		return
	}

	rUser, err := auth.GetCurrentUser(c)
	if err != nil {
		message := fmt.Sprintf("Authentication error")
		c.JSON(http.StatusForbidden, gin.H{"status": "error", "message": message})
		return
	}

	order, err := orderService.GetById(id)
	if err != nil {
		message := fmt.Sprintf("Could not find Order with [id:%s]", id)
		c.JSON(http.StatusNotFound, gin.H{"status": "error", "message": message})
		return
	}

	if !*rUser.IsAdmin && order.UserID.String() != rUser.ID.String() {
		message := fmt.Sprintf("You may not access this resource")
		c.JSON(http.StatusForbidden, gin.H{"status": "error", "message": message})
		return
	}

	if order.FulfilmentID == nil {
		message := fmt.Sprintf("Order [id:%s] is not assigned yet", id)
		c.JSON(http.StatusNotFound, gin.H{"status": "error", "message": message})
		return
	}

	fulfilment, err := fulfilmentService.GetById(*order.FulfilmentID)
	if err != nil {
		message := fmt.Sprintf("Error Fetching order[id:%s] fulfilment details: %s", id, err.Error())
		c.JSON(http.StatusNotFound, gin.H{"status": "error", "message": message})
		return
	}

	// if err := fulfilmentService.Save(fulfilment); err != nil {
	// 	message := fmt.Sprintf("Task failed: %s", err.Error())
	// 	c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": message})
	// 	return
	// }

	if err := orderService.CompleteOrder(order, false, &input); err != nil {
		message := fmt.Sprintf("Task failed: %s", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": message})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": fulfilment, "status": "success"})
}

func ClientCancelOrder(c *gin.Context) {

	orderRepo := repository.GetOrderRepository()
	orderService := service.GetOrderService()
	fulfilmentService := service.GetOrderFulfilmentService()

	id, err := uuid.Parse(c.Param("id"))
	if nil != err {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "Invalid Id provided"})
		return
	}

	rUser, err := auth.GetCurrentUser(c)
	if err != nil {
		message := fmt.Sprintf("Authentication error")
		c.JSON(http.StatusForbidden, gin.H{"status": "error", "message": message})
		return
	}

	order, err := orderService.GetById(id)
	if err != nil {
		message := fmt.Sprintf("Could not find Order with [id:%s]", id)
		c.JSON(http.StatusNotFound, gin.H{"status": "error", "message": message})
		return
	}

	if order.UserID.String() != rUser.ID.String() {
		message := fmt.Sprintf("You may not access this resource")
		c.JSON(http.StatusForbidden, gin.H{"status": "error", "message": message})
		return
	}

	if order.FulfilmentID == nil {
		message := fmt.Sprintf("Order [id:%s] is not assigned yet", id)
		c.JSON(http.StatusNotFound, gin.H{"status": "error", "message": message})
		return
	}

	fulfilment, err := fulfilmentService.GetById(*order.FulfilmentID)
	if err != nil {
		message := fmt.Sprintf("Error Fetching order[id:%s] fulfilment details: %s", id, err.Error())
		c.JSON(http.StatusNotFound, gin.H{"status": "error", "message": message})
		return
	}

	if err := orderRepo.Update(order, map[string]interface{}{"fulfilment_id": nil}); err != nil {
		message := fmt.Sprintf("Task failed: %s", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": message})
		return
	}

	if err := fulfilmentService.Delete(fulfilment); err != nil {
		message := fmt.Sprintf("Task failed: %s", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": message})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": fulfilment, "status": "success"})
}
