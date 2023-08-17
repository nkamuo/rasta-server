package repository

import (
	"sync"

	"github.com/google/uuid"
	"github.com/nkamuo/rasta-server/model"
	"gorm.io/gorm"
)

var towingPlaceRateRepo TowingPlaceRateRepository
var towingPlaceRateRepoMutext *sync.Mutex = &sync.Mutex{}

func GetTowingPlaceRateRepository() TowingPlaceRateRepository {
	towingPlaceRateRepoMutext.Lock()
	if towingPlaceRateRepo == nil {
		towingPlaceRateRepo = &towingPlaceRateRepository{db: model.DB}
	}
	towingPlaceRateRepoMutext.Unlock()
	return towingPlaceRateRepo
}

type TowingPlaceRateRepository interface {
	FindAll(page int, limit int) (towingPlaceRates []model.TowingPlaceRate, total int64, err error)
	GetById(id uuid.UUID) (towingPlaceRate *model.TowingPlaceRate, err error)
	GetByCode(code string) (towingPlaceRate *model.TowingPlaceRate, err error)
	GetByPlaceAndDistance(place model.Place, distance int64) (rate *model.TowingPlaceRate, err error)
	Save(towingPlaceRate *model.TowingPlaceRate) (err error)
	Delete(towingPlaceRate *model.TowingPlaceRate) (error error)
	DeleteById(id uuid.UUID) (towingPlaceRate *model.TowingPlaceRate, err error)
}

type towingPlaceRateRepository struct {
	db *gorm.DB
}

func (repo *towingPlaceRateRepository) FindAll(page int, limit int) (towingPlaceRates []model.TowingPlaceRate, total int64, err error) {
	offset := (page - 1) * limit

	err = repo.db.
		Model(&model.TowingPlaceRate{}).
		Count(&total).
		Limit(limit).
		Offset(offset).
		Order("created_at desc").
		Find(&towingPlaceRates).Error
	if err != nil {
		return
	}
	return
}

func (repo *towingPlaceRateRepository) GetById(id uuid.UUID) (towingPlaceRate *model.TowingPlaceRate, err error) {
	if err = model.DB.Where("id = ?", id).First(&towingPlaceRate).Error; err != nil {
		return nil, err
	}
	return towingPlaceRate, nil
}
func (repo *towingPlaceRateRepository) GetByCode(code string) (towingPlaceRate *model.TowingPlaceRate, err error) {
	if err = model.DB.Where("code = ?", code).First(&towingPlaceRate).Error; err != nil {
		return nil, err
	}
	return towingPlaceRate, nil
}

func (repo *towingPlaceRateRepository) GetByPlaceAndDistance(place model.Place, distance int64) (rate *model.TowingPlaceRate, err error) {
	// var rates []model.TowingPlaceRate
	const QUERY = "place_id = ? AND (? BETWEEN min_distance AND max_distance)"
	if err = model.DB.Where(QUERY, place.ID, distance). /*.Order("sequence ASC")*/
								First(&rate).Error; err != nil {
		return nil, err
	}

	return rate, nil
}

func (repo *towingPlaceRateRepository) Save(towingPlaceRate *model.TowingPlaceRate) (err error) {
	if (uuid.UUID{} == towingPlaceRate.ID) {
		//NEW - No ID yet
		return repo.db.Create(&towingPlaceRate).Error
	}
	return repo.db.Updates(&towingPlaceRate).Error
}

func (repo *towingPlaceRateRepository) Delete(towingPlaceRate *model.TowingPlaceRate) (err error) {
	return repo.db.Delete(&towingPlaceRate).Error
}

func (repo *towingPlaceRateRepository) DeleteById(id uuid.UUID) (towingPlaceRate *model.TowingPlaceRate, err error) {
	towingPlaceRate, err = repo.GetById(id)
	if err != nil {
		return nil, err
	}
	err = repo.db.Delete(towingPlaceRate).Error
	return towingPlaceRate, err
}
