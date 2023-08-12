package repository

import (
	"sync"

	"github.com/google/uuid"
	"github.com/nkamuo/rasta-server/model"
	"gorm.io/gorm"
)

var requestRepo RequestRepository
var requestRepoMutext *sync.Mutex = &sync.Mutex{}

func GetRequestRepository() RequestRepository {
	requestRepoMutext.Lock()
	if requestRepo == nil {
		requestRepo = &requestRepository{db: model.DB}
	}
	requestRepoMutext.Unlock()
	return requestRepo
}

type RequestRepository interface {
	FindAll(page int, limit int) (requests []model.Request, total int64, err error)
	GetById(id uuid.UUID) (request *model.Request, err error)
	Save(request *model.Request) (err error)
	Delete(request *model.Request) (error error)
	DeleteById(id uuid.UUID) (request *model.Request, err error)
}

type requestRepository struct {
	db *gorm.DB
}

func (repo *requestRepository) FindAll(page int, limit int) (requests []model.Request, total int64, err error) {
	offset := (page - 1) * limit

	err = repo.db.
		Model(&model.Request{}).
		Count(&total).
		Limit(limit).
		Offset(offset).
		Order("created_at desc").
		Find(&requests).Error
	if err != nil {
		return
	}
	return
}

func (repo *requestRepository) GetById(id uuid.UUID) (request *model.Request, err error) {
	if err = model.DB. /*.Preload("Place")*/ Where("id = ?", id).First(&request).Error; err != nil {
		return nil, err
	}
	return request, nil
}

func (repo *requestRepository) Save(request *model.Request) (err error) {
	if (uuid.UUID{} == request.ID) {
		//NEW - No ID yet
		return repo.db.Create(&request).Error
	}
	return repo.db.Updates(&request).Error
}

func (repo *requestRepository) Delete(request *model.Request) (err error) {
	return repo.db.Delete(&request).Error
}

func (repo *requestRepository) DeleteById(id uuid.UUID) (request *model.Request, err error) {
	request, err = repo.GetById(id)
	if err != nil {
		return nil, err
	}
	err = repo.db.Delete(request).Error
	return request, err
}
