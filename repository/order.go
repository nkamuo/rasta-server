package repository

import (
	"sync"

	"github.com/google/uuid"
	"github.com/nkamuo/rasta-server/model"
	"gorm.io/gorm"
)

var orderRepo OrderRepository
var orderRepoMutext *sync.Mutex = &sync.Mutex{}

func GetOrderRepository() OrderRepository {
	orderRepoMutext.Lock()
	if orderRepo == nil {
		orderRepo = &orderRepository{db: model.DB}
	}
	orderRepoMutext.Unlock()
	return orderRepo
}

type OrderRepository interface {
	FindAll(page int, limit int) (orders []model.Order, total int64, err error)
	GetById(id uuid.UUID) (order *model.Order, err error)
	Save(order *model.Order) (err error)
	Delete(order *model.Order) (error error)
	DeleteById(id uuid.UUID) (order *model.Order, err error)
}

type orderRepository struct {
	db *gorm.DB
}

func (repo *orderRepository) FindAll(page int, limit int) (orders []model.Order, total int64, err error) {
	offset := (page - 1) * limit

	err = repo.db.
		Model(&model.Order{}).
		Count(&total).
		Limit(limit).
		Offset(offset).
		Order("created_at desc").
		Find(&orders).Error
	if err != nil {
		return
	}
	return
}

func (repo *orderRepository) GetById(id uuid.UUID) (order *model.Order, err error) {
	if err = model.DB. /*.Preload("Place")*/ Where("id = ?", id).First(&order).Error; err != nil {
		return nil, err
	}
	return order, nil
}

func (repo *orderRepository) Save(order *model.Order) (err error) {
	if (uuid.UUID{} == order.ID) {
		//NEW - No ID yet
		return repo.db.Create(&order).Error
	}
	return repo.db.Updates(&order).Error
}

func (repo *orderRepository) Delete(order *model.Order) (err error) {
	return repo.db.Delete(&order).Error
}

func (repo *orderRepository) DeleteById(id uuid.UUID) (order *model.Order, err error) {
	order, err = repo.GetById(id)
	if err != nil {
		return nil, err
	}
	err = repo.db.Delete(order).Error
	return order, err
}
