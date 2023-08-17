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
	"github.com/nkamuo/rasta-server/service"

	"github.com/gin-gonic/gin"
	// "github.com/gin-gonic/gin"
)

func FindLocations(c *gin.Context) {
	var locations []model.Location
	var page pagination.Page
	if err := c.ShouldBindQuery(&page); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": err.Error()})
		return
	}

	query := model.DB.Model(&model.Location{})

	if page.Search != "" {
		nameSearchQuery := strings.Join([]string{"%", page.Search, "%"}, "")
		query = query.Where("name LIKE ?", nameSearchQuery)
	}

	query = query.Scopes(pagination.Paginate(locations, &page, query))

	if err := query.Find(&locations).Error; nil != err {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": err.Error()})
		return
	}
	page.Rows = locations
	c.JSON(http.StatusOK, gin.H{"data": page})
}

func ResolveDistanc(c *gin.Context) {

	locationService := service.GetLocationService()

	var input dto.DistanceMatrixRequestInput
	if err := c.ShouldBindQuery(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": err.Error()})
		return
	}

	origin, err := locationService.Search(input.Origin)
	if nil != err {
		message := fmt.Sprintf("Could not resolve origin location[%s]:%s", input.Origin, err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": message})
		return
	}
	destination, err := locationService.Search(input.Destination)
	if nil != err {
		message := fmt.Sprintf("Could not resolve destination location[%s]:%s", input.Destination, err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": message})
		return
	}

	response, err := locationService.GetDistance(origin, destination)
	if nil != err {
		message := fmt.Sprintf("Could not resolve distance: %s", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": message})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "success", "data": response})
}

func ResolveDistanceMatrix(c *gin.Context) {

	locationService := service.GetLocationService()

	var input dto.DistanceMatrixRequestInput
	if err := c.ShouldBindQuery(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": err.Error()})
		return
	}

	origin, err := locationService.Search(input.Origin)
	if nil != err {
		message := fmt.Sprintf("Could not resolve origin location[%s]:%s", input.Origin, err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": message})
		return
	}
	destination, err := locationService.Search(input.Destination)
	if nil != err {
		message := fmt.Sprintf("Could not resolve destination location[%s]:%s", input.Destination, err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": message})
		return
	}

	response, err := locationService.ResolveDistanceMatrix(origin, destination)
	if nil != err {
		message := fmt.Sprintf("Could not resolve distance matrix:%s", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": message})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "success", "data": response})
}

func FindLocation(c *gin.Context) {
	locationService := service.GetLocationService()

	id, err := uuid.Parse(c.Param("id"))
	if nil != err {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "Invalid Id provided"})
		return
	}
	location, err := locationService.GetById(id)
	if nil != err {
		message := fmt.Sprintf("Could not find location with [id:%s]", id)
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": message})
	}
	c.JSON(http.StatusOK, gin.H{"status": "success", "data": location})
}

// func UpdateLocation(c *gin.Context) {
// 	locationService := service.GetLocationService()

// 	var input dto.LocationUpdateInput

// 	if err := c.Copy().BindJSON(&input); err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": err.Error()})
// 		return
// 	}

// 	id, err := uuid.Parse(c.Param("id"))
// 	if nil != err {
// 		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "Invalid Id provided"})
// 		return
// 	}
// 	location, err := locationService.GetById(id)
// 	if nil != err {
// 		message := fmt.Sprintf("Could not find location with [id:%s]", id)
// 		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": message})
// 	}

// 	if nil != input.Coordinates {
// 		location.Coordinates = *input.Coordinates
// 	}

// 	if nil != input.GoogleID {
// 		location.GoggleID = *input.GoogleID
// 	}

// 	if nil != input.Code {
// 		location.Code = *input.Code
// 	}
// 	if nil != input.Name {
// 		location.Name = *input.Name
// 	}
// 	if nil != input.ShortName {
// 		location.ShortName = *input.ShortName
// 	}
// 	if nil != input.LongName {
// 		location.LongName = *input.LongName
// 	}
// 	if nil != input.Description {
// 		location.Description = *input.Description
// 	}
// 	if nil != input.Active {
// 		location.Active = input.Active
// 	}
// 	if nil != input.Category {
// 		if err := ValidateLocationCategory(*input.Category); err != nil {
// 			c.JSON(http.StatusBadRequest, gin.H{"message": err.Error(), "status": "error"})
// 			return
// 		}
// 		location.Category = *input.Category
// 	}

// 	if err := locationService.Save(location); nil != err {
// 		c.JSON(http.StatusOK, gin.H{"status": "error", "message": err.Error()})
// 		return
// 	}

// 	message := fmt.Sprintf("Updated location \"%s\"", location.ID)
// 	c.JSON(http.StatusOK, gin.H{"data": location, "status": "success", "message": message})
// }

func DeleteLocation(c *gin.Context) {
	locationService := service.GetLocationService()

	id, err := uuid.Parse(c.Param("id"))
	if nil != err {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "Invalid Id provided"})
		return
	}
	location, err := locationService.GetById(id)
	if nil != err {
		message := fmt.Sprintf("Could not find location with [id:%s]", id)
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": message})
	}

	if err := locationService.Delete(location); nil != err {
		c.JSON(http.StatusOK, gin.H{"status": "error", "message": err.Error()})
		return
	}

	message := fmt.Sprintf("Deleted location \"%s\"", location.ID)
	c.JSON(http.StatusOK, gin.H{"data": location, "status": "success", "message": message})
}

func ValidateLocationCategory(category model.ProductCategory) (err error) {
	switch category {
	case model.PLACE_CITY:
		return nil
	case model.PLACE_STATE:
		return nil
	case model.PLACE_COUNTRY:
		return nil
	}
	return errors.New(fmt.Sprintf("Unsupported Location Category \"%s\"", category))
}
