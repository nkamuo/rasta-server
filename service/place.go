package service

import (
	"sync"

	"github.com/google/uuid"
	"github.com/nkamuo/rasta-server/model"
	"github.com/nkamuo/rasta-server/repository"
)

var placeService PlaceService
var placeRepoMutext *sync.Mutex = &sync.Mutex{}

func GetPlaceService() PlaceService {
	placeRepoMutext.Lock()
	if placeService == nil {
		placeService = &placeServiceImpl{repo: repository.GetPlaceRepository()}
	}
	placeRepoMutext.Unlock()
	return placeService
}

type PlaceService interface {
	GetById(id uuid.UUID) (place *model.Place, err error)
	Save(place *model.Place) (err error)
	Delete(place *model.Place) (error error)
}

type placeServiceImpl struct {
	repo repository.PlaceRepository
}

func (service *placeServiceImpl) GetById(id uuid.UUID) (place *model.Place, err error) {
	return service.repo.GetById(id)
}

func (service *placeServiceImpl) Save(place *model.Place) (err error) {
	return service.repo.Save(place)
}

func (service *placeServiceImpl) Delete(place *model.Place) (err error) {
	err = service.repo.Delete(place)

	return err
}

func (service *placeServiceImpl) DeleteById(id uuid.UUID) (place *model.Place, err error) {
	place, err = service.repo.DeleteById(id)
	return place, err
}
