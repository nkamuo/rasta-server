package controller

import (
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

func FindRespondentServiceReviews(c *gin.Context) {

	placeRepo := repository.GetPlaceRepository()

	var respondentServiceReviews []model.RespondentServiceReview
	var page pagination.Page
	if err := c.ShouldBindQuery(&page); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": err.Error()})
		return
	}
	query := model.DB.Preload("Respondent").Preload("Request")

	if respondent_id := c.Query("respondent_id"); respondent_id != "" {
		respondentID, err := uuid.Parse(respondent_id)
		if nil != err {
			message := fmt.Sprintf("Error parsing respondent_id[%s] query: %s", respondent_id, err.Error())
			c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": message})
			return
		}
		if _, err := placeRepo.GetById(respondentID); err != nil {
			message := fmt.Sprintf("Could not find referenced place[%s]: %s", respondentID, err.Error())
			c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": message})
			return
		}
		query = query.Where("respondent_id = ?", respondentID)
	}

	if err := query.Scopes(pagination.Paginate(respondentServiceReviews, &page, query)).Find(&respondentServiceReviews).Error; nil != err {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": err.Error()})
		return
	}
	page.Rows = respondentServiceReviews
	c.JSON(http.StatusOK, gin.H{"data": page})
}

func CreateRespondentServiceReview(c *gin.Context) {

	userService := service.GetUserService()
	orderService := service.GetOrderService()
	requestService := service.GetRequestService()
	respondentServiceReviewService := service.GetRespondentServiceReviewService()

	var author *model.User

	var input dto.RespondentServiceReviewCreationInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := auth.GetCurrentUser(c)
	if err != nil {
		message := fmt.Sprintf("Could not resolve requeting user account: %s", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": message})
		return
	}

	if nil != input.AuthorID {
		author, err = userService.GetById(*input.AuthorID)
		if err != nil {
			message := fmt.Sprintf("Could not find author with [id:%s]: %s", input.AuthorID, err.Error())
			c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": message})
			return
		}
	} else {
		author = user
	}

	request, err := requestService.GetById(input.RequestID)
	if err != nil {
		message := fmt.Sprintf("Could not find Request with [id:%s]: %s", input.RequestID, err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": message})
		return
	}

	order, err := orderService.GetById(*request.OrderID)
	if err != nil {
		message := fmt.Sprintf("Interner Server Error: %s", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": message})
		return
	}

	if order.UserID != input.AuthorID {
		message := fmt.Sprintf("The author of a review must be the user associated witht the order %s", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": message})
		return
	}

	var respondentServiceReview *model.RespondentServiceReview

	if err := model.DB.
		Where("request_id = ? AND author_id = ?", request.ID, author.ID).
		First(&respondentServiceReview).Error; nil != err {
		if err.Error() != "record not found" {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
	}

	if nil == respondentServiceReview {
		respondentServiceReview = &model.RespondentServiceReview{}
	}

	// respondentServiceReview.RequestID = &input.RequestID
	// respondentServiceReview.AuthorID = input.AuthorID
	respondentServiceReview.OrderID = &order.ID
	respondentServiceReview.Description = input.Description
	respondentServiceReview.Published = input.Published
	respondentServiceReview.Rating = input.Rating

	// Description: input.Description,
	// Rate:        input.Rate,
	// Active: &input.Active,

	if err := respondentServiceReviewService.Save(respondentServiceReview); nil != err {
		c.JSON(http.StatusOK, gin.H{"status": "error", "message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": respondentServiceReview, "status": "success"})
}

func FindRespondentServiceReview(c *gin.Context) {
	respondentServiceReviewService := service.GetRespondentServiceReviewService()

	id, err := uuid.Parse(c.Param("id"))
	if nil != err {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "Invalid Id provided"})
		return
	}

	respondentServiceReview, err := respondentServiceReviewService.GetById(id)
	if err != nil {
		message := fmt.Sprintf("Could not find respondentServiceReview with [id:%s]", id)
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": message})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": respondentServiceReview})
}

func UpdateRespondentServiceReview(c *gin.Context) {
	respondentServiceReviewService := service.GetRespondentServiceReviewService()

	var input dto.RespondentServiceReviewUpdateInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	id, err := uuid.Parse(c.Param("id"))
	if nil != err {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "Invalid Id provided"})
		return
	}
	respondentServiceReview, err := respondentServiceReviewService.GetById(id)
	if nil != err {
		message := fmt.Sprintf("Could not find respondentServiceReview with [id:%s]", id)
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": message})
	}

	if nil != input.Rating {
		respondentServiceReview.Rating = *input.Rating
	}
	if nil != input.Published {
		respondentServiceReview.Published = *&input.Published
	}
	if nil != input.Description {
		respondentServiceReview.Description = input.Description
	}

	if err := respondentServiceReviewService.Save(respondentServiceReview); nil != err {
		c.JSON(http.StatusOK, gin.H{"status": "error", "message": err.Error()})
		return
	}

	message := fmt.Sprintf("Updated model \"%s\"", respondentServiceReview.ID)
	c.JSON(http.StatusOK, gin.H{"data": respondentServiceReview, "status": "success", "message": message})
}

func DeleteRespondentServiceReview(c *gin.Context) {
	id := c.Param("id")

	var respondentServiceReview model.RespondentServiceReview

	if err := model.DB.Where("id = ?", id).First(&respondentServiceReview).Error; err != nil {
		message := fmt.Sprintf("Could not find respondentServiceReview with [id:%s]", id)
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": message})
		return
	}
	model.DB.Delete(&respondentServiceReview)
	message := fmt.Sprintf("Deleted respondentServiceReview \"%s\"", respondentServiceReview.ID)
	c.JSON(http.StatusOK, gin.H{"data": respondentServiceReview, "status": "success", "message": message})
}
