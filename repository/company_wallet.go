package repository

import (
	"sync"

	"github.com/google/uuid"
	"github.com/nkamuo/rasta-server/model"
	"gorm.io/gorm"
)

var companyWalletRepo CompanyWalletRepository
var companyWalletRepoMutext *sync.Mutex = &sync.Mutex{}

func GetCompanyWalletRepository() CompanyWalletRepository {
	companyWalletRepoMutext.Lock()
	if companyWalletRepo == nil {
		companyWalletRepo = &companyWalletRepository{db: model.DB}
	}
	companyWalletRepoMutext.Unlock()
	return companyWalletRepo
}

type CompanyWalletRepository interface {
	FindAll(page int, limit int) (companyWallets []model.CompanyWallet, total int64, err error)
	GetById(id uuid.UUID) (companyWallet *model.CompanyWallet, err error)
	GetByRequest(request model.Request) (companyWallet *model.CompanyWallet, err error)
	Save(companyWallet *model.CompanyWallet) (err error)
	Delete(companyWallet *model.CompanyWallet) (error error)
	DeleteById(id uuid.UUID) (companyWallet *model.CompanyWallet, err error)
}

type companyWalletRepository struct {
	db *gorm.DB
}

func (repo *companyWalletRepository) FindAll(page int, limit int) (companyWallets []model.CompanyWallet, total int64, err error) {
	offset := (page - 1) * limit

	err = repo.db.
		Model(&model.CompanyWallet{}).
		Count(&total).
		Limit(limit).
		Offset(offset).
		Order("created_at desc").
		Find(&companyWallets).Error
	if err != nil {
		return
	}
	return
}

func (repo *companyWalletRepository) GetById(id uuid.UUID) (companyWallet *model.CompanyWallet, err error) {
	if err = model.DB. /*.Preload("Place")*/ Where("id = ?", id).First(&companyWallet).Error; err != nil {
		return nil, err
	}
	return companyWallet, nil
}

func (repo *companyWalletRepository) GetByRequest(request model.Request) (companyWallet *model.CompanyWallet, err error) {
	if err = model.DB.Where("request_id = ?", request.ID).First(&companyWallet).Error; err != nil {
		return nil, err
	}
	return companyWallet, nil
}

func (repo *companyWalletRepository) Save(companyWallet *model.CompanyWallet) (err error) {
	if (uuid.UUID{} == companyWallet.ID) {
		//NEW - No ID yet
		return repo.db.Create(&companyWallet).Error
	}
	return repo.db.Updates(&companyWallet).Error
}

func (repo *companyWalletRepository) Delete(companyWallet *model.CompanyWallet) (err error) {
	return repo.db.Delete(&companyWallet).Error
}

func (repo *companyWalletRepository) DeleteById(id uuid.UUID) (companyWallet *model.CompanyWallet, err error) {
	companyWallet, err = repo.GetById(id)
	if err != nil {
		return nil, err
	}
	err = repo.db.Delete(companyWallet).Error
	return companyWallet, err
}
