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
	GetById(id uuid.UUID, preload ...string) (order *model.Order, err error)
	GetByFulfilment(fulfilment model.OrderFulfilment) (order *model.Order, err error)
	CountByRespondent(respondent *model.Respondent, isActive *bool) (count int64, err error)
	Save(order *model.Order) (err error)
	Update(order *model.Order, fields map[string]interface{}) (err error)
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

func (repo *orderRepository) GetById(id uuid.UUID, preload ...string) (order *model.Order, err error) {
	query := model.DB

	for _, pLoad := range preload {
		query = query.Preload(pLoad)
	}

	if err = query.Where("id = ?", id).First(&order).Error; err != nil {
		return nil, err
	}
	return order, nil
}

func (repo *orderRepository) GetByFulfilment(fulfilment model.OrderFulfilment) (order *model.Order, err error) {
	if err = model.DB. /*.Preload("Place")*/ Where("fulfilment_id = ?", fulfilment.ID).First(&order).Error; err != nil {
		return nil, err
	}
	return order, nil
}

func (repo *orderRepository) CountByRespondent(respondent *model.Respondent, isActive *bool) (count int64, err error) {
	query := repo.db.Model(&model.Order{})
	query = query.Joins("LEFT JOIN order_fulfilments ON orders.fulfilment_id = order_fulfilments.id")
	query = query.Where("order_fulfilments.responder_id = ?", respondent.ID)
	if isActive != nil {
		if *isActive {
			query = query.Where("(order_fulfilments.client_confirmed_at IS NULL) AND (order_fulfilments.auto_confirmed_at IS NULL)")
		} else {
			query = query.Where("(order_fulfilments.client_confirmed_at IS NOT NULL) OR (order_fulfilments.auto_confirmed_at IS NOT NULL)")
		}
	}
	err = query.Count(&count).Error
	return count, err
}

func (repo *orderRepository) Save(order *model.Order) (err error) {
	if (uuid.UUID{} == order.ID) {
		//NEW - No ID yet
		return repo.db.Create(&order).Error
	}
	return repo.db.Updates(&order).Error
}

func (repo *orderRepository) Update(order *model.Order, fields map[string]interface{}) (err error) {
	return repo.db.Model(order).Updates(fields).Error
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
