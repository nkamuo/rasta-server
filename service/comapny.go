package service

import (
	"sync"

	"github.com/google/uuid"
	"github.com/nkamuo/rasta-server/model"
	"github.com/nkamuo/rasta-server/repository"
)

var companyService CompanyService
var companyRepoMutext *sync.Mutex = &sync.Mutex{}

func GetCompanyService() CompanyService {
	companyRepoMutext.Lock()
	if companyService == nil {
		companyService = &companyServiceImpl{repo: repository.GetCompanyRepository()}
	}
	companyRepoMutext.Unlock()
	return companyService
}

type CompanyService interface {
	GetById(id uuid.UUID) (company *model.Company, err error)
	Save(company *model.Company) (err error)
	Delete(company *model.Company) (error error)
}

type companyServiceImpl struct {
	repo repository.CompanyRepository
}

func (service *companyServiceImpl) GetById(id uuid.UUID) (company *model.Company, err error) {
	return service.repo.GetById(id)
}

func (service *companyServiceImpl) Save(company *model.Company) (err error) {
	return service.repo.Save(company)
}

func (service *companyServiceImpl) ValidateEmail(company *model.Company) (err error) {
	return nil
}

func (service *companyServiceImpl) Delete(company *model.Company) (err error) {
	err = service.repo.Delete(company)

	return err
}

func (service *companyServiceImpl) DeleteById(id uuid.UUID) (company *model.Company, err error) {
	company, err = service.repo.DeleteById(id)
	return company, err
}
