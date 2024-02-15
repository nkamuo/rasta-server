package repository

import (
	"sync"

	"github.com/google/uuid"
	"github.com/nkamuo/rasta-server/model"
	"gorm.io/gorm"
)

var priceRepo RespondentAccessProductPriceRepository
var priceRepoMutext *sync.Mutex = &sync.Mutex{}

func GetRespondentAccessProductPriceRepository() RespondentAccessProductPriceRepository {
	priceRepoMutext.Lock()
	if priceRepo == nil {
		priceRepo = &priceRepository{db: model.DB}
	}
	priceRepoMutext.Unlock()
	return priceRepo
}

type RespondentAccessProductPriceRepository interface {
	FindAll(page int, limit int) (prices []model.RespondentAccessProductPrice, total int64, err error)
	GetById(id uuid.UUID, prefetch ...string) (price *model.RespondentAccessProductPrice, err error)
	GetActiveByRespondent(respondent model.Respondent, prefetch ...string) (price *model.RespondentAccessProductPrice, err error)
	Save(price *model.RespondentAccessProductPrice) (err error)
	Delete(price *model.RespondentAccessProductPrice) (error error)
	DeleteById(id uuid.UUID) (price *model.RespondentAccessProductPrice, err error)
}

type priceRepository struct {
	db *gorm.DB
}

func (repo *priceRepository) FindAll(page int, limit int) (prices []model.RespondentAccessProductPrice, total int64, err error) {
	offset := (page - 1) * limit

	err = repo.db.
		Model(&model.RespondentAccessProductPrice{}).
		Count(&total).
		Limit(limit).
		Offset(offset).
		Order("created_at desc").
		Find(&prices).Error
	if err != nil {
		return
	}
	return
}

func (repo *priceRepository) GetById(id uuid.UUID, prefetch ...string) (price *model.RespondentAccessProductPrice, err error) {
	query := model.DB

	for _, pFetch := range prefetch {
		query = query.Preload(pFetch)
	}

	if err = query.First(&price, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return price, nil
}

func (repo *priceRepository) GetActiveByRespondent(respondent model.Respondent, prefetch ...string) (price *model.RespondentAccessProductPrice, err error) {
	query := model.DB //.Preload("Assignments.Assignment.Product").
	query = query.Where("respondent_id = ?", respondent.ID)

	for _, pFetch := range prefetch {
		query = query.Preload(pFetch)
	}

	if err = query.First(&price).Error; err != nil {
		return nil, err
	}
	return price, nil
}

func (repo *priceRepository) Save(price *model.RespondentAccessProductPrice) (err error) {
	if (uuid.UUID{} == price.ID) {
		//NEW - No ID yet
		return repo.db.Create(&price).Error
	}
	return repo.db.Updates(&price).Error
}

func (repo *priceRepository) Delete(price *model.RespondentAccessProductPrice) (err error) {
	return repo.db.Delete(&price).Error
}

func (repo *priceRepository) DeleteById(id uuid.UUID) (price *model.RespondentAccessProductPrice, err error) {
	price, err = repo.GetById(id)
	if err != nil {
		return nil, err
	}
	err = repo.db.Delete(&price).Error
	return price, err
}
