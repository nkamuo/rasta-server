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

var companyWithdrawalService CompanyWithdrawalService
var companyWithdrawalRepoMutext *sync.Mutex = &sync.Mutex{}

func GetCompanyWithdrawalService() CompanyWithdrawalService {
	companyWithdrawalRepoMutext.Lock()
	if companyWithdrawalService == nil {
		companyWithdrawalService = &companyWithdrawalServiceImpl{repo: repository.GetCompanyWithdrawalRepository()}
	}
	companyWithdrawalRepoMutext.Unlock()
	return companyWithdrawalService
}

type CompanyWithdrawalService interface {
	GetById(id uuid.UUID) (companyWithdrawal *model.CompanyWithdrawal, err error)
	// GetByEmail(email string) (companyWithdrawal *model.CompanyWithdrawal, err error)
	// GetByPhone(phone string) (companyWithdrawal *model.CompanyWithdrawal, err error)
	Init(wallet model.CompanyWallet, amount uint64, description string) (withdrawal *model.CompanyWithdrawal, err error)
	//THIS PUTS THE MONEY INTO THE TARGETS WALLET BALANCE
	Commite(withdrawal *model.CompanyWithdrawal) (err error)
	// Revert(withdrawal *model.CompanyWithdrawal) ( err error)
	Save(companyWithdrawal *model.CompanyWithdrawal) (err error)
	Delete(companyWithdrawal *model.CompanyWithdrawal) (error error)
}

type companyWithdrawalServiceImpl struct {
	repo repository.CompanyWithdrawalRepository
}

func (service *companyWithdrawalServiceImpl) GetById(id uuid.UUID) (companyWithdrawal *model.CompanyWithdrawal, err error) {
	return service.repo.GetById(id)
}

func (service *companyWithdrawalServiceImpl) Save(companyWithdrawal *model.CompanyWithdrawal) (err error) {
	return service.repo.Save(companyWithdrawal)
}

func (service *companyWithdrawalServiceImpl) Commite(withdrawal *model.CompanyWithdrawal) (err error) {
	walletRepo := repository.GetCompanyWalletRepository()

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

func (service *companyWithdrawalServiceImpl) Init(wallet model.CompanyWallet, amount uint64, description string) (withdrawal *model.CompanyWithdrawal, err error) {
	walletService := GetCompanyWalletService()
	if err = walletService.Refresh(&wallet); err != nil {
		return nil, err
	}

	if amount > wallet.Balance { // YOU MIGHT USE wallet.CalculateResultantBalance() for less freedom
		fBalance := float64(wallet.Balance) / 100.0
		fAmount := float64(amount) / 100.0
		message := fmt.Sprintf("You cannot withdraw %2f from your wallet when your balance is %2f", fAmount, fBalance)
		return nil, errors.New(message)
	}
	withdrawal = &model.CompanyWithdrawal{
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

func (service *companyWithdrawalServiceImpl) Delete(companyWithdrawal *model.CompanyWithdrawal) (err error) {
	err = service.repo.Delete(companyWithdrawal)

	return err
}

func (service *companyWithdrawalServiceImpl) DeleteById(id uuid.UUID) (companyWithdrawal *model.CompanyWithdrawal, err error) {
	companyWithdrawal, err = service.repo.DeleteById(id)
	return companyWithdrawal, err
}
