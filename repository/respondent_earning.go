package repository

import (
	"sync"

	"github.com/google/uuid"
	"github.com/nkamuo/rasta-server/model"
	"gorm.io/gorm"
)

var respondentEarningRepo RespondentEarningRepository
var respondentEarningRepoMutext *sync.Mutex = &sync.Mutex{}

func GetRespondentEarningRepository() RespondentEarningRepository {
	respondentEarningRepoMutext.Lock()
	if respondentEarningRepo == nil {
		respondentEarningRepo = &respondentEarningRepository{db: model.DB}
	}
	respondentEarningRepoMutext.Unlock()
	return respondentEarningRepo
}

type RespondentEarningRepository interface {
	FindAll(page int, limit int) (respondentEarnings []model.RespondentEarning, total int64, err error)
	GetById(id uuid.UUID) (respondentEarning *model.RespondentEarning, err error)
	GetByRequest(request model.Request) (respondentEarning *model.RespondentEarning, err error)
	FindByRespondent(respondent model.Respondent) (withdrawals *[]model.RespondentEarning, err error)
	Save(respondentEarning *model.RespondentEarning) (err error)
	Delete(respondentEarning *model.RespondentEarning) (error error)
	DeleteById(id uuid.UUID) (respondentEarning *model.RespondentEarning, err error)
}

type respondentEarningRepository struct {
	db *gorm.DB
}

func (repo *respondentEarningRepository) FindAll(page int, limit int) (respondentEarnings []model.RespondentEarning, total int64, err error) {
	offset := (page - 1) * limit

	err = repo.db.
		Model(&model.RespondentEarning{}).
		Count(&total).
		Limit(limit).
		Offset(offset).
		Order("created_at desc").
		Find(&respondentEarnings).Error
	if err != nil {
		return
	}
	return
}

func (repo *respondentEarningRepository) FindByRespondent(respondent model.Respondent) (withdrawals *[]model.RespondentEarning, err error) {
	if err = model.DB.Where("respondent_id = ?", respondent.ID).Find(&withdrawals).Error; err != nil {
		return nil, err
	}
	return withdrawals, nil
}

func (repo *respondentEarningRepository) GetById(id uuid.UUID) (respondentEarning *model.RespondentEarning, err error) {
	if err = model.DB. /*.Preload("Place")*/ Where("id = ?", id).First(&respondentEarning).Error; err != nil {
		return nil, err
	}
	return respondentEarning, nil
}

func (repo *respondentEarningRepository) GetByRequest(request model.Request) (respondentEarning *model.RespondentEarning, err error) {
	if err = model.DB.Where("request_id = ?", request.ID).First(&respondentEarning).Error; err != nil {
		return nil, err
	}
	return respondentEarning, nil
}

func (repo *respondentEarningRepository) Save(respondentEarning *model.RespondentEarning) (err error) {
	if (uuid.UUID{} == respondentEarning.ID) {
		//NEW - No ID yet
		return repo.db.Create(&respondentEarning).Error
	}
	return repo.db.Updates(&respondentEarning).Error
}

func (repo *respondentEarningRepository) Delete(respondentEarning *model.RespondentEarning) (err error) {
	return repo.db.Delete(&respondentEarning).Error
}

func (repo *respondentEarningRepository) DeleteById(id uuid.UUID) (respondentEarning *model.RespondentEarning, err error) {
	respondentEarning, err = repo.GetById(id)
	if err != nil {
		return nil, err
	}
	err = repo.db.Delete(respondentEarning).Error
	return respondentEarning, err
}
