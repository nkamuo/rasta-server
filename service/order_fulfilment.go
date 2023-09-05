package service

import (
	"sync"

	"github.com/google/uuid"
	"github.com/nkamuo/rasta-server/model"
	"github.com/nkamuo/rasta-server/repository"
	// "github.com/stripe/stripe-go/v74/fulfilmentmethod"
)

var fulfilmentService OrderFulfilmentService
var fulfilmentRepoMutext *sync.Mutex = &sync.Mutex{}

func GetOrderFulfilmentService() OrderFulfilmentService {
	fulfilmentRepoMutext.Lock()
	if fulfilmentService == nil {
		fulfilmentService = &fulfilmentServiceImpl{repo: repository.GetOrderFulfilmentRepository()}
	}
	fulfilmentRepoMutext.Unlock()
	return fulfilmentService
}

type OrderFulfilmentService interface {
	GetById(id uuid.UUID, preload ...string) (fulfilment *model.OrderFulfilment, err error)
	// GetByEmail(email string) (fulfilment *model.OrderFulfilment, err error)
	// GetByPhone(phone string) (fulfilment *model.OrderFulfilment, err error)
	Save(fulfilment *model.OrderFulfilment) (err error)
	Delete(fulfilment *model.OrderFulfilment) (error error)
}

type fulfilmentServiceImpl struct {
	repo repository.OrderFulfilmentRepository
}

func (service *fulfilmentServiceImpl) GetById(id uuid.UUID, preload ...string) (fulfilment *model.OrderFulfilment, err error) {
	return service.repo.GetById(id)
}

func (service *fulfilmentServiceImpl) Save(fulfilment *model.OrderFulfilment) (err error) {
	return service.repo.Save(fulfilment)
}

func (service *fulfilmentServiceImpl) Delete(fulfilment *model.OrderFulfilment) (err error) {
	err = service.repo.Delete(fulfilment)

	return err
}

func (service *fulfilmentServiceImpl) DeleteById(id uuid.UUID) (fulfilment *model.OrderFulfilment, err error) {
	fulfilment, err = service.repo.DeleteById(id)
	return fulfilment, err
}
