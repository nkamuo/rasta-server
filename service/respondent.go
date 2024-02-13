package service

import (
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/nkamuo/rasta-server/model"
	"github.com/nkamuo/rasta-server/repository"
)

var respondentService RespondentService
var respondentRepoMutext *sync.Mutex = &sync.Mutex{}

func GetRespondentService() RespondentService {
	respondentRepoMutext.Lock()
	if respondentService == nil {
		respondentService = &respondentServiceImpl{repo: repository.GetRespondentRepository()}
	}
	respondentRepoMutext.Unlock()
	return respondentService
}

type RespondentService interface {
	GetById(id uuid.UUID, preload ...string) (respondent *model.Respondent, err error)
	Save(respondent *model.Respondent) (err error)
	AssignToCompany(respondent *model.Respondent, company *model.Company) (err error)
	RemoveFromCompany(respondent *model.Respondent, company *model.Company) (err error)
	CanHandleMotoristRequest(respondent *model.Respondent) (result bool, err error)
	Delete(respondent *model.Respondent) (error error)
}

type respondentServiceImpl struct {
	repo repository.RespondentRepository
}

func (service *respondentServiceImpl) GetById(id uuid.UUID, preload ...string) (respondent *model.Respondent, err error) {
	return service.repo.GetById(id, preload...)
}

func (service *respondentServiceImpl) Save(respondent *model.Respondent) (err error) {
	respondentWalletService := GetRespondentWalletService()
	if err := service.repo.Save(respondent); err != nil {
		return err
	}
	if err := respondentWalletService.CreateNewFor(respondent); err != nil {
		return err
	}
	return nil
}

func (service *respondentServiceImpl) AssignToCompany(respondent *model.Respondent, company *model.Company) (err error) {
	respondent.CompanyID = &company.ID
	return service.Save(respondent)
}

func (service *respondentServiceImpl) RemoveFromCompany(respondent *model.Respondent, company *model.Company) (err error) {
	if respondent.CompanyID != &company.ID {
		return errors.New("Respondent is currently not assigned to this company")
	}
	respondent.CompanyID = &uuid.Nil
	return service.Save(respondent)
}

func (service *respondentServiceImpl) ValidateEmail(respondent *model.Respondent) (err error) {
	return nil
}

func (service *respondentServiceImpl) Delete(respondent *model.Respondent) (err error) {
	err = service.repo.Delete(respondent)

	return err
}

func (service *respondentServiceImpl) DeleteById(id uuid.UUID) (respondent *model.Respondent, err error) {
	respondent, err = service.repo.DeleteById(id)
	return respondent, err
}

func (service *respondentServiceImpl) CanHandleMotoristRequest(respondent *model.Respondent) (result bool, err error) {

	respondentService := GetRespondentService()
	accessBalanceService := GetRespondentAccessProductBalanceService()
	subscriptionService := GetRespondentAccessProductSubscriptionService()

	fullRespondent, err := respondentService.GetById(respondent.ID, "User")
	if err != nil {
		message := fmt.Sprintf("An error occured refetching responder: %s", err.Error())
		return false, errors.New(message)
	}

	if balance, err := accessBalanceService.GetByRespondent(respondent); err != nil {
		return false, err
	} else if balance.Balance != nil && *balance.Balance > 0 {
		return true, nil
	}

	if _, err := subscriptionService.GetActiveByRespondentAndTime(*fullRespondent, time.Now()); err != nil {
		return false, err
	} else {

		return true, nil
	}

	// return false, nil
}
