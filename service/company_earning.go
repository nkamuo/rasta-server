package service

import (
	"sync"

	"github.com/google/uuid"
	"github.com/nkamuo/rasta-server/model"
	"github.com/nkamuo/rasta-server/repository"
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

func (service *companyEarningServiceImpl) Delete(companyEarning *model.CompanyEarning) (err error) {
	err = service.repo.Delete(companyEarning)

	return err
}

func (service *companyEarningServiceImpl) DeleteById(id uuid.UUID) (companyEarning *model.CompanyEarning, err error) {
	companyEarning, err = service.repo.DeleteById(id)
	return companyEarning, err
}
