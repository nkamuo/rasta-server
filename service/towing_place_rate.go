package service

import (
	"sync"

	"github.com/google/uuid"
	"github.com/nkamuo/rasta-server/model"
	"github.com/nkamuo/rasta-server/repository"
)

var towingPlaceRateService TowingPlaceRateService
var towingPlaceRateRepoMutext *sync.Mutex = &sync.Mutex{}

func GetTowingPlaceRateService() TowingPlaceRateService {
	towingPlaceRateRepoMutext.Lock()
	if towingPlaceRateService == nil {
		towingPlaceRateService = &towingPlaceRateServiceImpl{repo: repository.GetTowingPlaceRateRepository()}
	}
	towingPlaceRateRepoMutext.Unlock()
	return towingPlaceRateService
}

type TowingPlaceRateService interface {
	GetById(id uuid.UUID) (towingPlaceRate *model.TowingPlaceRate, err error)
	GetByPlaceAndDistance(place model.Place, distance int64) (rate *model.TowingPlaceRate, err error)
	Save(towingPlaceRate *model.TowingPlaceRate) (err error)
	Delete(towingPlaceRate *model.TowingPlaceRate) (error error)
}

type towingPlaceRateServiceImpl struct {
	repo repository.TowingPlaceRateRepository
}

func (service *towingPlaceRateServiceImpl) GetById(id uuid.UUID) (towingPlaceRate *model.TowingPlaceRate, err error) {
	return service.repo.GetById(id)
}
func (service *towingPlaceRateServiceImpl) GetByPlaceAndDistance(place model.Place, distance int64) (rate *model.TowingPlaceRate, err error) {
	return service.repo.GetByPlaceAndDistance(place, distance)
}

func (service *towingPlaceRateServiceImpl) Save(towingPlaceRate *model.TowingPlaceRate) (err error) {
	return service.repo.Save(towingPlaceRate)
}

func (service *towingPlaceRateServiceImpl) Delete(towingPlaceRate *model.TowingPlaceRate) (err error) {
	err = service.repo.Delete(towingPlaceRate)

	return err
}

func (service *towingPlaceRateServiceImpl) DeleteById(id uuid.UUID) (towingPlaceRate *model.TowingPlaceRate, err error) {
	towingPlaceRate, err = service.repo.DeleteById(id)
	return towingPlaceRate, err
}
