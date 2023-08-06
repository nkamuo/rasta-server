package repository

import (
	"sync"

	"github.com/google/uuid"
	"github.com/nkamuo/rasta-server/model"
	"gorm.io/gorm"
)

var vehicleModelRepo VehicleModelRepository
var vehicleModelRepoMutext *sync.Mutex = &sync.Mutex{}

func GetVehicleModelRepository() VehicleModelRepository {
	vehicleModelRepoMutext.Lock()
	if vehicleModelRepo == nil {
		vehicleModelRepo = &vehicleModelRepository{db: model.DB}
	}
	vehicleModelRepoMutext.Unlock()
	return vehicleModelRepo
}

type VehicleModelRepository interface {
	FindAll(page int, limit int) (vehicleModels []model.VehicleModel, total int64, err error)
	GetById(id uuid.UUID) (vehicleModel *model.VehicleModel, err error)
	Save(vehicleModel *model.VehicleModel) (err error)
	Delete(vehicleModel *model.VehicleModel) (error error)
	DeleteById(id uuid.UUID) (vehicleModel *model.VehicleModel, err error)
}

type vehicleModelRepository struct {
	db *gorm.DB
}

func (repo *vehicleModelRepository) FindAll(page int, limit int) (vehicleModels []model.VehicleModel, total int64, err error) {
	offset := (page - 1) * limit

	err = repo.db.
		Model(&model.VehicleModel{}).
		Count(&total).
		Limit(limit).
		Offset(offset).
		Order("created_at desc").
		Find(&vehicleModels).Error
	if err != nil {
		return
	}
	return
}

func (repo *vehicleModelRepository) GetById(id uuid.UUID) (vehicleModel *model.VehicleModel, err error) {
	if err = model.DB. /*.Preload("Place")*/ Where("id = ?", id).First(&vehicleModel).Error; err != nil {
		return nil, err
	}
	return vehicleModel, nil
}

func (repo *vehicleModelRepository) Save(vehicleModel *model.VehicleModel) (err error) {
	if (uuid.UUID{} == vehicleModel.ID) {
		return repo.db.Create(&vehicleModel).Error
	}
	return repo.db.Updates(&vehicleModel).Error
}

func (repo *vehicleModelRepository) Delete(vehicleModel *model.VehicleModel) (err error) {
	return repo.db.Delete(&vehicleModel).Error
}

func (repo *vehicleModelRepository) DeleteById(id uuid.UUID) (vehicleModel *model.VehicleModel, err error) {
	vehicleModel, err = repo.GetById(id)
	if err != nil {
		return nil, err
	}
	err = repo.db.Delete(vehicleModel).Error
	return vehicleModel, err
}
