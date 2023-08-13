package controller

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/google/uuid"

	// "gorm.io/gorm"

	// "github.com/mitchellh/mapstructure"
	"github.com/nkamuo/rasta-server/data/pagination"
	"github.com/nkamuo/rasta-server/dto"
	"github.com/nkamuo/rasta-server/model"
	"github.com/nkamuo/rasta-server/repository"
	"github.com/nkamuo/rasta-server/service"

	"github.com/gin-gonic/gin"
	// "github.com/gin-gonic/gin"
)

func FindPlaces(c *gin.Context) {
	var places []model.Place
	var page pagination.Page
	if err := c.ShouldBindQuery(&page); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": err.Error()})
		return
	}

	query := model.DB.Model(&model.Place{})

	if page.Search != "" {
		nameSearchQuery := strings.Join([]string{"%", page.Search, "%"}, "")
		query = query.Where("name LIKE ?", nameSearchQuery)
	}

	query = query.Scopes(pagination.Paginate(places, &page, query))

	if err := query.Find(&places).Error; nil != err {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": err.Error()})
		return
	}
	page.Rows = places
	c.JSON(http.StatusOK, gin.H{"data": page})
}

func CreatePlace(c *gin.Context) {

	placeService := service.GetPlaceService()

	var input dto.PlaceCreationInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": err.Error()})
		return
	}

	if err := ValidatePlaceCategory(input.Category); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error(), "status": "error"})
		return
	}

	// Create place
	place := model.Place{
		Active:      &input.Active,
		Code:        input.Code,
		Name:        input.Name,
		ShortName:   input.ShortName,
		LongName:    input.LongName,
		Description: input.Description,
		Coordinates: input.Coordinates,
		Category:    input.Category,
	}

	if nil != input.GoogleID {
		place.GoggleID = *input.GoogleID
	}

	// fmt.Printf("Input USer ID: %s\n user.ID: %s\n place.UserId: %s\n", input.UserId, user.ID, place.UserID)

	if err := placeService.Save(&place); nil != err {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "success", "data": place})
}

func FindPlace(c *gin.Context) {
	placeService := service.GetPlaceService()

	id, err := uuid.Parse(c.Param("id"))
	if nil != err {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "Invalid Id provided"})
		return
	}
	place, err := placeService.GetById(id)
	if nil != err {
		message := fmt.Sprintf("Could not find place with [id:%s]", id)
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": message})
	}
	c.JSON(http.StatusOK, gin.H{"status": "success", "data": place})
}

func FindPlaceByLocation(c *gin.Context) {
	placeRepo := repository.GetPlaceRepository()
	locationService := service.GetLocationService()

	location_ref := c.Query("location")

	if location_ref == "" {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "location query parameter not provided"})
		return
	}

	location, err := locationService.Search(location_ref)
	if nil != err {
		message := fmt.Sprintf("Could not resolve the specified location[%s]", location_ref)
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": message})
		return
	}

	place, err := placeRepo.GetByLocation(location)
	if nil != err {
		message := fmt.Sprintf("Could not find place with for location[address:%s]", location.Address)
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": message})
	}
	c.JSON(http.StatusOK, gin.H{"status": "success", "data": place})
}

func UpdatePlace(c *gin.Context) {
	placeService := service.GetPlaceService()

	var input dto.PlaceUpdateInput

	if err := c.Copy().BindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": err.Error()})
		return
	}

	id, err := uuid.Parse(c.Param("id"))
	if nil != err {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "Invalid Id provided"})
		return
	}
	place, err := placeService.GetById(id)
	if nil != err {
		message := fmt.Sprintf("Could not find place with [id:%s]", id)
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": message})
	}

	if nil != input.Coordinates {
		place.Coordinates = *input.Coordinates
	}

	if nil != input.GoogleID {
		place.GoggleID = *input.GoogleID
	}

	if nil != input.Code {
		place.Code = *input.Code
	}
	if nil != input.Name {
		place.Name = *input.Name
	}
	if nil != input.ShortName {
		place.ShortName = *input.ShortName
	}
	if nil != input.LongName {
		place.LongName = *input.LongName
	}
	if nil != input.Description {
		place.Description = *input.Description
	}
	if nil != input.Active {
		place.Active = input.Active
	}
	if nil != input.Category {
		if err := ValidatePlaceCategory(*input.Category); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": err.Error(), "status": "error"})
			return
		}
		place.Category = *input.Category
	}

	if err := placeService.Save(place); nil != err {
		c.JSON(http.StatusOK, gin.H{"status": "error", "message": err.Error()})
		return
	}

	message := fmt.Sprintf("Updated place \"%s\"", place.ID)
	c.JSON(http.StatusOK, gin.H{"data": place, "status": "success", "message": message})
}

func DeletePlace(c *gin.Context) {
	placeService := service.GetPlaceService()

	id, err := uuid.Parse(c.Param("id"))
	if nil != err {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "Invalid Id provided"})
		return
	}
	place, err := placeService.GetById(id)
	if nil != err {
		message := fmt.Sprintf("Could not find place with [id:%s]", id)
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": message})
	}

	if err := placeService.Delete(place); nil != err {
		c.JSON(http.StatusOK, gin.H{"status": "error", "message": err.Error()})
		return
	}

	message := fmt.Sprintf("Deleted place \"%s\"", place.ID)
	c.JSON(http.StatusOK, gin.H{"data": place, "status": "success", "message": message})
}

func ValidatePlaceCategory(category model.ProductCategory) (err error) {
	switch category {
	case model.PLACE_CITY:
		return nil
	case model.PLACE_STATE:
		return nil
	case model.PLACE_COUNTRY:
		return nil
	}
	return errors.New(fmt.Sprintf("Unsupported Place Category \"%s\"", category))
}
