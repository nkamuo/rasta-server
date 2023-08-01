package repository

import (
	"sync"

	"github.com/google/uuid"
	"github.com/nkamuo/rasta-server/model"
	"gorm.io/gorm"
)

var fuelTypeRepo FuelTypeRepository
var fuelTypeRepoMutext *sync.Mutex = &sync.Mutex{}

func GetFuelTypeRepository() FuelTypeRepository {
	fuelTypeRepoMutext.Lock()
	if fuelTypeRepo == nil {
		fuelTypeRepo = &fuelTypeRepository{db: model.DB}
	}
	fuelTypeRepoMutext.Unlock()
	return fuelTypeRepo
}

type FuelTypeRepository interface {
	FindAll(page int, limit int) (fuelTypes []model.FuelType, total int64, err error)
	GetById(id uuid.UUID) (fuelType *model.FuelType, err error)
	GetByCode(code string) (fuelType *model.FuelType, err error)
	GetRateForTypeInPlace(fuelType model.FuelType, place model.Place) (rate *model.FuelTypePlaceRate, err error)
	Save(fuelType *model.FuelType) (err error)
	Delete(fuelType *model.FuelType) (error error)
	DeleteById(id uuid.UUID) (fuelType *model.FuelType, err error)
}

type fuelTypeRepository struct {
	db *gorm.DB
}

func (repo *fuelTypeRepository) FindAll(page int, limit int) (fuelTypes []model.FuelType, total int64, err error) {
	offset := (page - 1) * limit

	err = repo.db.
		Model(&model.FuelType{}).
		Count(&total).
		Limit(limit).
		Offset(offset).
		Order("created_at desc").
		Find(&fuelTypes).Error
	if err != nil {
		return
	}
	return
}

func (repo *fuelTypeRepository) GetById(id uuid.UUID) (fuelType *model.FuelType, err error) {
	if err = model.DB.Where("id = ?", id).First(&fuelType).Error; err != nil {
		return nil, err
	}
	return fuelType, nil
}
func (repo *fuelTypeRepository) GetByCode(code string) (fuelType *model.FuelType, err error) {
	if err = model.DB.Where("code = ?", code).First(&fuelType).Error; err != nil {
		return nil, err
	}
	return fuelType, nil
}

func (repo *fuelTypeRepository) GetRateForTypeInPlace(fuelType model.FuelType, place model.Place) (rate *model.FuelTypePlaceRate, err error) {
	if err = model.DB.Where("fuel_type_id = ? AND place_id = ?", fuelType.ID, place.ID).First(&rate).Error; err != nil {
		return nil, err
	}
	return rate, nil
}

func (repo *fuelTypeRepository) Save(fuelType *model.FuelType) (err error) {
	if (uuid.UUID{} == fuelType.ID) {
		//NEW - No ID yet
		repo.db.Create(&fuelType)
		return nil
	}
	repo.db.Updates(&fuelType)
	return nil
}

func (repo *fuelTypeRepository) Delete(fuelType *model.FuelType) (err error) {
	repo.db.Delete(&fuelType)
	return nil
}

func (repo *fuelTypeRepository) DeleteById(id uuid.UUID) (fuelType *model.FuelType, err error) {
	fuelType, err = repo.GetById(id)
	if err != nil {
		return nil, err
	}
	return fuelType, nil
}
