package service

import (
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/nkamuo/rasta-server/dto"
	"github.com/nkamuo/rasta-server/model"
	"github.com/nkamuo/rasta-server/repository"
	"gorm.io/gorm"
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

	CompleteOrder(order *model.Order, isAuto bool, feedback *dto.ClientOrderConfirmationRequest) (err error)
	Process(order *model.Order) (err error)
	UpdateResponderLocationEntry(order *model.Order, locationEntry model.RespondentSessionLocationEntry) (err error)
	AssignResponder(order *model.Order, session *model.RespondentSession) (err error)
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

func (service *orderServiceImpl) AssignResponder(order *model.Order, session *model.RespondentSession) (err error) {
	respondentService := GetRespondentService()
	fulfilmentService := GetOrderFulfilmentService()

	var fulfilment *model.OrderFulfilment

	respondent, err := respondentService.GetById(*session.RespondentID)
	if err != nil {
		message := fmt.Sprintf("Error loading respondent: %s", err.Error())
		return errors.New(message)
	}

	if order.FulfilmentID != nil {
		if fulfilment, err = fulfilmentService.GetById(*order.FulfilmentID); err != nil {
			if err.Error() != "record not found" {
				return err
			}
		} else {
			if (fulfilment.SessionID != nil && *fulfilment.SessionID == session.ID) || (fulfilment.ResponderID != nil && *fulfilment.ResponderID == respondent.ID) {

			} else {
				return errors.New("Cannot overidde order responder assignment")
			}
		}
	}

	if fulfilment != nil && fulfilment.ResponderID.String() != respondent.ID.String() {
		if err := fulfilmentService.Delete(fulfilment); err != nil {
			return err
		}
		fulfilment = nil
	}

	if fulfilment == nil {
		fulfilment = &model.OrderFulfilment{}
	}
	fulfilment.SessionID = &session.ID
	fulfilment.ResponderID = &respondent.ID
	order.Status = model.ORDER_STATUS_RESPONDENT_ASSIGNED

	err = model.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Save(fulfilment).Error; err != nil {
			return err
		}
		order.FulfilmentID = &fulfilment.ID
		if err := tx.Save(order).Error; err != nil {
			return err
		}
		return nil
	})
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

func (service *orderServiceImpl) UpdateResponderLocationEntry(order *model.Order, locationEntry model.RespondentSessionLocationEntry) (err error) {

	// fulfilmentService := GetOrderFulfilmentService()
	// if order.FulfilmentID == nil {
	// 	return
	// }

	// fulfilment, err := fulfilmentService.GetById(*order.FulfilmentID)
	// if err != nil {
	// 	return err
	// }

	// fulfilment.Coordinates = &locationEntry.Coordinates

	// return fulfilmentService.Save(fulfilment)
	return nil
}

func (service *orderServiceImpl) Delete(order *model.Order) (err error) {
	err = service.repo.Delete(order)

	return err
}

func (service *orderServiceImpl) CompleteOrder(order *model.Order, isAuto bool, feedback *dto.ClientOrderConfirmationRequest) (err error) {

	fulfilmentService := GetOrderFulfilmentService()
	earningService := GetGeneralEarningService()
	// orderService := GetOrderService()

	if order == nil {
		return errors.New("Order is nil or not fulfilled")
	}
	fulfilment, err := fulfilmentService.GetById(*order.FulfilmentID)
	if err != nil {
		return err
	}

	if fulfilment.IsComplete() {
		// return errors.New("Order is already complete")
	} else {

		now := time.Now()
		if isAuto {
			fulfilment.AutoConfirmedAt = &now
		} else {
			fulfilment.ClientConfirmedAt = &now
		}

		err := model.DB.Transaction(func(tx *gorm.DB) error {
			if err := tx.Save(fulfilment).Error; err != nil {
				return err
			}
			if feedback != nil {
				review := &model.RespondentServiceReview{
					RespondentID: fulfilment.ResponderID,
					OrderID:      &order.ID,
					Rating:       feedback.Rating,
					Description:  feedback.Description,
				}
				if err := tx.Save(review).Error; err != nil {
					return err
				}
			}

			return nil
		})

		if err != nil {
			return err
		}

	}

	if err := earningService.ProcessEarnings(order); err != nil {
		return err
	}

	return err
}

func (service *orderServiceImpl) DeleteById(id uuid.UUID) (order *model.Order, err error) {
	order, err = service.repo.DeleteById(id)
	return order, err
}
