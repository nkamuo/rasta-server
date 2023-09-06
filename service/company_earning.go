package service

import (
	"sync"

	"github.com/google/uuid"
	"github.com/nkamuo/rasta-server/model"
	"github.com/nkamuo/rasta-server/repository"
	"gorm.io/gorm"
)

var companyEarningService CompanyEarningService
var companyEarningRepoMutext *sync.Mutex = &sync.Mutex{}

func GetCompanyEarningService() CompanyEarningService {
	companyEarningRepoMutext.Lock()
	if companyEarningService == nil {
		companyEarningService = &companyEarningServiceImpl{repo: repository.GetCompanyEarningRepository()}
	}
	companyEarningRepoMutext.Unlock()
	return companyEarningService
}

type CompanyEarningService interface {
	GetById(id uuid.UUID) (companyEarning *model.CompanyEarning, err error)
	// GetByEmail(email string) (companyEarning *model.CompanyEarning, err error)
	// GetByPhone(phone string) (companyEarning *model.CompanyEarning, err error)

	//THIS PUTS THE MONEY INTO THE TARGETS WALLET BALANCE
	Commite(earning *model.CompanyEarning) (err error)
	// Revert(earning *model.CompanyEarning) ( err error)
	Save(companyEarning *model.CompanyEarning) (err error)
	Delete(companyEarning *model.CompanyEarning) (error error)
}

type companyEarningServiceImpl struct {
	repo repository.CompanyEarningRepository
}

func (service *companyEarningServiceImpl) GetById(id uuid.UUID) (companyEarning *model.CompanyEarning, err error) {
	return service.repo.GetById(id)
}

func (service *companyEarningServiceImpl) Save(companyEarning *model.CompanyEarning) (err error) {
	return service.repo.Save(companyEarning)
}

func (service *companyEarningServiceImpl) Commite(earning *model.CompanyEarning) (err error) {
	walletRepo := repository.GetCompanyWalletRepository()
	companyService := GetCompanyService()

	company, err := companyService.GetById(*earning.CompanyID)
	if err != nil {
		return err
	}

	if wallet, err := walletRepo.GetByCompany(*company); err != nil {
		return err
	} else {
		err := model.DB.Transaction(func(tx *gorm.DB) error {
			if err := wallet.CommiteEarning(earning.OrderEarning); err != nil {
				return err
			}
			if err := tx.Save(wallet).Error; err != nil {
				return err
			}
			if err := tx.Save(earning).Error; err != nil {
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

func (service *companyEarningServiceImpl) Delete(companyEarning *model.CompanyEarning) (err error) {
	err = service.repo.Delete(companyEarning)

	return err
}

func (service *companyEarningServiceImpl) DeleteById(id uuid.UUID) (companyEarning *model.CompanyEarning, err error) {
	companyEarning, err = service.repo.DeleteById(id)
	return companyEarning, err
}
