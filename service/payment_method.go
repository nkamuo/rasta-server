package service

import (
	"sync"

	"github.com/google/uuid"
	"github.com/nkamuo/rasta-server/model"
	"github.com/nkamuo/rasta-server/repository"
)

var paymentMethodService PaymentMethodService
var paymentMethodRepoMutext *sync.Mutex = &sync.Mutex{}

func GetPaymentMethodService() PaymentMethodService {
	paymentMethodRepoMutext.Lock()
	if paymentMethodService == nil {
		paymentMethodService = &paymentMethodServiceImpl{repo: repository.GetPaymentMethodRepository()}
	}
	paymentMethodRepoMutext.Unlock()
	return paymentMethodService
}

type PaymentMethodService interface {
	GetById(id uuid.UUID) (paymentMethod *model.PaymentMethod, err error)
	// GetByEmail(email string) (paymentMethod *model.PaymentMethod, err error)
	// GetByPhone(phone string) (paymentMethod *model.PaymentMethod, err error)
	Save(paymentMethod *model.PaymentMethod) (err error)
	Delete(paymentMethod *model.PaymentMethod) (error error)
}

type paymentMethodServiceImpl struct {
	repo repository.PaymentMethodRepository
}

func (service *paymentMethodServiceImpl) GetById(id uuid.UUID) (paymentMethod *model.PaymentMethod, err error) {
	return service.repo.GetById(id)
}

func (service *paymentMethodServiceImpl) Save(paymentMethod *model.PaymentMethod) (err error) {
	return service.repo.Save(paymentMethod)
}

func (service *paymentMethodServiceImpl) Delete(paymentMethod *model.PaymentMethod) (err error) {
	err = service.repo.Delete(paymentMethod)

	return err
}

func (service *paymentMethodServiceImpl) DeleteById(id uuid.UUID) (paymentMethod *model.PaymentMethod, err error) {
	paymentMethod, err = service.repo.DeleteById(id)
	return paymentMethod, err
}
