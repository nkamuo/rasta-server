package repository

import (
	"sync"

	"github.com/google/uuid"
	"github.com/nkamuo/rasta-server/model"
	"gorm.io/gorm"
)

var productRepo ProductRepository
var productRepoMutext *sync.Mutex = &sync.Mutex{}

func GetProductRepository() ProductRepository {
	productRepoMutext.Lock()
	if productRepo == nil {
		productRepo = &productRepository{db: model.DB}
	}
	productRepoMutext.Unlock()
	return productRepo
}

type ProductRepository interface {
	FindAll(page int, limit int) (products []model.Product, total int64, err error)
	GetById(id uuid.UUID, preload ...string) (product *model.Product, err error)
	GetByPlaceIdAndCategory(id uuid.UUID, category model.ProductCategory) (product *model.Product, err error)
	Save(product *model.Product) (err error)
	Delete(product *model.Product) (error error)
	DeleteById(id uuid.UUID) (product *model.Product, err error)
}

type productRepository struct {
	db *gorm.DB
}

func (repo *productRepository) FindAll(page int, limit int) (products []model.Product, total int64, err error) {
	offset := (page - 1) * limit

	err = repo.db.
		Model(&model.Product{}).
		Count(&total).
		Limit(limit).
		Offset(offset).
		Order("created_at desc").
		Find(&products).Error
	if err != nil {
		return
	}
	return
}

func (repo *productRepository) GetById(id uuid.UUID, preload ...string) (product *model.Product, err error) {
	query := model.DB
	for _, pLoad := range preload {
		query = query.Preload(pLoad)
	}

	if err = query.Where("id = ?", id).First(&product).Error; err != nil {
		return nil, err
	}
	return product, nil
}

func (repo *productRepository) GetByPlaceIdAndCategory(placeID uuid.UUID, category model.ProductCategory) (product *model.Product, err error) {
	if err = model.DB.Where("place_id = ? AND category = ?", placeID, category).First(&product).Error; err != nil {
		return nil, err
	}
	return product, nil
}

func (repo *productRepository) Save(product *model.Product) (err error) {
	if (uuid.UUID{} == product.ID) {
		//NEW - No ID yet
		return repo.db.Create(&product).Error
	}
	return repo.db.Updates(&product).Error
}

func (repo *productRepository) Delete(product *model.Product) (err error) {
	return repo.db.Delete(&product).Error
}

func (repo *productRepository) DeleteById(id uuid.UUID) (product *model.Product, err error) {
	product, err = repo.GetById(id)
	if err != nil {
		return nil, err
	}
	err = repo.db.Delete(product).Error
	return product, err
}
