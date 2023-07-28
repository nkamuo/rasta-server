package service

import (
	"sync"

	"github.com/google/uuid"
	"github.com/nkamuo/rasta-server/model"
	"github.com/nkamuo/rasta-server/repository"
)

var productService ProductService
var productRepoMutext *sync.Mutex = &sync.Mutex{}

func GetProductService() ProductService {
	productRepoMutext.Lock()
	if productService == nil {
		productService = &productServiceImpl{repo: repository.GetProductRepository()}
	}
	productRepoMutext.Unlock()
	return productService
}

type ProductService interface {
	GetById(id uuid.UUID) (product *model.Product, err error)
	// GetByEmail(email string) (product *model.Product, err error)
	// GetByPhone(phone string) (product *model.Product, err error)
	Save(product *model.Product) (err error)
	Delete(product *model.Product) (error error)
}

type productServiceImpl struct {
	repo repository.ProductRepository
}

func (service *productServiceImpl) GetById(id uuid.UUID) (product *model.Product, err error) {
	return service.repo.GetById(id)
}

func (service *productServiceImpl) Save(product *model.Product) (err error) {
	return service.repo.Save(product)
}

func (service *productServiceImpl) Delete(product *model.Product) (err error) {
	err = service.repo.Delete(product)

	return err
}

func (service *productServiceImpl) DeleteById(id uuid.UUID) (product *model.Product, err error) {
	product, err = service.repo.DeleteById(id)
	return product, err
}
