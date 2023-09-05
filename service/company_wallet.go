package service

import (
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

func (service *companyWalletServiceImpl) DeleteById(id uuid.UUID) (companyWallet *model.CompanyWallet, err error) {
	companyWallet, err = service.repo.DeleteById(id)
	return companyWallet, err
}
