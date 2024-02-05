package service

import (
	"sync"

	"github.com/google/uuid"
	"github.com/nkamuo/rasta-server/model"
	"github.com/nkamuo/rasta-server/repository"
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
	// Close(balance *model.RespondentAccessProductBalance) (err error)
	Save(balance *model.RespondentAccessProductBalance) (err error)
	Delete(balance *model.RespondentAccessProductBalance) (error error)
}

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
