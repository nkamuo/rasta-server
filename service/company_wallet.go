package service

import (
	"errors"
	"fmt"
	"sync"

	"github.com/google/uuid"
	"github.com/nkamuo/rasta-server/model"
	"github.com/nkamuo/rasta-server/repository"
)

var companyWalletService CompanyWalletService
var companyWalletRepoMutext *sync.Mutex = &sync.Mutex{}

func GetCompanyWalletService() CompanyWalletService {
	companyWalletRepoMutext.Lock()
	if companyWalletService == nil {
		companyWalletService = &companyWalletServiceImpl{repo: repository.GetCompanyWalletRepository()}
	}
	companyWalletRepoMutext.Unlock()
	return companyWalletService
}

type CompanyWalletService interface {
	GetById(id uuid.UUID) (companyWallet *model.CompanyWallet, err error)
	// GetByEmail(email string) (companyWallet *model.CompanyWallet, err error)
	// GetByPhone(phone string) (companyWallet *model.CompanyWallet, err error)
	CreateNewFor(company *model.Company) (err error)
	Refresh(wallet *model.CompanyWallet) (err error)
	Save(companyWallet *model.CompanyWallet) (err error)
	Delete(companyWallet *model.CompanyWallet) (error error)
}

type companyWalletServiceImpl struct {
	repo repository.CompanyWalletRepository
}

func (service *companyWalletServiceImpl) GetById(id uuid.UUID) (companyWallet *model.CompanyWallet, err error) {
	return service.repo.GetById(id)
}

func (service *companyWalletServiceImpl) Save(companyWallet *model.CompanyWallet) (err error) {
	return service.repo.Save(companyWallet)
}

func (service *companyWalletServiceImpl) Delete(companyWallet *model.CompanyWallet) (err error) {
	err = service.repo.Delete(companyWallet)

	return err
}

func (service *companyWalletServiceImpl) CreateNewFor(company *model.Company) (err error) {
	if _, err := service.repo.GetByCompany(*company); err != nil {
		if err.Error() != "record not found" {
			return err
		}
	}
	wallet := &model.CompanyWallet{
		Wallet: model.Wallet{
			Status:             model.WALLET_STATUS_ACTIVE,
			Balance:            0,
			PendingCreditTotal: 0,
			PendingDebitTotal:  0,
		},
		CompanyID: &company.ID,
	}
	return service.Save(wallet)
}

func (service *companyWalletServiceImpl) DeleteById(id uuid.UUID) (companyWallet *model.CompanyWallet, err error) {
	companyWallet, err = service.repo.DeleteById(id)
	return companyWallet, err
}

func (repo *companyWalletServiceImpl) Refresh(wallet *model.CompanyWallet) (err error) {

	companyService := GetCompanyService()
	earningRepo := repository.GetCompanyEarningRepository()
	withdrawalRepo := repository.GetCompanyWithdrawalRepository()

	var balance, pendingCredit, committedCredit, pendingDebit, commitedDebit uint64

	company, err := companyService.GetById(*wallet.CompanyID)
	if err != nil {
		return err
	}

	if earnings, err := earningRepo.FindByCompany(*company); err == nil {
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

	if withdrawals, err := withdrawalRepo.FindByCompany(*company); err == nil {
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
