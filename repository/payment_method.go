package repository

import (
	"sync"

	"github.com/google/uuid"
	"github.com/nkamuo/rasta-server/model"
	"gorm.io/gorm"
)

var paymentMethodRepo PaymentMethodRepository
var paymentMethodRepoMutext *sync.Mutex = &sync.Mutex{}

func GetPaymentMethodRepository() PaymentMethodRepository {
	paymentMethodRepoMutext.Lock()
	if paymentMethodRepo == nil {
		paymentMethodRepo = &paymentMethodRepository{db: model.DB}
	}
	paymentMethodRepoMutext.Unlock()
	return paymentMethodRepo
}

type PaymentMethodRepository interface {
	FindAll(page int, limit int) (paymentMethods []model.PaymentMethod, total int64, err error)
	GetById(id uuid.UUID) (paymentMethod *model.PaymentMethod, err error)
	Save(paymentMethod *model.PaymentMethod) (err error)
	Delete(paymentMethod *model.PaymentMethod) (error error)
	DeleteById(id uuid.UUID) (paymentMethod *model.PaymentMethod, err error)
}

type paymentMethodRepository struct {
	db *gorm.DB
}

func (repo *paymentMethodRepository) FindAll(page int, limit int) (paymentMethods []model.PaymentMethod, total int64, err error) {
	offset := (page - 1) * limit

	err = repo.db.
		Model(&model.PaymentMethod{}).
		Count(&total).
		Limit(limit).
		Offset(offset).
		Order("created_at desc").
		Find(&paymentMethods).Error
	if err != nil {
		return
	}
	return
}

func (repo *paymentMethodRepository) GetById(id uuid.UUID) (paymentMethod *model.PaymentMethod, err error) {
	if err = model.DB. /*.Preload("Place")*/ Where("id = ?", id).First(&paymentMethod).Error; err != nil {
		return nil, err
	}
	return paymentMethod, nil
}

func (repo *paymentMethodRepository) Save(paymentMethod *model.PaymentMethod) (err error) {
	if (uuid.UUID{} == paymentMethod.ID) {
		//NEW - No ID yet
		return repo.db.Create(&paymentMethod).Error
	}
	return repo.db.Updates(&paymentMethod).Error
}

func (repo *paymentMethodRepository) Delete(paymentMethod *model.PaymentMethod) (err error) {
	return repo.db.Delete(&paymentMethod).Error
}

func (repo *paymentMethodRepository) DeleteById(id uuid.UUID) (paymentMethod *model.PaymentMethod, err error) {
	paymentMethod, err = repo.GetById(id)
	if err != nil {
		return nil, err
	}
	err = repo.db.Delete(paymentMethod).Error
	return paymentMethod, err
}
