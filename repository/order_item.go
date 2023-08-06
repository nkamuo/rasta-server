package repository

import (
	"sync"

	"github.com/google/uuid"
	"github.com/nkamuo/rasta-server/model"
	"gorm.io/gorm"
)

var orderItemRepo OrderItemRepository
var orderItemRepoMutext *sync.Mutex = &sync.Mutex{}

func GetOrderItemRepository() OrderItemRepository {
	orderItemRepoMutext.Lock()
	if orderItemRepo == nil {
		orderItemRepo = &orderItemRepository{db: model.DB}
	}
	orderItemRepoMutext.Unlock()
	return orderItemRepo
}

type OrderItemRepository interface {
	FindAll(page int, limit int) (orderItems []model.OrderItem, total int64, err error)
	GetById(id uuid.UUID) (orderItem *model.OrderItem, err error)
	Save(orderItem *model.OrderItem) (err error)
	Delete(orderItem *model.OrderItem) (error error)
	DeleteById(id uuid.UUID) (orderItem *model.OrderItem, err error)
}

type orderItemRepository struct {
	db *gorm.DB
}

func (repo *orderItemRepository) FindAll(page int, limit int) (orderItems []model.OrderItem, total int64, err error) {
	offset := (page - 1) * limit

	err = repo.db.
		Model(&model.OrderItem{}).
		Count(&total).
		Limit(limit).
		Offset(offset).
		Order("created_at desc").
		Find(&orderItems).Error
	if err != nil {
		return
	}
	return
}

func (repo *orderItemRepository) GetById(id uuid.UUID) (orderItem *model.OrderItem, err error) {
	if err = model.DB. /*.Preload("Place")*/ Where("id = ?", id).First(&orderItem).Error; err != nil {
		return nil, err
	}
	return orderItem, nil
}

func (repo *orderItemRepository) Save(orderItem *model.OrderItem) (err error) {
	if (uuid.UUID{} == orderItem.ID) {
		//NEW - No ID yet
		return repo.db.Create(&orderItem).Error
	}
	return repo.db.Updates(&orderItem).Error
}

func (repo *orderItemRepository) Delete(orderItem *model.OrderItem) (err error) {
	return repo.db.Delete(&orderItem).Error
}

func (repo *orderItemRepository) DeleteById(id uuid.UUID) (orderItem *model.OrderItem, err error) {
	orderItem, err = repo.GetById(id)
	if err != nil {
		return nil, err
	}
	err = repo.db.Delete(orderItem).Error
	return orderItem, err
}
