package service

import (
	"errors"
	"fmt"
	"sync"

	"github.com/google/uuid"
	"github.com/nkamuo/rasta-server/model"
	"github.com/nkamuo/rasta-server/repository"
)

var generalEarningService GeneralEarningService
var generalEarningRepoMutext *sync.Mutex = &sync.Mutex{}

func GetGeneralEarningService() GeneralEarningService {
	generalEarningRepoMutext.Lock()
	if generalEarningService == nil {
		generalEarningService = &generalEarningServiceImpl{repo: repository.GetGeneralEarningRepository()}
	}
	generalEarningRepoMutext.Unlock()
	return generalEarningService
}

type GeneralEarningService interface {
	GetById(id uuid.UUID) (generalEarning *model.OrderEarning, err error)
	Save(generalEarning *model.OrderEarning) (err error)
	ProcessEarnings(order *model.Order) (err error)
	Delete(generalEarning *model.OrderEarning) (error error)
}

type generalEarningServiceImpl struct {
	repo repository.GeneralEarningRepository
}

func (service *generalEarningServiceImpl) GetById(id uuid.UUID) (generalEarning *model.OrderEarning, err error) {
	return service.repo.GetById(id)
}

func (service *generalEarningServiceImpl) Save(generalEarning *model.OrderEarning) (err error) {
	return service.repo.Save(generalEarning)
}

func (service *generalEarningServiceImpl) Delete(generalEarning *model.OrderEarning) (err error) {
	err = service.repo.Delete(generalEarning)

	return err
}

func (service *generalEarningServiceImpl) DeleteById(id uuid.UUID) (generalEarning *model.OrderEarning, err error) {
	generalEarning, err = service.repo.DeleteById(id)
	return generalEarning, err
}

func (service *generalEarningServiceImpl) ProcessEarnings(order *model.Order) (err error) {

	companyEarningRepository := repository.GetCompanyEarningRepository()
	respondentEarningRepository := repository.GetRespondentEarningRepository()
	responderService := GetRespondentService()
	fulfilmentService := GetOrderFulfilmentService()
	orderRepository := repository.GetOrderRepository()

	if order.FulfilmentID == nil {
		return errors.New(fmt.Sprint("Order fulfilment not started yet"))
	}

	fulfilment, err := fulfilmentService.GetById(*order.FulfilmentID, "Responder.Company", "Session.Respondent.Company")
	if err != nil {
		return errors.New(fmt.Sprintf("Could not load fulfilment Plan: %s", err.Error()))
	}

	var companyEarnings = make([]model.CompanyEarning, 1)
	var respondentEarnings = make([]model.RespondentEarning, 1)

	if fulfilment.ResponderID == nil {
		return errors.New("Order is not assigned to responder yet")
	}

	responder, err := responderService.GetById(*fulfilment.ResponderID)
	if err != nil {
		return err
	}

	if !fulfilment.IsComplete() {
		return errors.New("Order is not completed yet")
	}

	fullOrder, err := orderRepository.GetById(order.ID, "Items.Product")
	if err != nil {
		return err
	}

	if !responder.Independent() {
		for _, request := range *fullOrder.Items {
			earning, err := companyEarningRepository.GetByRequest(request)
			if err != nil {
				if err.Error() != "record not found" {
					return err
				}
			} else {
				continue //EARNING ALREADY CREATED
			}
			earning, err = buildCompanyEarning(*fulfilment, request, *responder.Company)
			if err != nil {
				message := fmt.Sprintf("Error creating earning: %s", err.Error())
				return errors.New(message)
			}
			companyEarnings = append(companyEarnings, *earning)
		}
	} else {
		for _, request := range *fullOrder.Items {
			earning, err := respondentEarningRepository.GetByRequest(request)
			if err != nil {
				if err.Error() != "record not found" {
					return err
				}
			} else {
				continue //EARNING ALREADY CREATED
			}

			earning, err = buildRespondentEarning(*fulfilment, request)
			if err != nil {
				message := fmt.Sprintf("Error creating earning: %s", err.Error())
				return errors.New(message)
			}
			respondentEarnings = append(respondentEarnings, *earning)
		}
	}

	return nil
	// return service.repo.Save(order)
}

func buildCompanyEarning(fulfilment model.OrderFulfilment, request model.Request, company model.Company) (earning *model.CompanyEarning, err error) {

	// if nil == company {
	// 	message := fmt.Sprintf("Could not find the Company this responder[%s] is associated with", fulfilment.Responder.User.FullName())
	// 	return nil, errors.New(message)
	// }

	earningService := GetCompanyEarningService()

	respondent := fulfilment.Responder

	baseEarning, err := buildEarning(fulfilment, request)
	if err != nil {
		return nil, err
	}
	earning = &model.CompanyEarning{
		OrderEarning: *baseEarning,
		CompanyID:    &company.ID,
		RespondentID: &respondent.ID,
	}

	if err := earningService.Save(earning); err != nil {
		return nil, err
	}

	return earning, err
}

func buildRespondentEarning(fulfilment model.OrderFulfilment, request model.Request) (earning *model.RespondentEarning, err error) {

	earningService := GetRespondentEarningService()

	respondent := fulfilment.Responder

	baseEarning, err := buildEarning(fulfilment, request)
	if err != nil {
		return nil, err
	}
	earning = &model.RespondentEarning{
		OrderEarning: *baseEarning,
		RespondentID: &respondent.ID,
	}

	if err := earningService.Save(earning); err != nil {
		return nil, err
	}

	return earning, err
}

func buildEarning(fulfilment model.OrderFulfilment, request model.Request) (earning *model.OrderEarning, err error) {
	amount := request.GetAmount()
	label := fmt.Sprintf("%s", request.Product.Label)
	description := fmt.Sprintf("Reward for %s", request.Product.Label)
	earning = &model.OrderEarning{
		Amount:      amount,
		Label:       label,
		Description: description,
		RequestID:   &request.ID,
		Status:      model.ORDER_EARNING_STATUS_PENDING,
	}
	return earning, err
}
