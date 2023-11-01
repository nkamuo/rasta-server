package repository

import (
	"sync"

	"github.com/google/uuid"
	"github.com/nkamuo/rasta-server/model"
	"gorm.io/gorm"
)

var respondentChargeRepo RespondentOrderChargeRepository
var respondentChargeRepoMutext *sync.Mutex = &sync.Mutex{}

func GetRespondentOrderChargeRepository() RespondentOrderChargeRepository {
	respondentChargeRepoMutext.Lock()
	if respondentChargeRepo == nil {
		respondentChargeRepo = &respondentChargeRepository{db: model.DB}
	}
	respondentChargeRepoMutext.Unlock()
	return respondentChargeRepo
}

type RespondentOrderChargeRepository interface {
	FindAll(page int, limit int) (respondentCharges []model.RespondentOrderCharge, total int64, err error)
	GetById(id uuid.UUID) (respondentCharge *model.RespondentOrderCharge, err error)
	GetByRequest(request model.Request) (respondentCharge *model.RespondentOrderCharge, err error)
	FindByRespondent(respondent model.Respondent) (withdrawals *[]model.RespondentOrderCharge, err error)
	Save(respondentCharge *model.RespondentOrderCharge) (err error)
	Delete(respondentCharge *model.RespondentOrderCharge) (error error)
	DeleteById(id uuid.UUID) (respondentCharge *model.RespondentOrderCharge, err error)
}

type respondentChargeRepository struct {
	db *gorm.DB
}

func (repo *respondentChargeRepository) FindAll(page int, limit int) (respondentCharges []model.RespondentOrderCharge, total int64, err error) {
	offset := (page - 1) * limit

	err = repo.db.
		Model(&model.RespondentOrderCharge{}).
		Count(&total).
		Limit(limit).
		Offset(offset).
		Order("created_at desc").
		Find(&respondentCharges).Error
	if err != nil {
		return
	}
	return
}

func (repo *respondentChargeRepository) FindByRespondent(respondent model.Respondent) (withdrawals *[]model.RespondentOrderCharge, err error) {
	if err = model.DB.Where("respondent_id = ?", respondent.ID).Find(&withdrawals).Error; err != nil {
		return nil, err
	}
	return withdrawals, nil
}

func (repo *respondentChargeRepository) GetById(id uuid.UUID) (respondentCharge *model.RespondentOrderCharge, err error) {
	if err = model.DB. /*.Preload("Place")*/ Where("id = ?", id).First(&respondentCharge).Error; err != nil {
		return nil, err
	}
	return respondentCharge, nil
}

func (repo *respondentChargeRepository) GetByRequest(request model.Request) (respondentCharge *model.RespondentOrderCharge, err error) {
	if err = model.DB.Where("request_id = ?", request.ID).First(&respondentCharge).Error; err != nil {
		return nil, err
	}
	return respondentCharge, nil
}

func (repo *respondentChargeRepository) Save(respondentCharge *model.RespondentOrderCharge) (err error) {
	if (uuid.UUID{} == respondentCharge.ID) {
		//NEW - No ID yet
		return repo.db.Create(&respondentCharge).Error
	}
	return repo.db.Updates(&respondentCharge).Error
}

func (repo *respondentChargeRepository) Delete(respondentCharge *model.RespondentOrderCharge) (err error) {
	return repo.db.Delete(&respondentCharge).Error
}

func (repo *respondentChargeRepository) DeleteById(id uuid.UUID) (respondentCharge *model.RespondentOrderCharge, err error) {
	respondentCharge, err = repo.GetById(id)
	if err != nil {
		return nil, err
	}
	err = repo.db.Delete(respondentCharge).Error
	return respondentCharge, err
}
