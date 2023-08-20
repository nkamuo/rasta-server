package controller

import (
	"fmt"
	"net/http"

	"github.com/nkamuo/rasta-server/data/pagination"
	"github.com/nkamuo/rasta-server/model"
	"github.com/nkamuo/rasta-server/repository"
	"github.com/nkamuo/rasta-server/utils/auth"

	"github.com/gin-gonic/gin"
)

func FindAvailableOrdersForRespondent(c *gin.Context) {

	respondentRepo := repository.GetRespondentRepository()
	sessionRepo := repository.GetRespondentSessionRepository()

	var requests []model.Request
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
		Preload("Order")
		// .Preload("Order.Items").
		// Preload("Order.User").Preload("Order.Items.Origin").Preload("Order.Items.Destination").
		// Preload("FuelTypeInfo").Preload("VehicleInfo")
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

	for _, request := range requests {
		orders = append(orders, *request.Order)
	}
	if orders == nil {
		orders = make([]model.Order, 0)
	}

	page.Rows = requests
	c.JSON(http.StatusOK, gin.H{"data": page})
}
