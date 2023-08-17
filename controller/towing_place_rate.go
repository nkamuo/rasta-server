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

	"github.com/gin-gonic/gin"
)

func FindTowingPlaceRates(c *gin.Context) {

	placeRepo := repository.GetPlaceRepository()

	var towingPlaceRates []model.TowingPlaceRate
	var page pagination.Page
	if err := c.ShouldBindQuery(&page); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": err.Error()})
		return
	}
	query := model.DB.Preload("Place") //.Preload("Towing")

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

	if err := query.Scopes(pagination.Paginate(towingPlaceRates, &page, query)).Order("max_distance ASC").Find(&towingPlaceRates).Error; nil != err {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": err.Error()})
		return
	}
	page.Rows = towingPlaceRates
	c.JSON(http.StatusOK, gin.H{"data": page})
}

func CreateTowingPlaceRate(c *gin.Context) {

	placeService := service.GetPlaceService()
	towingPlaceRateService := service.GetTowingPlaceRateService()

	var input dto.TowingPlaceRateCreationInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	place, err := placeService.GetById(input.PlaceID)
	if err != nil {
		message := fmt.Sprintf("Could not find place with [id:%s]: %s", input.PlaceID, err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": message})
		return
	}

	// if err := model.DB.
	// 	Where("fuel_type_id = ? AND place_id = ?", input.TowingID, input.PlaceID).
	// 	First(&towingPlaceRate).Error; nil != err {
	// 	if err.Error() != "record not found" {
	// 		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	// 		return
	// 	}
	// }

	towingPlaceRate := &model.TowingPlaceRate{
		PlaceID: &place.ID,
		Rate:    &input.Rate,
		//
		MinDistance: input.MinDistance,
		MaxDistance: &input.MaxDistance,
		//
		Active:      &input.Active,
		Description: input.Description,
	}

	if err := towingPlaceRateService.Save(towingPlaceRate); nil != err {
		c.JSON(http.StatusOK, gin.H{"status": "error", "message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": towingPlaceRate, "status": "success"})
}

func FindTowingRateByOriginAndDestination(c *gin.Context) {

	placeRepo := repository.GetPlaceRepository()
	locationService := service.GetLocationService()
	towingPlaceRateService := service.GetTowingPlaceRateService()

	var _query dto.TowingPlaceRateRequestQuery
	if err := c.ShouldBindQuery(&_query); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": err.Error()})
		return
	}

	origin, err := locationService.Search(_query.Origin)
	if err != nil {
		message := fmt.Sprintf("Could not find the specified origin with [id:%s]: %s", _query.Origin, err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": message})
		return
	}

	destination, err := locationService.Search(_query.Destination)
	if err != nil {
		message := fmt.Sprintf("Could not find the specified destination with [id:%s]: %s", _query.Destination, err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": message})
		return
	}

	distanceInfo, err := locationService.GetDistance(origin, destination)
	if err != nil {
		message := fmt.Sprintf("Could not resolve the distance between %s and %s : %s", *origin.Name, *destination.Name, err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": message})
		return
	}

	fmt.Println("ENSURE TO CONFIRM THAT THE ORIGIN AND DESTINATION ARE WITHIN THE SAME PLACE")

	place, err := placeRepo.GetByLocation(origin)
	if err != nil {
		message := fmt.Sprintf("Error Resolving Place for location[%s]: %s", *origin.Name, err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": message})
		return
	}

	towingPlaceRate, err := towingPlaceRateService.GetByPlaceAndDistance(*place, int64(distanceInfo.Distance.Value))
	if err != nil {
		message := fmt.Sprintf("Error Resolving rate %s", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "routing": distanceInfo, "message": message})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": towingPlaceRate, "routing": distanceInfo, "status": "success"})
}

func FindTowingRateByPlaceAndDistance(c *gin.Context) {

	placeService := service.GetPlaceService()
	towingPlaceRateService := service.GetTowingPlaceRateService()

	var _query dto.TowingPlaceRateDistanceRequestQuery
	if err := c.ShouldBindQuery(&_query); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": err.Error()})
		return
	}

	placeID, err := uuid.Parse(_query.PlaceID)

	if err != nil {
		message := fmt.Sprintf("Could not parse [place_id:%s] into UUID: %s", _query.PlaceID, err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": message})
		return
	}

	place, err := placeService.GetById(placeID)
	if err != nil {
		message := fmt.Sprintf("Could not find place with [id:%s]: %s", _query.PlaceID, err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": message})
		return
	}

	towingPlaceRate, err := towingPlaceRateService.GetByPlaceAndDistance(*place, int64(_query.Distance))
	if err != nil {
		message := fmt.Sprintf("Error Resolving rate %s", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": message})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": towingPlaceRate, "status": "success"})
}

func FindTowingPlaceRate(c *gin.Context) {
	towingPlaceRateService := service.GetTowingPlaceRateService()

	id, err := uuid.Parse(c.Param("id"))
	if nil != err {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "Invalid Id provided"})
		return
	}

	towingPlaceRate, err := towingPlaceRateService.GetById(id)
	if err != nil {
		message := fmt.Sprintf("Could not find towingPlaceRate with [id:%s]", id)
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": message})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": towingPlaceRate})
}

func UpdateTowingPlaceRate(c *gin.Context) {
	towingPlaceRateService := service.GetTowingPlaceRateService()

	var input dto.TowingPlaceRateUpdateInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	id, err := uuid.Parse(c.Param("id"))
	if nil != err {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "Invalid Id provided"})
		return
	}
	towingPlaceRate, err := towingPlaceRateService.GetById(id)
	if nil != err {
		message := fmt.Sprintf("Could not find towingPlaceRate with [id:%s]", id)
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": message})
	}

	if nil != input.Rate {
		towingPlaceRate.Rate = input.Rate
	}
	if nil != input.MinDistance {
		towingPlaceRate.MinDistance = input.MinDistance
	}
	if nil != input.MaxDistance {
		towingPlaceRate.MaxDistance = *&input.MaxDistance
	}
	if nil != input.Active {
		towingPlaceRate.Active = *&input.Active
	}
	if nil != input.Description {
		towingPlaceRate.Description = input.Description
	}

	if err := towingPlaceRateService.Save(towingPlaceRate); nil != err {
		c.JSON(http.StatusOK, gin.H{"status": "error", "message": err.Error()})
		return
	}

	message := fmt.Sprintf("Updated model \"%s\"", towingPlaceRate.ID)
	c.JSON(http.StatusOK, gin.H{"data": towingPlaceRate, "status": "success", "message": message})
}

func DeleteTowingPlaceRate(c *gin.Context) {
	id := c.Param("id")

	var towingPlaceRate model.TowingPlaceRate

	if err := model.DB.Where("id = ?", id).First(&towingPlaceRate).Error; err != nil {
		message := fmt.Sprintf("Could not find towingPlaceRate with [id:%s]", id)
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": message})
		return
	}
	model.DB.Delete(&towingPlaceRate)
	message := fmt.Sprintf("Deleted towingPlaceRate \"%s\"", towingPlaceRate.ID)
	c.JSON(http.StatusOK, gin.H{"data": towingPlaceRate, "status": "success", "message": message})
}
