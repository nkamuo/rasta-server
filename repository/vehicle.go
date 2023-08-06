package repository

import (
	"sync"

	"github.com/google/uuid"
	"github.com/nkamuo/rasta-server/model"
	"gorm.io/gorm"
)

var vehicleRepo VehicleRepository
var vehicleRepoMutext *sync.Mutex = &sync.Mutex{}

func GetVehicleRepository() VehicleRepository {
	vehicleRepoMutext.Lock()
	if vehicleRepo == nil {
		vehicleRepo = &vehicleRepository{db: model.DB}
	}
	vehicleRepoMutext.Unlock()
	return vehicleRepo
}

type VehicleRepository interface {
	FindAll(page int, limit int) (vehicles []model.Vehicle, total int64, err error)
	GetById(id uuid.UUID) (vehicle *model.Vehicle, err error)
	Save(vehicle *model.Vehicle) (err error)
	Delete(vehicle *model.Vehicle) (error error)
	DeleteById(id uuid.UUID) (vehicle *model.Vehicle, err error)
}

type vehicleRepository struct {
	db *gorm.DB
}

func (repo *vehicleRepository) FindAll(page int, limit int) (vehicles []model.Vehicle, total int64, err error) {
	offset := (page - 1) * limit

	err = repo.db.
		Model(&model.Vehicle{}).
		Count(&total).
		Limit(limit).
		Offset(offset).
		Order("created_at desc").
		Find(&vehicles).Error
	if err != nil {
		return
	}
	return
}

func (repo *vehicleRepository) GetById(id uuid.UUID) (vehicle *model.Vehicle, err error) {
	if err = model.DB. /*.Preload("Place")*/ Where("id = ?", id).First(&vehicle).Error; err != nil {
		return nil, err
	}
	return vehicle, nil
}

func (repo *vehicleRepository) Save(vehicle *model.Vehicle) (err error) {
	if (uuid.UUID{} == vehicle.ID) {
		//NEW - No ID yet
		return repo.db.Create(&vehicle).Error
	}
	return repo.db.Updates(&vehicle).Error
}

func (repo *vehicleRepository) Delete(vehicle *model.Vehicle) (err error) {
	return repo.db.Delete(&vehicle).Error
}

func (repo *vehicleRepository) DeleteById(id uuid.UUID) (vehicle *model.Vehicle, err error) {
	vehicle, err = repo.GetById(id)
	if err != nil {
		return nil, err
	}
	err = repo.db.Delete(vehicle).Error
	return vehicle, err
}
