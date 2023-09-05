package repository

import (
	"sync"

	"github.com/google/uuid"
	"github.com/nkamuo/rasta-server/model"
	"gorm.io/gorm"
)

var fulfilmentRepo OrderFulfilmentRepository
var fulfilmentRepoMutext *sync.Mutex = &sync.Mutex{}

func GetOrderFulfilmentRepository() OrderFulfilmentRepository {
	fulfilmentRepoMutext.Lock()
	if fulfilmentRepo == nil {
		fulfilmentRepo = &fulfilmentRepository{db: model.DB}
	}
	fulfilmentRepoMutext.Unlock()
	return fulfilmentRepo
}

type OrderFulfilmentRepository interface {
	FindAll(page int, limit int) (fulfilments []model.OrderFulfilment, total int64, err error)
	GetById(id uuid.UUID, preload ...string) (fulfilment *model.OrderFulfilment, err error)
	Save(fulfilment *model.OrderFulfilment) (err error)
	Delete(fulfilment *model.OrderFulfilment) (error error)
	DeleteById(id uuid.UUID) (fulfilment *model.OrderFulfilment, err error)
}

type fulfilmentRepository struct {
	db *gorm.DB
}

func (repo *fulfilmentRepository) FindAll(page int, limit int) (fulfilments []model.OrderFulfilment, total int64, err error) {
	offset := (page - 1) * limit

	err = repo.db.
		Model(&model.OrderFulfilment{}).
		Count(&total).
		Limit(limit).
		Offset(offset).
		Order("created_at desc").
		Find(&fulfilments).Error
	if err != nil {
		return
	}
	return
}

func (repo *fulfilmentRepository) GetById(id uuid.UUID, preload ...string) (fulfilment *model.OrderFulfilment, err error) {
	query := model.DB. /*.Preload("Place")*/ Where("id = ?", id)

	for _, entry := range preload {
		query = query.Preload(entry)
	}

	if err = query.First(&fulfilment).Error; err != nil {
		return nil, err
	}
	return fulfilment, nil
}

func (repo *fulfilmentRepository) Save(fulfilment *model.OrderFulfilment) (err error) {
	if (uuid.UUID{} == fulfilment.ID) {
		//NEW - No ID yet
		return repo.db.Create(&fulfilment).Error
	}
	return repo.db.Updates(&fulfilment).Error
}

func (repo *fulfilmentRepository) Delete(fulfilment *model.OrderFulfilment) (err error) {
	return repo.db.Delete(&fulfilment).Error
}

func (repo *fulfilmentRepository) DeleteById(id uuid.UUID) (fulfilment *model.OrderFulfilment, err error) {
	fulfilment, err = repo.GetById(id)
	if err != nil {
		return nil, err
	}
	err = repo.db.Delete(fulfilment).Error
	return fulfilment, err
}
