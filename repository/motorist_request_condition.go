package repository

import (
	"sync"

	"github.com/google/uuid"
	"github.com/nkamuo/rasta-server/model"
	"gorm.io/gorm"
)

var motoristRequestSituationRepo MotoristRequestSituationRepository
var motoristRequestSituationRepoMutext *sync.Mutex = &sync.Mutex{}

func GetMotoristRequestSituationRepository() MotoristRequestSituationRepository {
	motoristRequestSituationRepoMutext.Lock()
	if motoristRequestSituationRepo == nil {
		motoristRequestSituationRepo = &motoristRequestSituationRepository{db: model.DB}
	}
	motoristRequestSituationRepoMutext.Unlock()
	return motoristRequestSituationRepo
}

type MotoristRequestSituationRepository interface {
	FindAll(page int, limit int) (motoristRequestSituations []model.MotoristRequestSituation, total int64, err error)
	FindAllDefault() (motoristRequestSituations []model.MotoristRequestSituation, total int64, err error)
	GetById(id uuid.UUID) (motoristRequestSituation *model.MotoristRequestSituation, err error)
	Save(motoristRequestSituation *model.MotoristRequestSituation) (err error)
	Delete(motoristRequestSituation *model.MotoristRequestSituation) (error error)
	DeleteById(id uuid.UUID) (motoristRequestSituation *model.MotoristRequestSituation, err error)
}

type motoristRequestSituationRepository struct {
	db *gorm.DB
}

func (repo *motoristRequestSituationRepository) FindAll(page int, limit int) (motoristRequestSituations []model.MotoristRequestSituation, total int64, err error) {
	offset := (page - 1) * limit

	err = repo.db.
		Model(&model.MotoristRequestSituation{}).
		Count(&total).
		Limit(limit).
		Offset(offset).
		Order("created_at desc").
		Find(&motoristRequestSituations).Error
	if err != nil {
		return
	}
	return
}

func (repo *motoristRequestSituationRepository) FindAllDefault() (motoristRequestSituations []model.MotoristRequestSituation, total int64, err error) {
	Situations := model.GetDefaultMotoristRequestSituations()
	return Situations, int64(len(Situations)), nil
}

func (repo *motoristRequestSituationRepository) GetById(id uuid.UUID) (motoristRequestSituation *model.MotoristRequestSituation, err error) {
	if err = model.DB. /*.Preload("Place")*/ Where("id = ?", id).First(&motoristRequestSituation).Error; err != nil {
		return nil, err
	}
	return motoristRequestSituation, nil
}

func (repo *motoristRequestSituationRepository) Save(motoristRequestSituation *model.MotoristRequestSituation) (err error) {
	if (uuid.UUID{} == motoristRequestSituation.ID) {
		return repo.db.Create(&motoristRequestSituation).Error
	}
	return repo.db.Updates(&motoristRequestSituation).Error
}

func (repo *motoristRequestSituationRepository) Delete(motoristRequestSituation *model.MotoristRequestSituation) (err error) {
	return repo.db.Delete(&motoristRequestSituation).Error
}

func (repo *motoristRequestSituationRepository) DeleteById(id uuid.UUID) (motoristRequestSituation *model.MotoristRequestSituation, err error) {
	motoristRequestSituation, err = repo.GetById(id)
	if err != nil {
		return nil, err
	}
	err = repo.db.Delete(motoristRequestSituation).Error
	return motoristRequestSituation, err
}
