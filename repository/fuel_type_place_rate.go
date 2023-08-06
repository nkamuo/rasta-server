package repository

import (
	"sync"

	"github.com/google/uuid"
	"github.com/nkamuo/rasta-server/model"
	"gorm.io/gorm"
)

var fuelTypePlaceRateRepo FuelTypePlaceRateRepository
var fuelTypePlaceRateRepoMutext *sync.Mutex = &sync.Mutex{}

func GetFuelTypePlaceRateRepository() FuelTypePlaceRateRepository {
	fuelTypePlaceRateRepoMutext.Lock()
	if fuelTypePlaceRateRepo == nil {
		fuelTypePlaceRateRepo = &fuelTypePlaceRateRepository{db: model.DB}
	}
	fuelTypePlaceRateRepoMutext.Unlock()
	return fuelTypePlaceRateRepo
}

type FuelTypePlaceRateRepository interface {
	FindAll(page int, limit int) (fuelTypePlaceRates []model.FuelTypePlaceRate, total int64, err error)
	GetById(id uuid.UUID) (fuelTypePlaceRate *model.FuelTypePlaceRate, err error)
	GetByCode(code string) (fuelTypePlaceRate *model.FuelTypePlaceRate, err error)
	GetRateForFuelTypeInPlace(fuelType model.FuelType, place model.Place) (rate *model.FuelTypePlaceRate, err error)
	Save(fuelTypePlaceRate *model.FuelTypePlaceRate) (err error)
	Delete(fuelTypePlaceRate *model.FuelTypePlaceRate) (error error)
	DeleteById(id uuid.UUID) (fuelTypePlaceRate *model.FuelTypePlaceRate, err error)
}

type fuelTypePlaceRateRepository struct {
	db *gorm.DB
}

func (repo *fuelTypePlaceRateRepository) FindAll(page int, limit int) (fuelTypePlaceRates []model.FuelTypePlaceRate, total int64, err error) {
	offset := (page - 1) * limit

	err = repo.db.
		Model(&model.FuelTypePlaceRate{}).
		Count(&total).
		Limit(limit).
		Offset(offset).
		Order("created_at desc").
		Find(&fuelTypePlaceRates).Error
	if err != nil {
		return
	}
	return
}

func (repo *fuelTypePlaceRateRepository) GetById(id uuid.UUID) (fuelTypePlaceRate *model.FuelTypePlaceRate, err error) {
	if err = model.DB.Where("id = ?", id).First(&fuelTypePlaceRate).Error; err != nil {
		return nil, err
	}
	return fuelTypePlaceRate, nil
}
func (repo *fuelTypePlaceRateRepository) GetByCode(code string) (fuelTypePlaceRate *model.FuelTypePlaceRate, err error) {
	if err = model.DB.Where("code = ?", code).First(&fuelTypePlaceRate).Error; err != nil {
		return nil, err
	}
	return fuelTypePlaceRate, nil
}

func (repo *fuelTypePlaceRateRepository) GetRateForFuelTypeInPlace(fuelType model.FuelType, place model.Place) (rate *model.FuelTypePlaceRate, err error) {
	if err = model.DB.Where("fuel_type_id = ? AND place_id = ?", fuelType.ID, place.ID).First(&rate).Error; err != nil {
		return nil, err
	}
	return rate, nil
}

func (repo *fuelTypePlaceRateRepository) Save(fuelTypePlaceRate *model.FuelTypePlaceRate) (err error) {
	if (uuid.UUID{} == fuelTypePlaceRate.ID) {
		//NEW - No ID yet
		return repo.db.Create(&fuelTypePlaceRate).Error
	}
	return repo.db.Updates(&fuelTypePlaceRate).Error
}

func (repo *fuelTypePlaceRateRepository) Delete(fuelTypePlaceRate *model.FuelTypePlaceRate) (err error) {
	return repo.db.Delete(&fuelTypePlaceRate).Error
}

func (repo *fuelTypePlaceRateRepository) DeleteById(id uuid.UUID) (fuelTypePlaceRate *model.FuelTypePlaceRate, err error) {
	fuelTypePlaceRate, err = repo.GetById(id)
	if err != nil {
		return nil, err
	}
	err = repo.db.Delete(fuelTypePlaceRate).Error
	return fuelTypePlaceRate, err
}
