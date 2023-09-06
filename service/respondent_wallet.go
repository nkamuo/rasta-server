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
		respondentWalletService = &respondentWalletServiceImpl{repo: repository.GetRespondentWalletRepository()}
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
			}
		}
	} else {
		return err
	}

	if commitedDebit > committedCredit {
		message := fmt.Sprintf("State Error: Wallet balance can't be negetive - commited debit is greater than commited credit")
		return errors.New(message)
	}

	balance = committedCredit - commitedDebit

	wallet.Balance = balance
	wallet.PendingDebitTotal = pendingDebit
	wallet.PendingCreditTotal = pendingCredit

	return nil
}

func (service *respondentWalletServiceImpl) DeleteById(id uuid.UUID) (respondentWallet *model.RespondentWallet, err error) {
	respondentWallet, err = service.repo.DeleteById(id)
	return respondentWallet, err
}
