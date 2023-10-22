package controller

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/google/uuid"
	// "github.com/mitchellh/mapstructure"
	"github.com/nkamuo/rasta-server/data/pagination"
	"github.com/nkamuo/rasta-server/dto"
	"github.com/nkamuo/rasta-server/model"
	"github.com/nkamuo/rasta-server/repository"
	"github.com/nkamuo/rasta-server/service"
	"github.com/nkamuo/rasta-server/utils/auth"

	"github.com/gin-gonic/gin"
	// "github.com/gin-gonic/gin"
)

func FindRespondentSessionLocationEntries(c *gin.Context) {

	// respondentRepo := repository.GetRespondentRepository()
	sessionRepo := repository.GetRespondentSessionRepository()

	var locationEntries []model.RespondentSessionLocationEntry
	var page pagination.Page
	if err := c.ShouldBindQuery(&page); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": err.Error()})
		return
	}

	query := model.DB.Preload("Session.Respondent") //.Preload("Product") //.Preload("Place")

	if session_id := c.Query("session_id"); session_id != "" {
		sessionID, err := uuid.Parse(session_id)
		if nil != err {
			message := fmt.Sprintf("Error parsing session_id[%s] query: %s", session_id, err.Error())
			c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": message})
			return
		}
		if _, err := sessionRepo.GetById(sessionID); err != nil {
			message := fmt.Sprintf("Could not find referenced place[%s]: %s", sessionID, err.Error())
			c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": message})
			return
		}
		query = query.Where("session_id = ?", sessionID)
	}
	/*
		if repondent_id := c.Query("repondent_id"); repondent_id != "" {
			respondentID, err := uuid.Parse(repondent_id)
			if nil != err {
				message := fmt.Sprintf("Error parsing repondent_id[%s] query: %s", repondent_id, err.Error())
				c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": message})
				return
			}
			if _, err := respondentRepo.GetById(respondentID); err != nil {
				message := fmt.Sprintf("Could not find referenced Respondent[%s]: %s", respondentID, err.Error())
				c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": message})
				return
			}
			query = query.Where("repondent_id = ?", respondentID)
		}
	*/
	if err := query.Scopes(pagination.Paginate(locationEntries, &page, query)).Find(&locationEntries).Error; nil != err {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": err.Error()})
		return
	}
	page.Rows = locationEntries
	c.JSON(http.StatusOK, gin.H{"data": page})
}

func CreateRespondentSessionLocationEntry(c *gin.Context) {

	// userService := service.GetUserService()
	// placeService := service.GetPlaceService()
	// respondentService := service.GetProductService()
	sessionService := service.GetRespondentSessionService()
	respondentRepo := repository.GetRespondentRepository()
	sessionRepo := repository.GetRespondentSessionRepository()
	locationEntryService := service.GetRespondentSessionLocationEntryService()

	var input dto.RespondentSessionLocationEntryCreationInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": err.Error()})
		return
	}

	rUser, err := auth.GetCurrentUser(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": err.Error()})
		return
	}

	respondent, err := respondentRepo.GetByUser(*rUser)
	if nil != err {
		message := fmt.Sprintf("Could not resolve the matching responder account: %s", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"message": message, "status": "error"})
		return
	}

	session, err := sessionRepo.GetActiveByRespondent(*respondent)
	if nil != err {
		message := fmt.Sprintf("Could not resolve active session for this responder: %s", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"message": message, "status": "error"})
		return
	}

	locationEntry := model.RespondentSessionLocationEntry{
		SessionID: &session.ID,
		Coordinates: model.LocationCoordinates{
			Latitude:  input.Coordinates.Latitude,
			Longitude: input.Coordinates.Longitude,
		},
	}

	// fmt.Printf("Input USer ID: %s\n user.ID: %s\n locationEntry.UserId: %s\n", input.UserId, user.ID, locationEntry.UserID)

	if err := locationEntryService.Save(&locationEntry); nil != err {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": err.Error()})
		return
	}

	if err := sessionService.UpdateLocationEntry(session, locationEntry); nil != err {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": err.Error()})
		return
	}

	// ()

	c.JSON(http.StatusOK, gin.H{"status": "success", "data": locationEntry})
}

func FindRespondentSessionLocationEntry(c *gin.Context) {
	locationEntryService := service.GetRespondentSessionLocationEntryService()

	id, err := uuid.Parse(c.Param("id"))
	if nil != err {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "Invalid Id provided"})
		return
	}
	locationEntry, err := locationEntryService.GetById(id)
	if nil != err {
		message := fmt.Sprintf("Could not find locationEntry with [id:%s]", id)
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": message})
	}
	c.JSON(http.StatusOK, gin.H{"status": "success", "data": locationEntry})
}

func DeleteRespondentSessionLocationEntry(c *gin.Context) {
	locationEntryService := service.GetRespondentSessionLocationEntryService()

	id, err := uuid.Parse(c.Param("id"))
	if nil != err {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "Invalid Id provided"})
		return
	}
	locationEntry, err := locationEntryService.GetById(id)
	if nil != err {
		message := fmt.Sprintf("Could not find locationEntry with [id:%s]", id)
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": message})
	}

	if err := locationEntryService.Delete(locationEntry); nil != err {
		c.JSON(http.StatusOK, gin.H{"status": "error", "message": err.Error()})
		return
	}

	message := fmt.Sprintf("Deleted locationEntry \"%s\"", locationEntry.ID)
	c.JSON(http.StatusOK, gin.H{"data": locationEntry, "status": "success", "message": message})
}

func ValidateRespondentSessionLocationEntryCategory(category model.ProductCategory) (err error) {
	switch category {
	case model.PLACE_CITY:
		return nil
	case model.PLACE_STATE:
		return nil
	case model.PLACE_COUNTRY:
		return nil
	}
	return errors.New(fmt.Sprintf("Unsupported RespondentSessionLocationEntry Category \"%s\"", category))
}
