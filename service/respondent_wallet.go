package service

import (
	"sync"

	"github.com/google/uuid"
	"github.com/nkamuo/rasta-server/model"
	"github.com/nkamuo/rasta-server/repository"
)

var respondentWalletService RespondentWalletService
var respondentWalletRepoMutext *sync.Mutex = &sync.Mutex{}

func GetRespondentWalletService() RespondentWalletService {
	respondentWalletRepoMutext.Lock()
	if respondentWalletService == nil {
		respondentWalletService = &respondentWalletServiceImpl{repo: repository.GetRespondentWalletRepository()}
	}
	respondentWalletRepoMutext.Unlock()
	return respondentWalletService
}

type RespondentWalletService interface {
	GetById(id uuid.UUID) (respondentWallet *model.RespondentWallet, err error)
	// GetByEmail(email string) (respondentWallet *model.RespondentWallet, err error)
	// GetByPhone(phone string) (respondentWallet *model.RespondentWallet, err error)
	Save(respondentWallet *model.RespondentWallet) (err error)
	Delete(respondentWallet *model.RespondentWallet) (error error)
}

type respondentWalletServiceImpl struct {
	repo repository.RespondentWalletRepository
}

func (service *respondentWalletServiceImpl) GetById(id uuid.UUID) (respondentWallet *model.RespondentWallet, err error) {
	return service.repo.GetById(id)
}

func (service *respondentWalletServiceImpl) Save(respondentWallet *model.RespondentWallet) (err error) {
	return service.repo.Save(respondentWallet)
}

func (service *respondentWalletServiceImpl) Delete(respondentWallet *model.RespondentWallet) (err error) {
	err = service.repo.Delete(respondentWallet)

	return err
}

func (service *respondentWalletServiceImpl) DeleteById(id uuid.UUID) (respondentWallet *model.RespondentWallet, err error) {
	respondentWallet, err = service.repo.DeleteById(id)
	return respondentWallet, err
}
