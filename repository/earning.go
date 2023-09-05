package repository

import (
	"sync"

	"github.com/google/uuid"
	"github.com/nkamuo/rasta-server/model"
	"gorm.io/gorm"
)

var generalEarningRepo GeneralEarningRepository
var generalEarningRepoMutext *sync.Mutex = &sync.Mutex{}

func GetGeneralEarningRepository() GeneralEarningRepository {
	generalEarningRepoMutext.Lock()
	if generalEarningRepo == nil {
		generalEarningRepo = &generalEarningRepository{db: model.DB}
	}
	generalEarningRepoMutext.Unlock()
	return generalEarningRepo
}

type GeneralEarningRepository interface {
	FindAll(page int, limit int) (generalEarnings []model.OrderEarning, total int64, err error)
	GetById(id uuid.UUID) (generalEarning *model.OrderEarning, err error)
	Save(generalEarning *model.OrderEarning) (err error)
	Delete(generalEarning *model.OrderEarning) (error error)
	DeleteById(id uuid.UUID) (generalEarning *model.OrderEarning, err error)
}

type generalEarningRepository struct {
	db *gorm.DB
}

func (repo *generalEarningRepository) FindAll(page int, limit int) (generalEarnings []model.OrderEarning, total int64, err error) {
	offset := (page - 1) * limit

	err = repo.db.
		Model(&model.OrderEarning{}).
		Count(&total).
		Limit(limit).
		Offset(offset).
		Order("created_at desc").
		Find(&generalEarnings).Error
	if err != nil {
		return
	}
	return
}

func (repo *generalEarningRepository) GetById(id uuid.UUID) (generalEarning *model.OrderEarning, err error) {
	if err = model.DB. /*.Preload("Place")*/ Where("id = ?", id).First(&generalEarning).Error; err != nil {
		return nil, err
	}
	return generalEarning, nil
}

func (repo *generalEarningRepository) Save(generalEarning *model.OrderEarning) (err error) {
	if (uuid.UUID{} == generalEarning.ID) {
		//NEW - No ID yet
		return repo.db.Create(&generalEarning).Error
	}
	return repo.db.Updates(&generalEarning).Error
}

func (repo *generalEarningRepository) Delete(generalEarning *model.OrderEarning) (err error) {
	return repo.db.Delete(&generalEarning).Error
}

func (repo *generalEarningRepository) DeleteById(id uuid.UUID) (generalEarning *model.OrderEarning, err error) {
	generalEarning, err = repo.GetById(id)
	if err != nil {
		return nil, err
	}
	err = repo.db.Delete(generalEarning).Error
	return generalEarning, err
}
