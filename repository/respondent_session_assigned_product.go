package repository

import (
	"sync"

	"github.com/google/uuid"
	"github.com/nkamuo/rasta-server/model"
	"gorm.io/gorm"
)

var assignedProductRepo RespondentSessionAssignedProductRepository
var assignedProductRepoMutext *sync.Mutex = &sync.Mutex{}

func GetRespondentSessionAssignedProductRepository() RespondentSessionAssignedProductRepository {
	assignedProductRepoMutext.Lock()
	if assignedProductRepo == nil {
		assignedProductRepo = &assignedProductRepository{db: model.DB}
	}
	assignedProductRepoMutext.Unlock()
	return assignedProductRepo
}

type RespondentSessionAssignedProductRepository interface {
	FindAll(page int, limit int) (assignedProducts []model.RespondentSessionAssignedProduct, total int64, err error)
	GetById(id uuid.UUID) (assignedProduct *model.RespondentSessionAssignedProduct, err error)
	Save(assignedProduct *model.RespondentSessionAssignedProduct) (err error)
	Delete(assignedProduct *model.RespondentSessionAssignedProduct) (error error)
	DeleteById(id uuid.UUID) (assignedProduct *model.RespondentSessionAssignedProduct, err error)
}

type assignedProductRepository struct {
	db *gorm.DB
}

func (repo *assignedProductRepository) FindAll(page int, limit int) (assignedProducts []model.RespondentSessionAssignedProduct, total int64, err error) {
	offset := (page - 1) * limit

	err = repo.db.
		Model(&model.RespondentSessionAssignedProduct{}).
		Count(&total).
		Limit(limit).
		Offset(offset).
		Order("created_at desc").
		Find(&assignedProducts).Error
	if err != nil {
		return
	}
	return
}

func (repo *assignedProductRepository) GetById(id uuid.UUID) (assignedProduct *model.RespondentSessionAssignedProduct, err error) {
	if err = model.DB. /*.Joins("OperatorUser")*/ First(&assignedProduct, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return assignedProduct, nil
}

func (repo *assignedProductRepository) Save(assignedProduct *model.RespondentSessionAssignedProduct) (err error) {
	if (uuid.UUID{} == assignedProduct.ID) {
		//NEW - No ID yet
		return repo.db.Create(&assignedProduct).Error
	}
	return repo.db.Updates(&assignedProduct).Error
}

func (repo *assignedProductRepository) Delete(assignedProduct *model.RespondentSessionAssignedProduct) (err error) {
	return repo.db.Delete(&assignedProduct).Error
}

func (repo *assignedProductRepository) DeleteById(id uuid.UUID) (assignedProduct *model.RespondentSessionAssignedProduct, err error) {
	assignedProduct, err = repo.GetById(id)
	if err != nil {
		return nil, err
	}
	err = repo.db.Delete(&assignedProduct).Error
	return assignedProduct, err
}
