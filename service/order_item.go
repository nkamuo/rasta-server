package service

import (
	"sync"

	"github.com/google/uuid"
	"github.com/nkamuo/rasta-server/model"
	"github.com/nkamuo/rasta-server/repository"
)

var orderItemService OrderItemService
var orderItemRepoMutext *sync.Mutex = &sync.Mutex{}

func GetOrderItemService() OrderItemService {
	orderItemRepoMutext.Lock()
	if orderItemService == nil {
		orderItemService = &orderItemServiceImpl{repo: repository.GetOrderItemRepository()}
	}
	orderItemRepoMutext.Unlock()
	return orderItemService
}

type OrderItemService interface {
	GetById(id uuid.UUID) (orderItem *model.OrderItem, err error)
	// GetByEmail(email string) (orderItem *model.OrderItem, err error)
	// GetByPhone(phone string) (orderItem *model.OrderItem, err error)
	Save(orderItem *model.OrderItem) (err error)
	Delete(orderItem *model.OrderItem) (error error)
}

type orderItemServiceImpl struct {
	repo repository.OrderItemRepository
}

func (service *orderItemServiceImpl) GetById(id uuid.UUID) (orderItem *model.OrderItem, err error) {
	return service.repo.GetById(id)
}

func (service *orderItemServiceImpl) Save(orderItem *model.OrderItem) (err error) {
	return service.repo.Save(orderItem)
}

func (service *orderItemServiceImpl) Delete(orderItem *model.OrderItem) (err error) {
	err = service.repo.Delete(orderItem)

	return err
}

func (service *orderItemServiceImpl) DeleteById(id uuid.UUID) (orderItem *model.OrderItem, err error) {
	orderItem, err = service.repo.DeleteById(id)
	return orderItem, err
}
