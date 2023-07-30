package service

import (
	"sync"

	"github.com/google/uuid"
	"github.com/nkamuo/rasta-server/model"
	"github.com/nkamuo/rasta-server/repository"
)

var vehicleModelService VehicleModelService
var vehicleModelRepoMutext *sync.Mutex = &sync.Mutex{}

func GetVehicleModelService() VehicleModelService {
	vehicleModelRepoMutext.Lock()
	if vehicleModelService == nil {
		vehicleModelService = &vehicleModelServiceImpl{repo: repository.GetVehicleModelRepository()}
	}
	vehicleModelRepoMutext.Unlock()
	return vehicleModelService
}

type VehicleModelService interface {
	GetById(id uuid.UUID) (vehicleModel *model.VehicleModel, err error)
	// GetByEmail(email string) (vehicleModel *model.VehicleModel, err error)
	// GetByPhone(phone string) (vehicleModel *model.VehicleModel, err error)
	Save(vehicleModel *model.VehicleModel) (err error)
	Delete(vehicleModel *model.VehicleModel) (error error)
}

type vehicleModelServiceImpl struct {
	repo repository.VehicleModelRepository
}

func (service *vehicleModelServiceImpl) GetById(id uuid.UUID) (vehicleModel *model.VehicleModel, err error) {
	return service.repo.GetById(id)
}

func (service *vehicleModelServiceImpl) Save(vehicleModel *model.VehicleModel) (err error) {
	return service.repo.Save(vehicleModel)
}

func (service *vehicleModelServiceImpl) Delete(vehicleModel *model.VehicleModel) (err error) {
	err = service.repo.Delete(vehicleModel)

	return err
}

func (service *vehicleModelServiceImpl) DeleteById(id uuid.UUID) (vehicleModel *model.VehicleModel, err error) {
	vehicleModel, err = service.repo.DeleteById(id)
	return vehicleModel, err
}
