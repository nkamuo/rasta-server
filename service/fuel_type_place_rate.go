package service

import (
	"sync"

	"github.com/google/uuid"
	"github.com/nkamuo/rasta-server/model"
	"github.com/nkamuo/rasta-server/repository"
)

var fuelTypePlaceRateService FuelTypePlaceRateService
var fuelTypePlaceRateRepoMutext *sync.Mutex = &sync.Mutex{}

func GetFuelTypePlaceRateService() FuelTypePlaceRateService {
	fuelTypePlaceRateRepoMutext.Lock()
	if fuelTypePlaceRateService == nil {
		fuelTypePlaceRateService = &fuelTypePlaceRateServiceImpl{repo: repository.GetFuelTypePlaceRateRepository()}
	}
	fuelTypePlaceRateRepoMutext.Unlock()
	return fuelTypePlaceRateService
}

type FuelTypePlaceRateService interface {
	GetById(id uuid.UUID) (fuelTypePlaceRate *model.FuelTypePlaceRate, err error)
	Save(fuelTypePlaceRate *model.FuelTypePlaceRate) (err error)
	Delete(fuelTypePlaceRate *model.FuelTypePlaceRate) (error error)
}

type fuelTypePlaceRateServiceImpl struct {
	repo repository.FuelTypePlaceRateRepository
}

func (service *fuelTypePlaceRateServiceImpl) GetById(id uuid.UUID) (fuelTypePlaceRate *model.FuelTypePlaceRate, err error) {
	return service.repo.GetById(id)
}

func (service *fuelTypePlaceRateServiceImpl) Save(fuelTypePlaceRate *model.FuelTypePlaceRate) (err error) {
	return service.repo.Save(fuelTypePlaceRate)
}

func (service *fuelTypePlaceRateServiceImpl) Delete(fuelTypePlaceRate *model.FuelTypePlaceRate) (err error) {
	err = service.repo.Delete(fuelTypePlaceRate)

	return err
}

func (service *fuelTypePlaceRateServiceImpl) DeleteById(id uuid.UUID) (fuelTypePlaceRate *model.FuelTypePlaceRate, err error) {
	fuelTypePlaceRate, err = service.repo.DeleteById(id)
	return fuelTypePlaceRate, err
}
