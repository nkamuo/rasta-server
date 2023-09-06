package repository

import (
	"sync"

	"github.com/google/uuid"
	"github.com/nkamuo/rasta-server/model"
	"gorm.io/gorm"
)

var companyWithdrawalRepo CompanyWithdrawalRepository
var companyWithdrawalRepoMutext *sync.Mutex = &sync.Mutex{}

func GetCompanyWithdrawalRepository() CompanyWithdrawalRepository {
	companyWithdrawalRepoMutext.Lock()
	if companyWithdrawalRepo == nil {
		companyWithdrawalRepo = &companyWithdrawalRepository{db: model.DB}
	}
	companyWithdrawalRepoMutext.Unlock()
	return companyWithdrawalRepo
}

type CompanyWithdrawalRepository interface {
	FindAll(page int, limit int) (companyWithdrawals []model.CompanyWithdrawal, total int64, err error)
	GetById(id uuid.UUID) (companyWithdrawal *model.CompanyWithdrawal, err error)
	GetByRequest(request model.Request) (companyWithdrawal *model.CompanyWithdrawal, err error)
	FindByCompany(company model.Company) (companyWithdrawal *[]model.CompanyWithdrawal, err error)
	Save(companyWithdrawal *model.CompanyWithdrawal) (err error)
	Delete(companyWithdrawal *model.CompanyWithdrawal) (error error)
	DeleteById(id uuid.UUID) (companyWithdrawal *model.CompanyWithdrawal, err error)
}

type companyWithdrawalRepository struct {
	db *gorm.DB
}

func (repo *companyWithdrawalRepository) FindAll(page int, limit int) (companyWithdrawals []model.CompanyWithdrawal, total int64, err error) {
	offset := (page - 1) * limit

	err = repo.db.
		Model(&model.CompanyWithdrawal{}).
		Count(&total).
		Limit(limit).
		Offset(offset).
		Order("created_at desc").
		Find(&companyWithdrawals).Error
	if err != nil {
		return
	}
	return
}

func (repo *companyWithdrawalRepository) GetById(id uuid.UUID) (companyWithdrawal *model.CompanyWithdrawal, err error) {
	if err = model.DB. /*.Preload("Place")*/ Where("id = ?", id).First(&companyWithdrawal).Error; err != nil {
		return nil, err
	}
	return companyWithdrawal, nil
}

func (repo *companyWithdrawalRepository) GetByRequest(request model.Request) (companyWithdrawal *model.CompanyWithdrawal, err error) {
	if err = model.DB.Where("request_id = ?", request.ID).First(&companyWithdrawal).Error; err != nil {
		return nil, err
	}
	return companyWithdrawal, nil
}

func (repo *companyWithdrawalRepository) FindByCompany(company model.Company) (withdrawals *[]model.CompanyWithdrawal, err error) {
	walletRepo := GetCompanyWalletRepository()
	wallet, err := walletRepo.GetByCompany(company)
	if err != nil {
		return nil, err
	}
	if err = model.DB.Where("wallet_id = ?", wallet.ID).Find(&withdrawals).Error; err != nil {
		return nil, err
	}
	return withdrawals, nil
}

func (repo *companyWithdrawalRepository) Save(companyWithdrawal *model.CompanyWithdrawal) (err error) {
	if (uuid.UUID{} == companyWithdrawal.ID) {
		//NEW - No ID yet
		return repo.db.Create(&companyWithdrawal).Error
	}
	return repo.db.Updates(&companyWithdrawal).Error
}

func (repo *companyWithdrawalRepository) Delete(companyWithdrawal *model.CompanyWithdrawal) (err error) {
	return repo.db.Delete(&companyWithdrawal).Error
}

func (repo *companyWithdrawalRepository) DeleteById(id uuid.UUID) (companyWithdrawal *model.CompanyWithdrawal, err error) {
	companyWithdrawal, err = repo.GetById(id)
	if err != nil {
		return nil, err
	}
	err = repo.db.Delete(companyWithdrawal).Error
	return companyWithdrawal, err
}
