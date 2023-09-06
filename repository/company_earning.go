package repository

import (
	"sync"

	"github.com/google/uuid"
	"github.com/nkamuo/rasta-server/model"
	"gorm.io/gorm"
)

var companyEarningRepo CompanyEarningRepository
var companyEarningRepoMutext *sync.Mutex = &sync.Mutex{}

func GetCompanyEarningRepository() CompanyEarningRepository {
	companyEarningRepoMutext.Lock()
	if companyEarningRepo == nil {
		companyEarningRepo = &companyEarningRepository{db: model.DB}
	}
	companyEarningRepoMutext.Unlock()
	return companyEarningRepo
}

type CompanyEarningRepository interface {
	FindAll(page int, limit int) (companyEarnings []model.CompanyEarning, total int64, err error)
	GetById(id uuid.UUID) (companyEarning *model.CompanyEarning, err error)
	GetByRequest(request model.Request) (companyEarning *model.CompanyEarning, err error)
	FindByCompany(company model.Company) (withdrawals *[]model.CompanyEarning, err error)
	Save(companyEarning *model.CompanyEarning) (err error)
	Delete(companyEarning *model.CompanyEarning) (error error)
	DeleteById(id uuid.UUID) (companyEarning *model.CompanyEarning, err error)
}

type companyEarningRepository struct {
	db *gorm.DB
}

func (repo *companyEarningRepository) FindAll(page int, limit int) (companyEarnings []model.CompanyEarning, total int64, err error) {
	offset := (page - 1) * limit

	err = repo.db.
		Model(&model.CompanyEarning{}).
		Count(&total).
		Limit(limit).
		Offset(offset).
		Order("created_at desc").
		Find(&companyEarnings).Error
	if err != nil {
		return
	}
	return
}

func (repo *companyEarningRepository) GetById(id uuid.UUID) (companyEarning *model.CompanyEarning, err error) {
	if err = model.DB. /*.Preload("Place")*/ Where("id = ?", id).First(&companyEarning).Error; err != nil {
		return nil, err
	}
	return companyEarning, nil
}

func (repo *companyEarningRepository) FindByCompany(company model.Company) (withdrawals *[]model.CompanyEarning, err error) {
	if err = model.DB.Where("company_id = ?", company.ID).Find(&withdrawals).Error; err != nil {
		return nil, err
	}
	return withdrawals, nil
}

func (repo *companyEarningRepository) GetByRequest(request model.Request) (companyEarning *model.CompanyEarning, err error) {
	if err = model.DB.Where("request_id = ?", request.ID).First(&companyEarning).Error; err != nil {
		return nil, err
	}
	return companyEarning, nil
}

func (repo *companyEarningRepository) Save(companyEarning *model.CompanyEarning) (err error) {
	if (uuid.UUID{} == companyEarning.ID) {
		//NEW - No ID yet
		return repo.db.Create(&companyEarning).Error
	}
	return repo.db.Updates(&companyEarning).Error
}

func (repo *companyEarningRepository) Delete(companyEarning *model.CompanyEarning) (err error) {
	return repo.db.Delete(&companyEarning).Error
}

func (repo *companyEarningRepository) DeleteById(id uuid.UUID) (companyEarning *model.CompanyEarning, err error) {
	companyEarning, err = repo.GetById(id)
	if err != nil {
		return nil, err
	}
	err = repo.db.Delete(companyEarning).Error
	return companyEarning, err
}
