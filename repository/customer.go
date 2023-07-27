package repository

import (
	"sync"

	"github.com/google/uuid"
	"github.com/nkamuo/rasta-server/model"
	"gorm.io/gorm"
)

var customerRepo *CustomerRepository
var customerRepoMutext *sync.Mutex

func GetCustomerRepository(db *gorm.DB) *CustomerRepository {
	customerRepoMutext.Lock()
	if customerRepo == nil {
		*customerRepo = &customerRepository{db: db}
	}
	customerRepoMutext.Unlock()
	return customerRepo
}

type CustomerRepository interface {
	FindAll(page int, limit int) (customers []model.Customer, total int64, err error)
	GetById(id uuid.UUID) (customer *model.Customer, err error)
	Save(customer *model.Customer) (err error)
	Delete(customer *model.Customer) (error error)
	DeleteById(id uuid.UUID) (customer *model.Customer, err error)
}

type customerRepository struct {
	db *gorm.DB
}

func (repo *customerRepository) FindAll(page int, limit int) (customers []model.Customer, total int64, err error) {
	offset := (page - 1) * limit

	err = repo.db.
		Model(&model.Customer{}).
		Count(&total).
		Limit(limit).
		Offset(offset).
		Order("created_at desc").
		Find(&customers).Error
	if err != nil {
		return
	}
	return
}

func (repo *customerRepository) GetById(id uuid.UUID) (customer *model.Customer, err error) {
	if err = model.DB.Where("id = ?", id).First(&customer).Error; err != nil {
		return nil, err
	}
	return customer, nil
}

func (repo *customerRepository) Save(customer *model.Customer) (err error) {
	if (uuid.UUID{} == customer.ID) {
		//NEW - No ID yet
		repo.db.Create(&customer)
		return nil
	}
	repo.db.Updates(&customer)
	return nil
}

func (repo *customerRepository) Delete(customer *model.Customer) (err error) {
	repo.db.Delete(&customer)
	return nil
}

func (repo *customerRepository) DeleteById(id uuid.UUID) (customer *model.Customer, err error) {
	customer, err = repo.GetById(id)
	if err != nil {
		return nil, err
	}
	return customer, nil
}
