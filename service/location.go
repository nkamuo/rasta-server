package service

import (
	"sync"

	"github.com/google/uuid"
	"github.com/nkamuo/rasta-server/model"
	"github.com/nkamuo/rasta-server/repository"
)

var locationService LocationService
var locationRepoMutext *sync.Mutex = &sync.Mutex{}

func GetLocationService() LocationService {
	locationRepoMutext.Lock()
	if locationService == nil {
		locationService = &locationServiceImpl{repo: repository.GetLocationRepository()}
	}
	locationRepoMutext.Unlock()
	return locationService
}

type LocationService interface {
	GetById(id uuid.UUID) (location *model.Location, err error)
	Resolve(data string) (location *model.Location, err error)
	Save(location *model.Location) (err error)
	Delete(location *model.Location) (error error)
}

type locationServiceImpl struct {
	repo repository.LocationRepository
}

func (service *locationServiceImpl) GetById(id uuid.UUID) (location *model.Location, err error) {
	return service.repo.GetById(id)
}

func (service *locationServiceImpl) Resolve(data string) (location *model.Location, err error) {

	return
}

func (service *locationServiceImpl) Save(location *model.Location) (err error) {
	return service.repo.Save(location)
}

func (service *locationServiceImpl) Delete(location *model.Location) (err error) {
	err = service.repo.Delete(location)

	return err
}

func (service *locationServiceImpl) DeleteById(id uuid.UUID) (location *model.Location, err error) {
	location, err = service.repo.DeleteById(id)
	return location, err
}
