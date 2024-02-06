package repository

import (
	"sync"

	"github.com/google/uuid"
	"github.com/nkamuo/rasta-server/model"
	"gorm.io/gorm"
)

var purchaseRepo RespondentAccessProductPurchaseRepository
var purchaseRepoMutext *sync.Mutex = &sync.Mutex{}

func GetRespondentAccessProductPurchaseRepository() RespondentAccessProductPurchaseRepository {
	purchaseRepoMutext.Lock()
	if purchaseRepo == nil {
		purchaseRepo = &purchaseRepository{db: model.DB}
	}
	purchaseRepoMutext.Unlock()
	return purchaseRepo
}

type RespondentAccessProductPurchaseRepository interface {
	FindAll(page int, limit int) (purchases []model.RespondentAccessProductPurchase, total int64, err error)
	GetById(id uuid.UUID, prefetch ...string) (purchase *model.RespondentAccessProductPurchase, err error)
	GetByStripeCheckoutID(id string, prefetch ...string) (purchase *model.RespondentAccessProductPurchase, err error)
	GetActiveByRespondent(respondent model.Respondent, prefetch ...string) (purchase *model.RespondentAccessProductPurchase, err error)
	Save(purchase *model.RespondentAccessProductPurchase) (err error)
	Delete(purchase *model.RespondentAccessProductPurchase) (error error)
	DeleteById(id uuid.UUID) (purchase *model.RespondentAccessProductPurchase, err error)
}

type purchaseRepository struct {
	db *gorm.DB
}

func (repo *purchaseRepository) FindAll(page int, limit int) (purchases []model.RespondentAccessProductPurchase, total int64, err error) {
	offset := (page - 1) * limit

	err = repo.db.
		Model(&model.RespondentAccessProductPurchase{}).
		Count(&total).
		Limit(limit).
		Offset(offset).
		Order("created_at desc").
		Find(&purchases).Error
	if err != nil {
		return
	}
	return
}

func (repo *purchaseRepository) GetById(id uuid.UUID, prefetch ...string) (purchase *model.RespondentAccessProductPurchase, err error) {
	query := model.DB

	for _, pFetch := range prefetch {
		query = query.Preload(pFetch)
	}

	if err = query.First(&purchase, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return purchase, nil
}

func (repo *purchaseRepository) GetByStripeCheckoutID(checkoutID string, prefetch ...string) (purchase *model.RespondentAccessProductPurchase, err error) {
	query := model.DB

	for _, pFetch := range prefetch {
		query = query.Preload(pFetch)
	}

	if err = query.First(&purchase, "stripe_checkout_session_id = ?", checkoutID).Error; err != nil {
		return nil, err
	}
	return purchase, nil
}

func (repo *purchaseRepository) GetActiveByRespondent(respondent model.Respondent, prefetch ...string) (purchase *model.RespondentAccessProductPurchase, err error) {
	query := model.DB //.Preload("Assignments.Assignment.Product").
	query = query.Where("respondent_id = ? AND active = ? AND ended_at IS NULL", respondent.ID, true)

	for _, pFetch := range prefetch {
		query = query.Preload(pFetch)
	}

	if err = query.First(&purchase).Error; err != nil {
		return nil, err
	}
	return purchase, nil
}

func (repo *purchaseRepository) Save(purchase *model.RespondentAccessProductPurchase) (err error) {
	if (uuid.UUID{} == purchase.ID) {
		//NEW - No ID yet
		return repo.db.Create(&purchase).Error
	}
	return repo.db.Updates(&purchase).Error
}

func (repo *purchaseRepository) Delete(purchase *model.RespondentAccessProductPurchase) (err error) {
	return repo.db.Delete(&purchase).Error
}

func (repo *purchaseRepository) DeleteById(id uuid.UUID) (purchase *model.RespondentAccessProductPurchase, err error) {
	purchase, err = repo.GetById(id)
	if err != nil {
		return nil, err
	}
	err = repo.db.Delete(&purchase).Error
	return purchase, err
}
