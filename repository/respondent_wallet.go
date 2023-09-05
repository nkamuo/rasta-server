package repository

import (
	"sync"

	"github.com/google/uuid"
	"github.com/nkamuo/rasta-server/model"
	"gorm.io/gorm"
)

var respondentWalletRepo RespondentWalletRepository
var respondentWalletRepoMutext *sync.Mutex = &sync.Mutex{}

func GetRespondentWalletRepository() RespondentWalletRepository {
	respondentWalletRepoMutext.Lock()
	if respondentWalletRepo == nil {
		respondentWalletRepo = &respondentWalletRepository{db: model.DB}
	}
	respondentWalletRepoMutext.Unlock()
	return respondentWalletRepo
}

type RespondentWalletRepository interface {
	FindAll(page int, limit int) (respondentWallets []model.RespondentWallet, total int64, err error)
	GetById(id uuid.UUID) (respondentWallet *model.RespondentWallet, err error)
	GetByRequest(request model.Request) (respondentWallet *model.RespondentWallet, err error)
	Save(respondentWallet *model.RespondentWallet) (err error)
	Delete(respondentWallet *model.RespondentWallet) (error error)
	DeleteById(id uuid.UUID) (respondentWallet *model.RespondentWallet, err error)
}

type respondentWalletRepository struct {
	db *gorm.DB
}

func (repo *respondentWalletRepository) FindAll(page int, limit int) (respondentWallets []model.RespondentWallet, total int64, err error) {
	offset := (page - 1) * limit

	err = repo.db.
		Model(&model.RespondentWallet{}).
		Count(&total).
		Limit(limit).
		Offset(offset).
		Order("created_at desc").
		Find(&respondentWallets).Error
	if err != nil {
		return
	}
	return
}

func (repo *respondentWalletRepository) GetById(id uuid.UUID) (respondentWallet *model.RespondentWallet, err error) {
	if err = model.DB. /*.Preload("Place")*/ Where("id = ?", id).First(&respondentWallet).Error; err != nil {
		return nil, err
	}
	return respondentWallet, nil
}

func (repo *respondentWalletRepository) GetByRequest(request model.Request) (respondentWallet *model.RespondentWallet, err error) {
	if err = model.DB.Where("request_id = ?", request.ID).First(&respondentWallet).Error; err != nil {
		return nil, err
	}
	return respondentWallet, nil
}

func (repo *respondentWalletRepository) Save(respondentWallet *model.RespondentWallet) (err error) {
	if (uuid.UUID{} == respondentWallet.ID) {
		//NEW - No ID yet
		return repo.db.Create(&respondentWallet).Error
	}
	return repo.db.Updates(&respondentWallet).Error
}

func (repo *respondentWalletRepository) Delete(respondentWallet *model.RespondentWallet) (err error) {
	return repo.db.Delete(&respondentWallet).Error
}

func (repo *respondentWalletRepository) DeleteById(id uuid.UUID) (respondentWallet *model.RespondentWallet, err error) {
	respondentWallet, err = repo.GetById(id)
	if err != nil {
		return nil, err
	}
	err = repo.db.Delete(respondentWallet).Error
	return respondentWallet, err
}
