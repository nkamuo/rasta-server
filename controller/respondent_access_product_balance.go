package controller

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/google/uuid"
	// "github.com/mitchellh/mapstructure"
	"github.com/nkamuo/rasta-server/data/pagination"
	"github.com/nkamuo/rasta-server/model"
	"github.com/nkamuo/rasta-server/repository"
	"github.com/nkamuo/rasta-server/service"
	"github.com/nkamuo/rasta-server/utils/auth"

	"github.com/gin-gonic/gin"
	// "github.com/gin-gonic/gin"
)

func FindRespondentAccessProductBalances(c *gin.Context) {

	respondentRepo := repository.GetRespondentRepository()
	placeRepo := repository.GetPlaceRepository()

	var balances []model.RespondentAccessProductBalance
	var page pagination.Page
	if err := c.ShouldBindQuery(&page); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": err.Error()})
		return
	}

	query := model.DB.Preload("Respondent.User").Preload("Assignments").Preload("Assignments.Assignment.Product") //.Preload("Place")

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
			message := fmt.Sprintf("Could not find referenced Responder[%s]: %s", respondentID, err.Error())
			c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": message})
			return
		}
		query = query.Where("respondent_id = ?", respondentID)
	}

	if status := c.Query("status"); status != "" {
		query = query.Where("active = ?", true)
	}

	if err := query.Scopes(pagination.Paginate(balances, &page, query)).Find(&balances).Error; nil != err {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": err.Error()})
		return
	}
	page.Rows = balances
	c.JSON(http.StatusOK, gin.H{"data": page})
}

func CreateRespondentAccessProductBalance(c *gin.Context) {

	// userService := service.GetUserService()
	// placeService := service.GetPlaceService()
	// productService := service.GetProductService()
	// respondentService := service.GetRespondentService()
	// balanceService := service.GetRespondentAccessProductBalanceService()
	// balanceRepo := repository.GetRespondentAccessProductBalanceRepository()
	// assignmentService := service.GetProductRespondentAssignmentService()

	// rUser, err := auth.GetCurrentUser(c)
	// if err != nil {
	// 	c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": err.Error()})
	// 	return
	// }

	// var input dto.RespondentAccessProductBalanceCreationInput
	// if err := c.ShouldBindJSON(&input); err != nil {
	// 	c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": err.Error()})
	// 	return
	// }

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

	// respondant, err := respondentService.GetById(input.RespondentID)
	// if nil != err {
	// 	c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": err.Error()})
	// 	return
	// }

	// user, err := userService.GetById(*respondant.UserID)
	// if nil != err {
	// 	c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": err.Error()})
	// 	return
	// }

	// var crntAssignment model.RespondentAccessProductBalance
	// var existingCount int64
	// err = model.DB.Where("respondent_id = ? AND active = ?", respondant.ID, true).Model(&crntAssignment).Count(&existingCount).Error
	// if nil != err {
	// 	c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": err.Error()})
	// 	return

	// }

	// if existingCount > 0 {
	// 	message := fmt.Sprintf(
	// 		"There is already an active balance for \"%v\"",
	// 		user.FullName(),
	// 	)
	// 	c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": message})
	// 	return

	// }

	// balance, err := balanceRepo.GetActiveByRespondent(*respondant)
	// if err != nil {
	// 	if err.Error() != "record not found" {
	// 		message := fmt.Sprintf(
	// 			"There was an error searching active balance for \"%v\"",
	// 			user.FullName(),
	// 		)
	// 		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": message})
	// 		return
	// 	}
	// } else {

	// 	if err := balanceService.Save(balance); err != nil {
	// 		message := fmt.Sprintf("An error occured: %s", err.Error())
	// 		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": message})
	// 		return
	// 	}

	// 	c.JSON(http.StatusOK, gin.H{"status": "success", "data": balance})
	// 	return

	// }
	// fmt.Printf("Input USer ID: %s\n user.ID: %s\n balance.UserId: %s\n", input.UserId, user.ID, balance.UserID)

	// if err := balanceService.Save(balance); nil != err {
	// 	c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": err.Error()})
	// 	return
	// }
	// c.JSON(http.StatusOK, gin.H{"status": "success", "data": balance})
}

// ////////
// // FIND SESSION
// //
// /////////////
func FindRespondentAccessProductBalance(c *gin.Context) {
	balanceService := service.GetRespondentAccessProductBalanceService()

	rUser, err := auth.GetCurrentUser(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": err.Error()})
		return
	}

	id, err := uuid.Parse(c.Param("id"))
	if nil != err {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "Invalid Id provided"})
		return
	}
	// Preload("Respondent.User").Preload("Assignments").Preload("Assignments.Assignment.Product")
	balance, err := balanceService.GetById(id, "Respondent.User", "Assignments.Assignment.Product")
	if nil != err {
		message := fmt.Sprintf("Could not find balance with [id:%s]", id)
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": message})
		return
	}

	if *rUser.IsAdmin {
		/// USER IS ADMIN - ADMIN CAN VIEW ANY SESSION
	} else {

		respondant, err := auth.GetCurrentRespondent(c)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": err.Error()})
			return
		}

		if respondant.ID.String() != (*balance.RespondentID).String() {
			message := fmt.Sprintf("Could not find your balance with [id:%s]", id)
			c.JSON(http.StatusForbidden, gin.H{"status": "error", "message": message})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{"status": "success", "data": balance})
}

func FindCurrentRespondentAccessProductBalance(c *gin.Context) {
	balanceRepo := repository.GetRespondentAccessProductBalanceRepository()

	respondant, err := auth.GetCurrentRespondent(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": err.Error()})
		return
	}

	balance, err := balanceRepo.GetActiveByRespondent(*respondant, "Assignments", "Assignments.Assignment", "Assignments.Assignment.Product")
	if nil != err {
		message := fmt.Sprintf("Could not find active balance for respondent[id:%s]", respondant.ID)
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": message})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "success", "data": balance})
}

func DeleteRespondentAccessProductBalance(c *gin.Context) {
	balanceService := service.GetRespondentAccessProductBalanceService()

	id, err := uuid.Parse(c.Param("id"))
	if nil != err {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "Invalid Id provided"})
		return
	}
	balance, err := balanceService.GetById(id)
	if nil != err {
		message := fmt.Sprintf("Could not find balance with [id:%s]", id)
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": message})
	}

	if err := balanceService.Delete(balance); nil != err {
		c.JSON(http.StatusOK, gin.H{"status": "error", "message": err.Error()})
		return
	}

	message := fmt.Sprintf("Deleted balance \"%s\"", balance.ID)
	c.JSON(http.StatusOK, gin.H{"data": balance, "status": "success", "message": message})
}

func ValidateRespondentAccessProductBalanceCategory(category model.ProductCategory) (err error) {
	switch category {
	case model.PLACE_CITY:
		return nil
	case model.PLACE_STATE:
		return nil
	case model.PLACE_COUNTRY:
		return nil
	}
	return errors.New(fmt.Sprintf("Unsupported RespondentAccessProductBalance Category \"%s\"", category))
}
