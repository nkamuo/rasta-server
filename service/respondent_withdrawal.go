package service

import (
	"errors"
	"fmt"
	"sync"

	"github.com/google/uuid"
	"github.com/nkamuo/rasta-server/model"
	"github.com/nkamuo/rasta-server/repository"
	"gorm.io/gorm"
)

var respondentWithdrawalService RespondentWithdrawalService
var respondentWithdrawalRepoMutext *sync.Mutex = &sync.Mutex{}

func GetRespondentWithdrawalService() RespondentWithdrawalService {
	respondentWithdrawalRepoMutext.Lock()
	if respondentWithdrawalService == nil {
		respondentWithdrawalService = &respondentWithdrawalServiceImpl{repo: repository.GetRespondentWithdrawalRepository()}
	}
	respondentWithdrawalRepoMutext.Unlock()
	return respondentWithdrawalService
}

type RespondentWithdrawalService interface {
	GetById(id uuid.UUID) (respondentWithdrawal *model.RespondentWithdrawal, err error)
	// GetByEmail(email string) (respondentWithdrawal *model.RespondentWithdrawal, err error)
	// GetByPhone(phone string) (respondentWithdrawal *model.RespondentWithdrawal, err error)
	Init(wallet model.RespondentWallet, amount uint64, description string) (withdrawal *model.RespondentWithdrawal, err error)
	//THIS PUTS THE MONEY INTO THE TARGETS WALLET BALANCE
	Commite(withdrawal *model.RespondentWithdrawal) (err error)
	// Revert(withdrawal *model.RespondentWithdrawal) ( err error)
	Save(respondentWithdrawal *model.RespondentWithdrawal) (err error)
	Delete(respondentWithdrawal *model.RespondentWithdrawal) (error error)
}

type respondentWithdrawalServiceImpl struct {
	repo repository.RespondentWithdrawalRepository
}

func (service *respondentWithdrawalServiceImpl) GetById(id uuid.UUID) (respondentWithdrawal *model.RespondentWithdrawal, err error) {
	return service.repo.GetById(id)
}

func (service *respondentWithdrawalServiceImpl) Save(respondentWithdrawal *model.RespondentWithdrawal) (err error) {
	return service.repo.Save(respondentWithdrawal)
}

func (service *respondentWithdrawalServiceImpl) Commite(withdrawal *model.RespondentWithdrawal) (err error) {
	walletRepo := repository.GetRespondentWalletRepository()

	if wallet, err := walletRepo.GetById(*withdrawal.WalletID); err != nil {
		return err
	} else {
		err := model.DB.Transaction(func(tx *gorm.DB) error {
			if err := wallet.CommiteWithdrawal(withdrawal.Withdrawal); err != nil {
				return err
			}
			if err := tx.Save(wallet).Error; err != nil {
				return err
			}
			if err := tx.Save(withdrawal).Error; err != nil {
				return err
			}
			return nil
		})
		if err != nil {
			return err
		}
		return nil
	}
}

func (service *respondentWithdrawalServiceImpl) Init(wallet model.RespondentWallet, amount uint64, description string) (withdrawal *model.RespondentWithdrawal, err error) {
	walletService := GetRespondentWalletService()
	if err = walletService.Refresh(&wallet); err != nil {
		return nil, err
	}

	if amount > wallet.Balance { // YOU MIGHT USE wallet.CalculateResultantBalance() for less freedom
		fBalance := float64(wallet.Balance) / 100.0
		fAmount := float64(amount) / 100.0
		message := fmt.Sprintf("You cannot withdraw %2f from your wallet when your balance is %2f", fAmount, fBalance)
		return nil, errors.New(message)
	}
	withdrawal = &model.RespondentWithdrawal{
		Withdrawal: model.Withdrawal{
			Amount:      amount,
			Description: description,
		},
		WalletID: &wallet.ID,
	}

	err = model.DB.Transaction(func(tx *gorm.DB) error {

		if err := tx.Save(withdrawal).Error; err != nil {
			return err
		}
		if err := wallet.InitWithdrawal(withdrawal.Withdrawal); err != nil {
			return err
		}
		if err := tx.Save(wallet).Error; err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return withdrawal, err
}

func (service *respondentWithdrawalServiceImpl) Delete(respondentWithdrawal *model.RespondentWithdrawal) (err error) {
	err = service.repo.Delete(respondentWithdrawal)

	return err
}

func (service *respondentWithdrawalServiceImpl) DeleteById(id uuid.UUID) (respondentWithdrawal *model.RespondentWithdrawal, err error) {
	respondentWithdrawal, err = service.repo.DeleteById(id)
	return respondentWithdrawal, err
}
