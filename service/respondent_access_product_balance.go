package service

import (
	"sync"

	"github.com/google/uuid"
	"github.com/nkamuo/rasta-server/model"
	"github.com/nkamuo/rasta-server/repository"
	"github.com/stripe/stripe-go/v74"
)

var balanceService RespondentAccessProductBalanceService
var balanceRepoMutext *sync.Mutex = &sync.Mutex{}

func GetRespondentAccessProductBalanceService() RespondentAccessProductBalanceService {
	balanceRepoMutext.Lock()
	if balanceService == nil {
		balanceService = &balanceServiceImpl{repo: repository.GetRespondentAccessProductBalanceRepository()}
	}
	balanceRepoMutext.Unlock()
	return balanceService
}

type RespondentAccessProductBalanceService interface {
	GetById(id uuid.UUID, preload ...string) (balance *model.RespondentAccessProductBalance, err error)
	GetByRespondent(respondent *model.Respondent, preload ...string) (balance *model.RespondentAccessProductBalance, err error)
	// Close(balance *model.RespondentAccessProductBalance) (err error)
	SetupForRespondent(respondent *model.Respondent) (balance *model.RespondentAccessProductBalance, err error)
	SetupForAllRespondents() (err error)
	Save(balance *model.RespondentAccessProductBalance) (err error)
	Delete(balance *model.RespondentAccessProductBalance) (error error)
}

// RespondentAccessProductBalance

type balanceServiceImpl struct {
	repo repository.RespondentAccessProductBalanceRepository
}

func (service *balanceServiceImpl) GetById(id uuid.UUID, preload ...string) (balance *model.RespondentAccessProductBalance, err error) {
	return service.repo.GetById(id, preload...)
}

// func (service *balanceServiceImpl) Close(balance *model.RespondentAccessProductBalance) (err error) {
// 	now := time.Now()
// 	balance.EndedAt = &now
// 	*balance.Active = false
// 	return service.repo.Save(balance)
// }

func (service balanceServiceImpl) GetByRespondent(respondent *model.Respondent, preload ...string) (balance *model.RespondentAccessProductBalance, err error) {
	return service.repo.GetActiveByRespondent(*respondent, preload...)
	// return nil, errors.New("Could not resolve Access Balance for the given responder");
}

func (service balanceServiceImpl) SetupForRespondent(respondent *model.Respondent) (balance *model.RespondentAccessProductBalance, err error) {

	balance, err = service.GetByRespondent(respondent)
	if err != nil {
		if err.Error() == "record not found" {

		} else {
			return nil, err
		}
	}

	if balance == nil {
		balance = &model.RespondentAccessProductBalance{
			RespondentID: &respondent.ID,
			Balance:      stripe.Int64(0),
		}

		if err = service.Save(balance); err != nil {
			return nil, err
		}
	}
	return balance, err
}

func (service balanceServiceImpl) SetupForAllRespondents() (err error) {
	respondantRepo := repository.GetRespondentRepository()
	respondents, _, err := respondantRepo.FindAll(1, 100000)
	if err != nil {
		return err
	}
	for _, respondent := range respondents {
		_, err = service.SetupForRespondent(&respondent)
		if err != nil {
			return err
		}
	}
	return err
}

func (service *balanceServiceImpl) Save(balance *model.RespondentAccessProductBalance) (err error) {
	return service.repo.Save(balance)
}

func (service *balanceServiceImpl) Delete(balance *model.RespondentAccessProductBalance) (err error) {
	err = service.repo.Delete(balance)

	return err
}

func (service *balanceServiceImpl) DeleteById(id uuid.UUID) (balance *model.RespondentAccessProductBalance, err error) {
	balance, err = service.repo.DeleteById(id)
	return balance, err
}
