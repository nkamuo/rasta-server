package repository

import (
	"sync"

	"github.com/google/uuid"
	"github.com/nkamuo/rasta-server/model"
	"gorm.io/gorm"
)

var placeRepo PlaceRepository
var placeRepoMutext *sync.Mutex = &sync.Mutex{}

func GetPlaceRepository() PlaceRepository {
	placeRepoMutext.Lock()
	if placeRepo == nil {
		placeRepo = &placeRepository{db: model.DB}
	}
	placeRepoMutext.Unlock()
	return placeRepo
}

type PlaceRepository interface {
	FindAll(page int, limit int) (places []model.Place, total int64, err error)
	GetById(id uuid.UUID) (place *model.Place, err error)
	Save(place *model.Place) (err error)
	Delete(place *model.Place) (error error)
	DeleteById(id uuid.UUID) (place *model.Place, err error)
}

type placeRepository struct {
	db *gorm.DB
}

func (repo *placeRepository) FindAll(page int, limit int) (places []model.Place, total int64, err error) {
	offset := (page - 1) * limit

	err = repo.db.
		Model(&model.Place{}).
		Count(&total).
		Limit(limit).
		Offset(offset).
		Order("created_at desc").
		Find(&places).Error
	if err != nil {
		return
	}
	return
}

func (repo *placeRepository) GetById(id uuid.UUID) (place *model.Place, err error) {
	if err = model.DB. /*.Joins("OperatorUser")*/ First(&place, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return place, nil
}

func (repo *placeRepository) Save(place *model.Place) (err error) {
	if (uuid.UUID{} == place.ID) {
		//NEW - No ID yet
		return repo.db.Create(&place).Error
	}
	return repo.db.Updates(&place).Error
}

func (repo *placeRepository) Delete(place *model.Place) (err error) {
	repo.db.Delete(&place)
	return repo.db.Error
}

func (repo *placeRepository) DeleteById(id uuid.UUID) (place *model.Place, err error) {
	place, err = repo.GetById(id)
	if err != nil {
		return nil, err
	}
	repo.db.Delete(&place)
	return place, repo.db.Error
}
