package repository

import (
	"sync"

	"github.com/google/uuid"
	"github.com/nkamuo/rasta-server/model"
	"gorm.io/gorm"
)

var documentRepo ImageDocumentRepository
var documentRepoMutext *sync.Mutex = &sync.Mutex{}

func GetImageDocumentRepository() ImageDocumentRepository {
	documentRepoMutext.Lock()
	if documentRepo == nil {
		documentRepo = &documentRepository{db: model.DB}
	}
	documentRepoMutext.Unlock()
	return documentRepo
}

type ImageDocumentRepository interface {
	FindAll(page int, limit int) (documents []model.ImageDocument, total int64, err error)
	GetById(id uuid.UUID) (document *model.ImageDocument, err error)
	Save(document *model.ImageDocument) (err error)
	Delete(document *model.ImageDocument) (error error)
	DeleteById(id uuid.UUID) (document *model.ImageDocument, err error)
}

type documentRepository struct {
	db *gorm.DB
}

func (repo *documentRepository) FindAll(page int, limit int) (documents []model.ImageDocument, total int64, err error) {
	offset := (page - 1) * limit

	err = repo.db.
		Model(&model.ImageDocument{}).
		Count(&total).
		Limit(limit).
		Offset(offset).
		Order("created_at desc").
		Find(&documents).Error
	if err != nil {
		return
	}
	return
}

func (repo *documentRepository) GetById(id uuid.UUID) (document *model.ImageDocument, err error) {
	if err = model.DB. /*.Joins("OperatorUser")*/ First(&document, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return document, nil
}

// func (repo *documentRepository) GetByCode(code string) (document *model.ImageDocument, err error) {
// 	if err = model.DB. /*.Joins("OperatorUser")*/ First(&document, "code = ?", code).Error; err != nil {
// 		return nil, err
// 	}
// 	return document, nil
// }

func (repo *documentRepository) Save(document *model.ImageDocument) (err error) {
	if (uuid.UUID{} == document.ID) {
		//NEW - No ID yet
		return repo.db.Create(&document).Error
	}
	return repo.db.Updates(&document).Error
}

func (repo *documentRepository) Delete(document *model.ImageDocument) (err error) {
	return repo.db.Delete(&document).Error
}

func (repo *documentRepository) DeleteById(id uuid.UUID) (document *model.ImageDocument, err error) {
	document, err = repo.GetById(id)
	if err != nil {
		return nil, err
	}
	err = repo.db.Delete(&document).Error
	return document, err
}
