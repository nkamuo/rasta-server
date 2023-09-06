package repository

import (
	"sync"

	"github.com/google/uuid"
	"github.com/nkamuo/rasta-server/model"
	"gorm.io/gorm"
)

var respondentWithdrawalRepo RespondentWithdrawalRepository
var respondentWithdrawalRepoMutext *sync.Mutex = &sync.Mutex{}

func GetRespondentWithdrawalRepository() RespondentWithdrawalRepository {
	respondentWithdrawalRepoMutext.Lock()
	if respondentWithdrawalRepo == nil {
		respondentWithdrawalRepo = &respondentWithdrawalRepository{db: model.DB}
	}
	respondentWithdrawalRepoMutext.Unlock()
	return respondentWithdrawalRepo
}

type RespondentWithdrawalRepository interface {
	FindAll(page int, limit int) (respondentWithdrawals []model.RespondentWithdrawal, total int64, err error)
	GetById(id uuid.UUID) (respondentWithdrawal *model.RespondentWithdrawal, err error)
	GetByRequest(request model.Request) (respondentWithdrawal *model.RespondentWithdrawal, err error)
	FindByRespondent(respondent model.Respondent) (withdrawals *[]model.RespondentWithdrawal, err error)
	Save(respondentWithdrawal *model.RespondentWithdrawal) (err error)
	Delete(respondentWithdrawal *model.RespondentWithdrawal) (error error)
	DeleteById(id uuid.UUID) (respondentWithdrawal *model.RespondentWithdrawal, err error)
}

type respondentWithdrawalRepository struct {
	db *gorm.DB
}

func (repo *respondentWithdrawalRepository) FindAll(page int, limit int) (respondentWithdrawals []model.RespondentWithdrawal, total int64, err error) {
	offset := (page - 1) * limit

	err = repo.db.
		Model(&model.RespondentWithdrawal{}).
		Count(&total).
		Limit(limit).
		Offset(offset).
		Order("created_at desc").
		Find(&respondentWithdrawals).Error
	if err != nil {
		return
	}
	return
}

func (repo *respondentWithdrawalRepository) FindByRespondent(respondent model.Respondent) (withdrawals *[]model.RespondentWithdrawal, err error) {
	walletRepo := GetRespondentWalletRepository()
	wallet, err := walletRepo.GetByRespondent(respondent)
	if err != nil {
		return nil, err
	}
	if err = model.DB.Where("wallet_id = ?", wallet.ID).Find(&withdrawals).Error; err != nil {
		return nil, err
	}
	return withdrawals, nil
}

func (repo *respondentWithdrawalRepository) GetById(id uuid.UUID) (respondentWithdrawal *model.RespondentWithdrawal, err error) {
	if err = model.DB. /*.Preload("Place")*/ Where("id = ?", id).First(&respondentWithdrawal).Error; err != nil {
		return nil, err
	}
	return respondentWithdrawal, nil
}

func (repo *respondentWithdrawalRepository) GetByRequest(request model.Request) (respondentWithdrawal *model.RespondentWithdrawal, err error) {
	if err = model.DB.Where("request_id = ?", request.ID).First(&respondentWithdrawal).Error; err != nil {
		return nil, err
	}
	return respondentWithdrawal, nil
}

func (repo *respondentWithdrawalRepository) Save(respondentWithdrawal *model.RespondentWithdrawal) (err error) {
	if (uuid.UUID{} == respondentWithdrawal.ID) {
		//NEW - No ID yet
		return repo.db.Create(&respondentWithdrawal).Error
	}
	return repo.db.Updates(&respondentWithdrawal).Error
}

func (repo *respondentWithdrawalRepository) Delete(respondentWithdrawal *model.RespondentWithdrawal) (err error) {
	return repo.db.Delete(&respondentWithdrawal).Error
}

func (repo *respondentWithdrawalRepository) DeleteById(id uuid.UUID) (respondentWithdrawal *model.RespondentWithdrawal, err error) {
	respondentWithdrawal, err = repo.GetById(id)
	if err != nil {
		return nil, err
	}
	err = repo.db.Delete(respondentWithdrawal).Error
	return respondentWithdrawal, err
}
