package controller

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/nkamuo/rasta-server/data/pagination"
	"github.com/nkamuo/rasta-server/dto"
	"github.com/nkamuo/rasta-server/model"
	"github.com/nkamuo/rasta-server/repository"
	"github.com/nkamuo/rasta-server/service"
	"github.com/nkamuo/rasta-server/utils/auth"

	"github.com/gin-gonic/gin"
)

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

	requestingUser, err := auth.GetCurrentUser(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": err.Error()})
		return
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

	query := model.DB.
		Preload("Order").
		Preload("Order.Items").Preload("Order.Items.Product").
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
			Order: *order,
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
