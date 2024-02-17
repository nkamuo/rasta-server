package controller

import (
	"errors"
	"fmt"
	"net/http"
	"time"

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

func RespondentClaimOrder(c *gin.Context) {

	orderService := service.GetOrderService()
	respondentRepo := repository.GetRespondentRepository()
	sessionRepo := repository.GetRespondentSessionRepository()
	// respondentService := service.GetRespondentService()

	requestingUser, err := auth.GetCurrentUser(c)
	if err != nil {
		message := fmt.Sprintf("Authentication error")
		c.JSON(http.StatusForbidden, gin.H{"status": "error", "message": message})
		return
	}

	respondent, err := respondentRepo.GetByUser(*requestingUser, "Vehicle")
	if err != nil {
		message := fmt.Sprintf("Authentication error")
		c.JSON(http.StatusForbidden, gin.H{"status": "error", "message": message})
		return
	}

	if respondent.Vehicle == nil || *respondent.Vehicle.Published != true {
		message := fmt.Sprintf("You need to select a verified vehicle before you can accept requests")
		c.JSON(http.StatusConflict, gin.H{"status": "error", "message": message})
		return
	}

	session, err := sessionRepo.GetActiveByRespondent(*respondent)
	if err != nil {
		message := fmt.Sprintf("Error identifying your session: %s", err.Error())
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

	if err := orderService.AssignResponder(order, session); err != nil {
		message := fmt.Sprintf("Task failed: %s", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": message})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": order, "status": "success"})
}

func RespondentVerifyOrderClientDetails(c *gin.Context) {

	orderService := service.GetOrderService()
	respondentRepo := repository.GetRespondentRepository()
	fulfilmentService := service.GetOrderFulfilmentService()
	// respondentService := service.GetRespondentService()

	requestingUser, err := auth.GetCurrentUser(c)
	if err != nil {
		message := fmt.Sprintf("Authentication error")
		c.JSON(http.StatusForbidden, gin.H{"status": "error", "message": message})
		return
	}

	respondent, err := respondentRepo.GetByUser(*requestingUser)
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

	if fulfilment.ResponderID.String() != respondent.ID.String() {
		message := fmt.Sprintf("You may not access this resource")
		c.JSON(http.StatusForbidden, gin.H{"status": "error", "message": message})
		return
	}

	now := time.Now()
	fulfilment.VerifiedClientAt = &now
	order.Status = model.ORDER_STATUS_CLIENT_CONFIRMED

	err = model.DB.Transaction(func(tx *gorm.DB) (err error) {
		if err = tx.Save(fulfilment).Error; err != nil {
			return err
		}
		if err = tx.Save(order).Error; err != nil {
			return err
		}
		return nil
	})

	if err != nil {
		message := fmt.Sprintf("Task failed: %s", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": message})
		return
	}

	// if err := fulfilmentService.Save(fulfilment); err != nil {
	// 	message := fmt.Sprintf("Task failed: %s", err.Error())
	// 	c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": message})
	// 	return
	// }

	c.JSON(http.StatusOK, gin.H{"data": fulfilment, "status": "success"})
}

func RespondentConfirmCompleteOrder(c *gin.Context) {

	orderService := service.GetOrderService()
	respondentRepo := repository.GetRespondentRepository()
	fulfilmentService := service.GetOrderFulfilmentService()
	// respondentService := service.GetRespondentService()

	var input dto.RespondentOrderCompletionInput
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

	respondent, err := respondentRepo.GetByUser(*rUser)
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

	if fulfilment.ResponderID.String() != respondent.ID.String() {
		message := fmt.Sprintf("You may not access this resource")
		c.JSON(http.StatusForbidden, gin.H{"status": "error", "message": message})
		return
	}

	now := time.Now()
	fulfilment.ResponderConfirmedAt = &now
	order.Status = model.ORDER_STATUS_COMPLETED_BY_RESPONDENT

	// if(input.ClientPaidCash){
	order.ClientPaidCash = &input.ClientPaidCash
	// }

	err = model.DB.Transaction(func(tx *gorm.DB) (err error) {
		if err = tx.Save(fulfilment).Error; err != nil {
			return err
		}
		if err = tx.Save(order).Error; err != nil {
			return err
		}
		return nil
	})

	if err := fulfilmentService.Save(fulfilment); err != nil {
		message := fmt.Sprintf("Task failed: %s", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": message})
		return
	}

	// nOrder, err := getOrderById(id)
	// if nil != err {
	// 	message := fmt.Sprintf("Error Fetching order[id:%s] details: %s", id, err.Error())
	// 	c.JSON(http.StatusNotFound, gin.H{"status": "error", "message": message})
	// 	return
	// }

	c.JSON(http.StatusOK, gin.H{"data": fulfilment, "status": "success"})
}

func RespondentCancelOrder(c *gin.Context) {

	// orderRepo := repository.GetOrderRepository()
	orderService := service.GetOrderService()
	respondentRepo := repository.GetRespondentRepository()
	fulfilmentService := service.GetOrderFulfilmentService()
	// respondentService := service.GetRespondentService()

	requestingUser, err := auth.GetCurrentUser(c)
	if err != nil {
		message := fmt.Sprintf("Authentication error")
		c.JSON(http.StatusForbidden, gin.H{"status": "error", "message": message})
		return
	}

	respondent, err := respondentRepo.GetByUser(*requestingUser)
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

	if fulfilment.ResponderID.String() != respondent.ID.String() {
		message := fmt.Sprintf("You may not access this resource")
		c.JSON(http.StatusForbidden, gin.H{"status": "error", "message": message})
		return
	}

	err = model.DB.Transaction(func(tx *gorm.DB) error {
		err := tx.Model(order).
			Updates(map[string]interface{}{
				"fulfilment_id": nil,
				"status":        model.ORDER_STATUS_PENDING,
			}).Error
		return err
	})

	if err != nil {
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

func RespondentUpdateOrderPayment(c *gin.Context) {

	orderRepo := repository.GetOrderRepository()
	// orderService := service.GetOrderService()
	orderPaymentService := service.GetOrderPaymentService()
	// orderRepo := repository.GetOrderRepository();
	respondentRepo := repository.GetRespondentRepository()
	fulfilmentService := service.GetOrderFulfilmentService()
	// respondentService := service.GetRespondentService()

	var input dto.RespondentOrderCompletionInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": err.Error()})
		return
	}

	requestingUser, err := auth.GetCurrentUser(c)
	if err != nil {
		message := fmt.Sprintf("Authentication error")
		c.JSON(http.StatusForbidden, gin.H{"status": "error", "message": message})
		return
	}

	respondent, err := respondentRepo.GetByUser(*requestingUser)
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

	order, err := orderRepo.GetById(id)
	if err != nil {
		message := fmt.Sprintf("Could not find Order with [id:%s]", id)
		c.JSON(http.StatusNotFound, gin.H{"status": "error", "message": message})
		return
	}

	rOrder, err := orderRepo.GetById(id, "Payment")
	if err != nil {
		message := fmt.Sprintf("Could not find Order with [id:%s]", id)
		c.JSON(http.StatusNotFound, gin.H{"status": "error", "message": message})
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

	if fulfilment.ResponderID.String() != respondent.ID.String() {
		message := fmt.Sprintf("You may not access this resource")
		c.JSON(http.StatusForbidden, gin.H{"status": "error", "message": message})
		return
	}

	var payment = rOrder.Payment

	if nil == payment {
		payment, err = orderPaymentService.InitOrderPayment(order)
		if err != nil {
			message := fmt.Sprintf("An error Occured: %s", err.Error())
			c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": message})
			return
		}
	}

	if input.ClientPaidCash {
		payment.Status = "paid_on_delivery"
	} else {
		if payment.Status == "paid_on_delivery" {
			payment.Status = "processing"
		}
	}

	if orderPaymentService.Save(payment); err != nil {
		message := fmt.Sprintf("An error Occured: %s", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": message})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": fulfilment, "status": "success"})
}

func FindOrderForRespondent(c *gin.Context) {

	respondentRepo := repository.GetRespondentRepository()
	sessionRepo := repository.GetRespondentSessionRepository()
	locationService := service.GetLocationService()

	id, err := uuid.Parse(c.Param("id"))
	if nil != err {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "Invalid Id provided"})
		return
	}

	requestingUser, err := auth.GetCurrentUser(c)
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{"status": "error", "message": err.Error()})
		return
	} else {

	}

	respondant, err := respondentRepo.GetByUser(*requestingUser)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": err.Error()})
		return
	}

	session, err := sessionRepo.GetActiveByRespondent(*respondant, "Assignments.Assignment.Product")
	if nil != err {
		message := fmt.Sprintf("Could not find active session for respondent[id:%s]", respondant.ID)
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": message})
		return
	}

	query := model.DB.Where("id = ?", id).
		Preload("Fulfilment.Responder.User").
		Preload("User").Preload("OrderMotoristRequestSituations.Situation").
		Preload("Payment").Preload("Items").
		Preload("Adjustments").Preload("Items.Product").
		Preload("Items.Origin").Preload("Items.Destination").
		Preload("Items.FuelTypeInfo").Preload("Items.VehicleInfo")

	var order model.Order
	if err := query.First(&order).Error; nil != err {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": err.Error()})
		return
	}

	if !*requestingUser.IsAdmin && (order.UserID != nil && order.UserID.String() != requestingUser.ID.String()) {
		// c.JSON(http.StatusForbidden, gin.H{"status": "error", "message": "Not Authorized"})
		// return
	}

	currentCoords := session.CurrentCoordinates

	if nil == currentCoords {
		currentCoords = &session.StartingCoordinates
	}

	curntLocation, err := model.CreateLocationFromCoordinates(currentCoords.Latitude, currentCoords.Longitude)

	location, err := getOrderPrimaryLocation(order)

	output := dto.CreateOrderOutput(order)

	var entry = dto.RespondentOrderEntryIOutput{
		Order: output,
	}

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": err.Error()})
		return
	}

	routingInfo, err := locationService.GetDistance(curntLocation, location)

	if nil != err {

	} else {
		entry.Routing = *routingInfo
	}

	c.JSON(http.StatusOK, gin.H{"data": entry, "status": "success"})
}

func FindAvailableOrdersForRespondent(c *gin.Context) {

	respondentRepo := repository.GetRespondentRepository()
	sessionRepo := repository.GetRespondentSessionRepository()
	locationService := service.GetLocationService()

	var requests []model.Request
	var orders []model.Order
	var result []dto.RespondentOrderEntryIOutput
	var page pagination.Page
	if err := c.ShouldBindQuery(&page); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": err.Error()})
		return
	}

	rUser, err := auth.GetCurrentUser(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": err.Error()})
		return
	}

	respondant, err := respondentRepo.GetByUser(*rUser)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": err.Error()})
		return
	}

	session, err := sessionRepo.GetActiveByRespondent(*respondant, "Assignments.Assignment.Product")
	if nil != err {
		message := fmt.Sprintf("Could not find active session for responder[id:%s]", respondant.ID)
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": message})
		return
	}

	query := model.DB.
		Preload("Order.Fulfilment.Responder").
		Preload("Order.Items").Preload("Order.Items.Product").
		// Preload("Order.OrderMotoristRequestSituations.Situation").
		Preload("Order.User").Preload("Order.Items.Origin").Preload("Order.Items.Destination").
		Preload("FuelTypeInfo").Preload("VehicleInfo")
	// Preload("Origin").Preload("Destination")

	//THE ASSIGNMENTS THIS RESPONDENT's SESSION IS CURRENTLY RESPONDING TO

	assignments := session.Assignments

	if assignments != nil {
		var productIds []string
		for _, elAssignment := range assignments {
			if elAssignment.Assignment == nil || elAssignment.Assignment.Product == nil {
				continue
			}

			product := elAssignment.Assignment.Product
			productIds = append(productIds, product.ID.String())
		}
		query.Where("product_id IN ?", productIds)
	}

	query = query.
		Joins("JOIN orders ON orders.id = requests.order_id", 1).
		Joins("LEFT JOIN order_fulfilments ON order_fulfilments.id = orders.fulfilment_id AND 1 = ?", 1)
	// // Joins("LEFT JOIN respondent_sessions ON order_fulfilments.session_id = respondent_sessions.id", 1).
	// Where("order_fulfilments.id IS NULL OR ((order_fulfilments.responder_id = ? /*OR respondent_sessions.respondent_id = ?*/) AND ((order_fulfilments.client_confirmed_at IS NULL) AND (order_fulfilments.auto_confirmed_at IS NULL)))", respondant.ID)

	statuses := page.Status
	status := "all"
	if len(statuses) > 0 {
		status = statuses[0]
	}

	switch status {
	case "open":
		query = query.Where("orders.status NOT IN ? AND order_fulfilments.id IS NULL", []string{model.ORDER_EARNING_STATUS_COMPLETED, model.ORDER_EARNING_STATUS_CANCELLED})
		break
	case "assigned":
		query = query.Where("order_fulfilments.responder_id = ? AND ((order_fulfilments.client_confirmed_at IS NULL) AND (order_fulfilments.auto_confirmed_at IS NULL) AND (order_fulfilments.responder_confirmed_at IS NULL))", respondant.ID)
		break
	case "completed":
		query = query.Where("order_fulfilments.responder_id = ? AND ((order_fulfilments.client_confirmed_at IS NOT NULL) OR (order_fulfilments.auto_confirmed_at IS NOT NULL)  OR (order_fulfilments.responder_confirmed_at IS NOT NULL))", respondant.ID)
		break
	default:
		query = query.Where("order_fulfilments.id IS NULL OR ((order_fulfilments.responder_id = ? /*OR respondent_sessions.respondent_id = ?*/) )", respondant.ID)
		break
	}
	// RESET STATUS
	page.Status = []string{}

	fmt.Println(query.Statement.SQL.String())
	if err := query.Scopes(pagination.Paginate(requests, &page, query)).Find(&requests).Error; nil != err {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": err.Error()})
		return
	}

	currentCoords := session.CurrentCoordinates

	if nil == currentCoords {
		currentCoords = &session.StartingCoordinates
	}

	curntLocation, err := model.CreateLocationFromCoordinates(currentCoords.Latitude, currentCoords.Longitude)

	for _, request := range requests {
		// orders = append(orders, *request.Order)

		order := request.Order
		location, err := getOrderPrimaryLocation(*order)

		var entry = dto.RespondentOrderEntryIOutput{
			Order: dto.CreateOrderOutput(*order),
		}

		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": err.Error()})
			return
		}

		routingInfo, err := locationService.GetDistance(curntLocation, location)

		if nil != err {

		} else {
			entry.Routing = *routingInfo
		}

		result = append(result, entry)

	}
	if orders == nil {
		orders = make([]model.Order, 0)
	}
	if result == nil {
		result = make([]dto.RespondentOrderEntryIOutput, 0)
	}

	page.Rows = result
	c.JSON(http.StatusOK, gin.H{"data": page})
}

func getOrderPrimaryLocation(order model.Order) (location *model.Location, err error) {

	if len(*order.Items) < 1 {
		return nil, errors.New("Order must have at least one Item")
	}
	item := (*order.Items)[0]

	if item.Origin != nil {
		return item.Origin, nil
	}

	if item.Destination != nil {
		return item.Destination, nil
	}

	return nil, errors.New("Order Item must have destination")

}

func getOrderById(id uuid.UUID) (order *model.Order, err error) {

	query := model.DB.Where("id = ?", id).
		Preload("Fulfilment.Responder.User").
		Preload("User").Preload("OrderMotoristRequestSituations.Situation").
		Preload("Payment").Preload("Items").
		Preload("Adjustments").Preload("Items.Product").
		Preload("Items.Origin").Preload("Items.Destination").
		Preload("Items.FuelTypeInfo").Preload("Items.VehicleInfo")

	if err = query.First(order).Error; err != nil {
		return nil, err
	}

	return order, nil
}
