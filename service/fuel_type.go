package service

import (
	"sync"

	"github.com/google/uuid"
	"github.com/nkamuo/rasta-server/model"
	"github.com/nkamuo/rasta-server/repository"
)

var fuelTypeService FuelTypeService
var fuelTypeRepoMutext *sync.Mutex = &sync.Mutex{}

func GetFuelTypeService() FuelTypeService {
	fuelTypeRepoMutext.Lock()
	if fuelTypeService == nil {
		fuelTypeService = &fuelTypeServiceImpl{repo: repository.GetFuelTypeRepository()}
	}
	fuelTypeRepoMutext.Unlock()
	return fuelTypeService
}

type FuelTypeService interface {
	GetById(id uuid.UUID) (fuelType *model.FuelType, err error)
	Save(fuelType *model.FuelType) (err error)
	Delete(fuelType *model.FuelType) (error error)
}

type fuelTypeServiceImpl struct {
	repo repository.FuelTypeRepository
}

func (service *fuelTypeServiceImpl) GetById(id uuid.UUID) (fuelType *model.FuelType, err error) {
	return service.repo.GetById(id)
}

func (service *fuelTypeServiceImpl) Save(fuelType *model.FuelType) (err error) {
	return service.repo.Save(fuelType)
}

func (service *fuelTypeServiceImpl) Delete(fuelType *model.FuelType) (err error) {
	err = service.repo.Delete(fuelType)

	return err
}

func (service *fuelTypeServiceImpl) DeleteById(id uuid.UUID) (fuelType *model.FuelType, err error) {
	fuelType, err = service.repo.DeleteById(id)
	return fuelType, err
}
