package controller

import (
	"errors"
	"fmt"
	"net/http"
	"time"

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

func FindRespondentSessions(c *gin.Context) {

	respondentRepo := repository.GetRespondentRepository()
	placeRepo := repository.GetPlaceRepository()

	var sessions []model.RespondentSession
	var page pagination.Page
	if err := c.ShouldBindQuery(&page); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": err.Error()})
		return
	}

	query := model.DB.Preload("Respondent").Preload("Assignments").Preload("Assignments.Assignment.Product") //.Preload("Place")

	if place_id := c.Query("place_id"); place_id != "" {
		placeID, err := uuid.Parse(place_id)
		if nil != err {
			message := fmt.Sprintf("Error parsing place_id[%s] query: %s", place_id, err.Error())
			c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": message})
			return
		}
		if _, err := placeRepo.GetById(placeID); err != nil {
			message := fmt.Sprintf("Could not find referenced place[%s]: %s", placeID, err.Error())
			c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": message})
			return
		}
		query = query.Where("place_id = ?", placeID)
	}

	if respondent_id := c.Query("respondent_id"); respondent_id != "" {
		respondentID, err := uuid.Parse(respondent_id)
		if nil != err {
			message := fmt.Sprintf("Error parsing respondent_id[%s] query: %s", respondent_id, err.Error())
			c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": message})
			return
		}
		if _, err := respondentRepo.GetById(respondentID); err != nil {
			message := fmt.Sprintf("Could not find referenced Respondent[%s]: %s", respondentID, err.Error())
			c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": message})
			return
		}
		query = query.Where("respondent_id = ?", respondentID)
	}

	if err := query.Scopes(pagination.Paginate(sessions, &page, query)).Find(&sessions).Error; nil != err {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": err.Error()})
		return
	}
	page.Rows = sessions
	c.JSON(http.StatusOK, gin.H{"data": page})
}

func CreateRespondentSession(c *gin.Context) {

	userService := service.GetUserService()
	// placeService := service.GetPlaceService()
	// productService := service.GetProductService()
	respondentService := service.GetRespondentService()
	sessionService := service.GetRespondentSessionService()
	sessionRepo := repository.GetRespondentSessionRepository()
	assignmentService := service.GetProductRespondentAssignmentService()

	// rUser, err := auth.GetCurrentUser(c)
	// if err != nil {
	// 	c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": err.Error()})
	// 	return
	// }

	var input dto.RespondentSessionCreationInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": err.Error()})
		return
	}

	// product, err := productService.GetById(input.ProductID)
	// if nil != err {
	// 	message := fmt.Sprintf("Could not resolve the specified Product with [id:%s]: %s", input.ProductID, err.Error())
	// 	c.JSON(http.StatusBadRequest, gin.H{"message": message, "status": "error"})
	// 	return
	// }

	// place, err := placeService.GetById(product.PlaceID)
	// if nil != err {
	// 	message := fmt.Sprintf("Could not resolve Place [id:%s]: %s", product.PlaceID, err.Error())
	// 	c.JSON(http.StatusBadRequest, gin.H{"message": message, "status": "error"})
	// 	return
	// }

	respondant, err := respondentService.GetById(input.RespondentID)
	if nil != err {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": err.Error()})
		return
	}

	user, err := userService.GetById(*respondant.UserID)
	if nil != err {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": err.Error()})
		return
	}

	// var crntAssignment model.RespondentSession
	// var existingCount int64
	// err = model.DB.Where("respondent_id = ? AND active = ?", respondant.ID, true).Model(&crntAssignment).Count(&existingCount).Error
	// if nil != err {
	// 	c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": err.Error()})
	// 	return

	// }

	// if existingCount > 0 {
	// 	message := fmt.Sprintf(
	// 		"There is already an active session for \"%v\"",
	// 		user.FullName(),
	// 	)
	// 	c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": message})
	// 	return

	// }

	crntSession, err := sessionRepo.GetActiveByRespondent(*respondant)
	if err != nil {
		if err.Error() != "record not found" {
			message := fmt.Sprintf(
				"There was an error searching active session for \"%v\"",
				user.FullName(),
			)
			c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": message})
			return
		}
	} else {
		if err := sessionService.Close(crntSession); err != nil {
			message := fmt.Sprintf(
				"There was an error closing active session for \"%v\"",
				user.FullName(),
			)
			c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": message})
			return
		}
	}

	var assignedProducts []model.RespondentSessionAssignedProduct

	for _, aInput := range input.Assignments {

		assignment, err := assignmentService.GetById(aInput.AssignmentID)
		if err != nil {
			message := fmt.Sprintf("Error resolving assignment with [id:%s]: %s", aInput.AssignmentID, err.Error())
			c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": message})
			return
		}

		asEntry := model.RespondentSessionAssignedProduct{
			AssignmentID: &assignment.ID,
			Active:       &aInput.Active,
			Note:         aInput.Note,
			Description:  aInput.Description,
		}
		assignedProducts = append(assignedProducts, asEntry)
	}

	var startingCoords = model.LocationCoordinates{
		Latitude:  input.StartingCoordinates.Latitude,
		Longitude: input.StartingCoordinates.Longitude,
		Altitude:  input.StartingCoordinates.Altitude,
		Accuracy:  input.StartingCoordinates.Accuracy,
	}

	now := time.Now()
	// Create session
	session := model.RespondentSession{
		RespondentID:        &respondant.ID,
		StartingCoordinates: startingCoords,
		Assignments:         assignedProducts,
		StartedAt:           &now,
		Note:                input.Note,
		Description:         input.Description,
	}

	if input.Active != nil {
		session.Active = input.Active
	} else {
		*session.Active = true
	}

	// if *rUser.IsAdmin {
	// 	session.Active = input.Active
	// } else {

	// }

	// fmt.Printf("Input USer ID: %s\n user.ID: %s\n session.UserId: %s\n", input.UserId, user.ID, session.UserID)

	if err := sessionService.Save(&session); nil != err {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "success", "data": session})
}

func FindRespondentSession(c *gin.Context) {
	sessionService := service.GetRespondentSessionService()

	id, err := uuid.Parse(c.Param("id"))
	if nil != err {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "Invalid Id provided"})
		return
	}
	session, err := sessionService.GetById(id)
	if nil != err {
		message := fmt.Sprintf("Could not find session with [id:%s]", id)
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": message})
	}
	c.JSON(http.StatusOK, gin.H{"status": "success", "data": session})
}

func FindCurrentRespondentSession(c *gin.Context) {
	respondentRepo := repository.GetRespondentRepository()
	sessionRepo := repository.GetRespondentSessionRepository()

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

	session, err := sessionRepo.GetActiveByRespondent(*respondant, "Assignments", "Assignments.Assignment", "Assignments.Assignment.Product")
	if nil != err {
		message := fmt.Sprintf("Could not find active session for respondent[id:%s]", respondant.ID)
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": message})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "success", "data": session})
}

func UpdateRespondentSession(c *gin.Context) {
	sessionService := service.GetRespondentSessionService()

	var input dto.RespondentSessionUpdateInput
	if err := c.Copy().ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": err.Error()})
		return
	}

	// var input dto.RespondentSessionUpdateInput
	// mapstructure.Decode(requestBody, input)

	// requestingUser, err := auth.GetCurrentUser(c)
	// if err != nil {
	// 	c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": err.Error()})
	// 	return
	// }

	id, err := uuid.Parse(c.Param("id"))
	if nil != err {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "Invalid Id provided"})
		return
	}
	session, err := sessionService.GetById(id)
	if nil != err {
		message := fmt.Sprintf("Could not find session with [id:%s]", id)
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": message})
		return
	}

	if input.Note != nil {
		session.Note = *input.Note
	}
	if input.Description != nil {
		session.Description = *input.Description
	}

	if input.Active != nil {
		session.Active = input.Active
	}

	if err := sessionService.Save(session); nil != err {
		c.JSON(http.StatusOK, gin.H{"status": "error", "message": err.Error()})
		return
	}

	message := fmt.Sprintf("Updated session \"%s\"", session.ID)
	c.JSON(http.StatusOK, gin.H{"data": session, "status": "success", "message": message})
}

func DeleteRespondentSession(c *gin.Context) {
	sessionService := service.GetRespondentSessionService()

	id, err := uuid.Parse(c.Param("id"))
	if nil != err {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "Invalid Id provided"})
		return
	}
	session, err := sessionService.GetById(id)
	if nil != err {
		message := fmt.Sprintf("Could not find session with [id:%s]", id)
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": message})
	}

	if err := sessionService.Delete(session); nil != err {
		c.JSON(http.StatusOK, gin.H{"status": "error", "message": err.Error()})
		return
	}

	message := fmt.Sprintf("Deleted session \"%s\"", session.ID)
	c.JSON(http.StatusOK, gin.H{"data": session, "status": "success", "message": message})
}

func ValidateRespondentSessionCategory(category model.ProductCategory) (err error) {
	switch category {
	case model.PLACE_CITY:
		return nil
	case model.PLACE_STATE:
		return nil
	case model.PLACE_COUNTRY:
		return nil
	}
	return errors.New(fmt.Sprintf("Unsupported RespondentSession Category \"%s\"", category))
}
