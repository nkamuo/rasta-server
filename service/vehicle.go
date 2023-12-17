package service

import (
	"sync"

	"github.com/google/uuid"
	"github.com/nkamuo/rasta-server/model"
	"github.com/nkamuo/rasta-server/repository"
)

var vehicleService VehicleService
var vehicleRepoMutext *sync.Mutex = &sync.Mutex{}

func GetVehicleService() VehicleService {
	vehicleRepoMutext.Lock()
	if vehicleService == nil {
		vehicleService = &vehicleServiceImpl{repo: repository.GetVehicleRepository()}
	}
	vehicleRepoMutext.Unlock()
	return vehicleService
}

type VehicleService interface {
	GetById(id uuid.UUID, preload ...string) (vehicle *model.Vehicle, err error)
	// GetByEmail(email string) (vehicle *model.Vehicle, err error)
	// GetByPhone(phone string) (vehicle *model.Vehicle, err error)
	Save(vehicle *model.Vehicle) (err error)
	Delete(vehicle *model.Vehicle) (error error)
}

type vehicleServiceImpl struct {
	repo repository.VehicleRepository
}

func (service *vehicleServiceImpl) GetById(id uuid.UUID, preload ...string) (vehicle *model.Vehicle, err error) {
	return service.repo.GetById(id, preload...)
}

func (service *vehicleServiceImpl) Save(vehicle *model.Vehicle) (err error) {
	return service.repo.Save(vehicle)
}

func (service *vehicleServiceImpl) Delete(vehicle *model.Vehicle) (err error) {
	err = service.repo.Delete(vehicle)

	return err
}

func (service *vehicleServiceImpl) DeleteById(id uuid.UUID) (vehicle *model.Vehicle, err error) {
	vehicle, err = service.repo.DeleteById(id)
	return vehicle, err
}
