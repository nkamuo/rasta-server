package service

import (
	"errors"
	"fmt"
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
		respondentWalletService = &respondentWalletServiceImpl{
			repo:          repository.GetRespondentWalletRepository(),
			walletMutexes: map[string]*sync.RWMutex{},
		}
	}
	respondentWalletRepoMutext.Unlock()
	return respondentWalletService
}

type RespondentWalletService interface {
	GetById(id uuid.UUID) (respondentWallet *model.RespondentWallet, err error)
	// GetByEmail(email string) (respondentWallet *model.RespondentWallet, err error)
	// GetByPhone(phone string) (respondentWallet *model.RespondentWallet, err error)

	CreateNewFor(respondent *model.Respondent) (err error)
	Refresh(wallet *model.RespondentWallet) (err error)
	GetMutex(wallet *model.RespondentWallet) (mutext *sync.RWMutex)
	Save(respondentWallet *model.RespondentWallet) (err error)
	Delete(respondentWallet *model.RespondentWallet) (error error)
}

type respondentWalletServiceImpl struct {
	repo          repository.RespondentWalletRepository
	walletMutexes map[string]*sync.RWMutex // = map[string]sync.RWMutex{}
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

func (service *respondentWalletServiceImpl) GetMutex(wallet *model.RespondentWallet) (mutext *sync.RWMutex) {
	walletId := wallet.ID.String()

	if walletId == "" {
		panic("CANNOT GET MUTEX FOR  WALLET WITHOUT A VALID ID")
	}

	mutex, ok := service.walletMutexes[walletId]
	if ok == false {
		mutex = &sync.RWMutex{}
		service.walletMutexes[walletId] = mutex
	}
	return mutex
}

func (service *respondentWalletServiceImpl) CreateNewFor(respondent *model.Respondent) (err error) {
	if _, err := service.repo.GetByRespondent(*respondent); err != nil {
		if err.Error() != "record not found" {
			return err
		}
	}
	wallet := &model.RespondentWallet{
		Wallet: model.Wallet{
			Status:             model.WALLET_STATUS_ACTIVE,
			Balance:            0,
			PendingCreditTotal: 0,
			PendingDebitTotal:  0,
		},
		RespondentID: &respondent.ID,
	}
	return service.Save(wallet)
}

func (service *respondentWalletServiceImpl) Refresh(wallet *model.RespondentWallet) (err error) {

	respondentService := GetRespondentService()
	earningRepo := repository.GetRespondentEarningRepository()
	withdrawalRepo := repository.GetRespondentWithdrawalRepository()

	var balance, pendingCredit, committedCredit, pendingDebit, commitedDebit uint64

	respondent, err := respondentService.GetById(*wallet.RespondentID)
	if err != nil {
		return err
	}

	if earnings, err := earningRepo.FindByRespondent(*respondent); err == nil {
		for _, earning := range *earnings {
			if earning.IsCommited() {
				committedCredit += earning.Amount
			} else if earning.IsPending() {
				pendingCredit += earning.Amount
			}
		}

	} else {
		return err
	}

	if withdrawals, err := withdrawalRepo.FindByRespondent(*respondent); err == nil {
		for _, withdrawal := range *withdrawals {
			if withdrawal.IsCommited() {
				commitedDebit += withdrawal.Amount
			} else if withdrawal.IsPending() {
				pendingDebit += withdrawal.Amount
				// commitedDebit += withdrawal.Amount
			}
		}
	} else {
		return err
	}

	if commitedDebit > committedCredit {
		message := fmt.Sprintf("State Error: Wallet balance can't be negetive - commited debit is greater than commited credit")
		return errors.New(message)
	}

	balance = committedCredit - commitedDebit - pendingDebit

	wallet.Balance = balance
	wallet.PendingDebitTotal = pendingDebit
	wallet.PendingCreditTotal = pendingCredit

	err = service.Save(wallet)

	return err
}

func (service *respondentWalletServiceImpl) DeleteById(id uuid.UUID) (respondentWallet *model.RespondentWallet, err error) {
	respondentWallet, err = service.repo.DeleteById(id)
	return respondentWallet, err
}
