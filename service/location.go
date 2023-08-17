package service

import (
	"errors"
	"fmt"
	"strings"
	"sync"

	"github.com/google/uuid"
	"github.com/nkamuo/rasta-server/initializers"
	"github.com/nkamuo/rasta-server/model"
	"github.com/nkamuo/rasta-server/repository"
	"github.com/nkamuo/rasta-server/utils/geo"
	"github.com/nkamuo/rasta-server/utils/slug"
)

var locationService LocationService
var locationRepoMutext *sync.Mutex = &sync.Mutex{}

func GetLocationService() LocationService {
	locationRepoMutext.Lock()
	if locationService == nil {
		locationService = &locationServiceImpl{
			repo:             repository.GetLocationRepository(),
			googleMapsAPIKey: initializers.CONFIG.GOOGLE_MAPS_API_KEY,
		}
	}
	locationRepoMutext.Unlock()
	return locationService
}

type LocationService interface {
	GetById(id uuid.UUID) (location *model.Location, err error)
	Resolve(data string) (location *model.Location, err error)
	Search(data string) (location *model.Location, err error)
	AssertLocationWithinPlace(location *model.Location, place *model.Place) (err error)
	ResolveDistanceMatrix(origin *model.Location, destination *model.Location) (response *geo.DistanceMatrixResponse, err error)
	GetDistance(origin *model.Location, destination *model.Location) (element *geo.DistanceMatrixElement, err error)
	Save(location *model.Location) (err error)
	Delete(location *model.Location) (error error)
}

type locationServiceImpl struct {
	repo             repository.LocationRepository
	googleMapsAPIKey string
}

func (service *locationServiceImpl) GetById(id uuid.UUID) (location *model.Location, err error) {
	return service.repo.GetById(id)
}

func (service locationServiceImpl) GetDistance(origin *model.Location, destination *model.Location) (element *geo.DistanceMatrixElement, err error) {

	data, err := service.ResolveDistanceMatrix(origin, destination)

	if nil != err {
		return nil, err
	}

	element = &data.Rows[0].Elements[0]
	return element, nil

	// distance := data.Rows[0].Elements[0].Distance
	// duration := data.Rows[0].Elements[0].Duration

	// fmt.Println("Distance Text:", distance.Text)
	// fmt.Println("Distance Value:", distance.Value)
	// fmt.Println("Distance Value:", element)

	// fmt.Println("Duration Text:", duration.Text)
	// fmt.Println("Duration Value:", duration.Value)

}

func (service *locationServiceImpl) ResolveDistanceMatrix(origin *model.Location, destination *model.Location) (response *geo.DistanceMatrixResponse, err error) {
	return geo.ResolveDistanceMatrix(origin, destination)
}

func (service *locationServiceImpl) Resolve(input string) (location *model.Location, err error) {
	location, err = service.Search(input)
	if err != nil {
		return nil, err
	}
	if err = service.repo.Save(location); nil != err {
		return nil, err
	}
	return location, err
}

func (service *locationServiceImpl) Search(input string) (location *model.Location, err error) {

	if location, err := service.repo.Search(input); nil == err {
		return location, nil
	} else {
		if err.Error() != "record not found" {
			return nil, err
		}
	}

	if service.googleMapsAPIKey == "" {
		panic("Please provide googleMapsAPIKey")
	}
	gResult, err := geo.ResolveGeocodingInfo(service.googleMapsAPIKey, input)

	if nil != err {
		return nil, err
	}

	var address, city, state, country, cityCode, stateCode, countryCode, streetName, streetNumber, googleId, postcode string

	coordinates := gResult.Results[0].Geometry.Location
	address = gResult.Results[0].FormattedAddress
	googleId = gResult.Results[0].PlaceID

	for _, addrComp := range gResult.Results[0].AddressComponents {
		for _, Type := range addrComp.Types {
			switch Type {
			case "locality":
				city = addrComp.LongName
				cityCode = slug.SlugifyWithUnderscore(addrComp.ShortName)
				break
			case "administrative_area_level_1":
				state = addrComp.LongName
				stateCode = slug.SlugifyWithUnderscore(addrComp.ShortName)
				break
			case "country":
				country = addrComp.LongName
				countryCode = slug.SlugifyWithUnderscore(addrComp.ShortName)
				break
			case "route":
				streetName = addrComp.LongName
				break
			case "street_number":
				streetNumber = addrComp.LongName
				break

			case "postal_code":
				postcode = addrComp.LongName
				break
			}
		}
	}

	var street string

	if streetName != "" || streetNumber != "" {
		street = strings.Trim(fmt.Sprintf("%s %s", streetNumber, streetName), " \n\t\r")
	}

	location = &model.Location{
		Address:     address,
		Street:      street,
		City:        city,
		State:       state,
		Country:     country,
		CityCode:    cityCode,
		StateCode:   stateCode,
		CountryCode: countryCode,
		PostCode:    postcode,
		GoogleID:    &googleId,
		Coordinates: &model.LocationCoordinates{
			Latitude:  coordinates.Lat,
			Longitude: coordinates.Lng,
		},
	}

	return location, nil
}

func (service *locationServiceImpl) Save(location *model.Location) (err error) {
	return service.repo.Save(location)
}

func (service *locationServiceImpl) AssertLocationWithinPlace(location *model.Location, place *model.Place) (err error) {

	PlaceCode := place.Code
	CodeParts := strings.Split(PlaceCode, "-")

	switch place.Category {

	case model.PLACE_CITY:
		if len(CodeParts) != 3 {
			return errors.New(fmt.Sprintf("Invalid code[%s] for place %s", PlaceCode, place.Name))
		}

		if !(CodeParts[0] == location.CountryCode &&
			CodeParts[1] == location.StateCode &&
			CodeParts[2] == location.CityCode) {
			return errors.New(fmt.Sprintf("%s is not within %s", location.Address, place.Name))
		}
		return nil

	case model.PLACE_STATE:
		if len(CodeParts) != 2 {
			return errors.New(fmt.Sprintf("Invalid code[%s] for place %s", PlaceCode, place.Name))
		}

		if !(CodeParts[0] == location.CountryCode &&
			CodeParts[1] == location.StateCode) {
			return errors.New(fmt.Sprintf("%s is not within %s", location.Address, place.Name))
		}
		return nil

	case model.PLACE_COUNTRY:
		if len(CodeParts) != 1 {
			return errors.New(fmt.Sprintf("Invalid code[%s] for place %s", PlaceCode, place.Name))
		}
		if CodeParts[0] != location.CountryCode {
			return errors.New(fmt.Sprintf("%s is not within %s", location.Address, place.Name))
		}
		return nil

	default:
		return errors.New(fmt.Sprintf("Unsupported Place Category[%s]", place.Category))
	}

}

func (service *locationServiceImpl) Delete(location *model.Location) (err error) {
	err = service.repo.Delete(location)

	return err
}

func (service *locationServiceImpl) DeleteById(id uuid.UUID) (location *model.Location, err error) {
	location, err = service.repo.DeleteById(id)
	return location, err
}
