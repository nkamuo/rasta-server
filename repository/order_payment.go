package repository

import (
	"sync"

	"github.com/google/uuid"
	"github.com/nkamuo/rasta-server/model"
	"gorm.io/gorm"
)

var paymentRepo OrderPaymentRepository
var paymentRepoMutext *sync.Mutex = &sync.Mutex{}

func GetOrderPaymentRepository() OrderPaymentRepository {
	paymentRepoMutext.Lock()
	if paymentRepo == nil {
		paymentRepo = &paymentRepository{db: model.DB}
	}
	paymentRepoMutext.Unlock()
	return paymentRepo
}

type OrderPaymentRepository interface {
	FindAll(page int, limit int) (payments []model.OrderPayment, total int64, err error)
	GetById(id uuid.UUID) (payment *model.OrderPayment, err error)
	Save(payment *model.OrderPayment) (err error)
	Delete(payment *model.OrderPayment) (error error)
	DeleteById(id uuid.UUID) (payment *model.OrderPayment, err error)
}

type paymentRepository struct {
	db *gorm.DB
}

func (repo *paymentRepository) FindAll(page int, limit int) (payments []model.OrderPayment, total int64, err error) {
	offset := (page - 1) * limit

	err = repo.db.
		Model(&model.OrderPayment{}).
		Count(&total).
		Limit(limit).
		Offset(offset).
		Order("created_at desc").
		Find(&payments).Error
	if err != nil {
		return
	}
	return
}

func (repo *paymentRepository) GetById(id uuid.UUID) (payment *model.OrderPayment, err error) {
	if err = model.DB. /*.Preload("Place")*/ Where("id = ?", id).First(&payment).Error; err != nil {
		return nil, err
	}
	return payment, nil
}

func (repo *paymentRepository) Save(payment *model.OrderPayment) (err error) {
	if (uuid.UUID{} == payment.ID) {
		//NEW - No ID yet
		return repo.db.Create(&payment).Error
	}
	return repo.db.Updates(&payment).Error
}

func (repo *paymentRepository) Delete(payment *model.OrderPayment) (err error) {
	return repo.db.Delete(&payment).Error
}

func (repo *paymentRepository) DeleteById(id uuid.UUID) (payment *model.OrderPayment, err error) {
	payment, err = repo.GetById(id)
	if err != nil {
		return nil, err
	}
	err = repo.db.Delete(payment).Error
	return payment, err
}
