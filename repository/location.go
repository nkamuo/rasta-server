package repository

import (
	"errors"
	"strings"
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
	Search(input string) (location *model.Location, err error)
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

func (repo *locationRepository) Search(input string) (location *model.Location, err error) {
	if err = model.DB.First(&location, "id = ?", input).Error; err != nil {
		if err.Error() != "record not found" {
			return nil, err
		}
	} else {
		return location, nil
	}
	if err = model.DB.First(&location, "google_id = ?", input).Error; err != nil {
		if err.Error() != "record not found" {
			return nil, err
		}
	} else {
		return location, nil
	}

	parts := strings.Split(input, ",")

	if len(parts) == 2 {
		if err = model.DB.First(&location, "latitude = ? AND longitude = ?", parts[0], parts[1]).Error; err != nil {
			if err.Error() != "record not found" {
				return nil, err
			}
		} else {
			return location, nil
		}
	}

	return nil, errors.New("record not found")
}

func (repo *locationRepository) Save(location *model.Location) (err error) {
	if (uuid.UUID{} == location.ID) {
		//NEW - No ID yet
		return repo.db.Create(&location).Error
	}
	return repo.db.Updates(&location).Error
}

func (repo *locationRepository) Delete(location *model.Location) (err error) {
	return repo.db.Delete(&location).Error
}

func (repo *locationRepository) DeleteById(id uuid.UUID) (location *model.Location, err error) {
	location, err = repo.GetById(id)
	if err != nil {
		return nil, err
	}
	err = repo.db.Delete(&location).Error
	return location, err
}
