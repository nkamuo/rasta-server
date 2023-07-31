package repository

import (
	"sync"

	"github.com/google/uuid"
	"github.com/nkamuo/rasta-server/model"
	"gorm.io/gorm"
)

var locationRepo LocationRepository
var locationRepoMutext *sync.Mutex = &sync.Mutex{}

func GetLocationRepository() LocationRepository {
	locationRepoMutext.Lock()
	if locationRepo == nil {
		locationRepo = &locationRepository{db: model.DB}
	}
	locationRepoMutext.Unlock()
	return locationRepo
}

type LocationRepository interface {
	FindAll(page int, limit int) (locations []model.Location, total int64, err error)
	GetById(id uuid.UUID) (location *model.Location, err error)
	Save(location *model.Location) (err error)
	Delete(location *model.Location) (error error)
	DeleteById(id uuid.UUID) (location *model.Location, err error)
}

type locationRepository struct {
	db *gorm.DB
}

func (repo *locationRepository) FindAll(page int, limit int) (locations []model.Location, total int64, err error) {
	offset := (page - 1) * limit

	err = repo.db.
		Model(&model.Location{}).
		Count(&total).
		Limit(limit).
		Offset(offset).
		Order("created_at desc").
		Find(&locations).Error
	if err != nil {
		return
	}
	return
}

func (repo *locationRepository) GetById(id uuid.UUID) (location *model.Location, err error) {
	if err = model.DB. /*.Joins("OperatorUser")*/ First(&location, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return location, nil
}

func (repo *locationRepository) Save(location *model.Location) (err error) {
	if (uuid.UUID{} == location.ID) {
		//NEW - No ID yet
		repo.db.Create(&location)
		return repo.db.Error
	}
	repo.db.Updates(&location)
	return repo.db.Error
}

func (repo *locationRepository) Delete(location *model.Location) (err error) {
	repo.db.Delete(&location)
	return repo.db.Error
}

func (repo *locationRepository) DeleteById(id uuid.UUID) (location *model.Location, err error) {
	location, err = repo.GetById(id)
	if err != nil {
		return nil, err
	}
	repo.db.Delete(&location)
	return location, repo.db.Error
}
