package service

import (
	"errors"
	"sync"

	"github.com/google/uuid"
	"github.com/nkamuo/rasta-server/model"
	"github.com/nkamuo/rasta-server/repository"
)

var orderService OrderService
var orderRepoMutext *sync.Mutex = &sync.Mutex{}

func GetOrderService() OrderService {
	orderRepoMutext.Lock()
	if orderService == nil {
		orderService = &orderServiceImpl{repo: repository.GetOrderRepository()}
	}
	orderRepoMutext.Unlock()
	return orderService
}

type OrderService interface {
	GetById(id uuid.UUID) (order *model.Order, err error)
	// GetByEmail(email string) (order *model.Order, err error)
	// GetByPhone(phone string) (order *model.Order, err error)

	Process(order *model.Order) (err error)
	AssignResponder(order *model.Order, responder *model.Respondent) (err error)
	Save(order *model.Order) (err error)
	Delete(order *model.Order) (error error)
}

type orderServiceImpl struct {
	repo repository.OrderRepository
}

func (service *orderServiceImpl) GetById(id uuid.UUID) (order *model.Order, err error) {
	return service.repo.GetById(id)
}

func (service *orderServiceImpl) Save(order *model.Order) (err error) {
	return service.repo.Save(order)
}

func (service *orderServiceImpl) AssignResponder(order *model.Order, responder *model.Respondent) (err error) {
	var fulfilment *model.OrderFulfilment

	fulfilmentService := GetOrderFulfilmentService()

	if order.FulfilmentID != nil {
		if fulfilment, err = fulfilmentService.GetById(*order.FulfilmentID); err != nil {
			if err.Error() != "record not found" {
				return err
			}
		}
		return errors.New("Cannot overidde order responder assignment")
	}

	if fulfilment != nil && fulfilment.ResponderID.String() != responder.ID.String() {
		if err := fulfilmentService.Delete(fulfilment); err != nil {
			return err
		}
		fulfilment = nil
	}

	if fulfilment == nil {
		fulfilment = &model.OrderFulfilment{}
	}
	fulfilment.ResponderID = &responder.ID
	order.Status = model.ORDER_STATUS_RESPONDENT_ASSIGNED

	if err := fulfilmentService.Save(fulfilment); err != nil {
		return err
	}

	order.FulfilmentID = &fulfilment.ID

	if err := service.Save(order); err != nil {
		return err
	}
	return err
}

func (service *orderServiceImpl) Process(order *model.Order) (err error) {
	var service_fee_description = "Charge for using HUQT service"
	// itemTotal := order.CalculateItemTotal();
	// adjustmentTotal := order.CalculateAdjustmentTotal();

	var serviceFee *model.OrderAdjustment

	for _, adj := range order.Adjustments {
		const SERVICE_FEE_ADJUSTMENT_CODE = "SERVICE_FEE"
		if adj.Code == SERVICE_FEE_ADJUSTMENT_CODE {
			serviceFee = &adj
			break
		}
	}

	if nil == serviceFee {
		serviceFee = &model.OrderAdjustment{
			Code:        model.SERVICE_FEE_ADJUSTMENT_CODE,
			Title:       "Service Fee",
			Amount:      900,
			Description: &service_fee_description,
		}
		order.AddAdjustment(*serviceFee)
	}

	return nil
}

func (service *orderServiceImpl) Delete(order *model.Order) (err error) {
	err = service.repo.Delete(order)

	return err
}

func (service *orderServiceImpl) DeleteById(id uuid.UUID) (order *model.Order, err error) {
	order, err = service.repo.DeleteById(id)
	return order, err
}
