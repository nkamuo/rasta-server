package repository

import (
	"sync"

	"github.com/google/uuid"
	"github.com/nkamuo/rasta-server/model"
	"gorm.io/gorm"
)

var assignmentRepo ProductRespondentAssignmentRepository
var assignmentRepoMutext *sync.Mutex = &sync.Mutex{}

func GetProductRespondentAssignmentRepository() ProductRespondentAssignmentRepository {
	assignmentRepoMutext.Lock()
	if assignmentRepo == nil {
		assignmentRepo = &assignmentRepository{db: model.DB}
	}
	assignmentRepoMutext.Unlock()
	return assignmentRepo
}

type ProductRespondentAssignmentRepository interface {
	FindAll(page int, limit int) (assignments []model.ProductRespondentAssignment, total int64, err error)
	GetById(id uuid.UUID) (assignment *model.ProductRespondentAssignment, err error)
	Save(assignment *model.ProductRespondentAssignment) (err error)
	Delete(assignment *model.ProductRespondentAssignment) (error error)
	DeleteById(id uuid.UUID) (assignment *model.ProductRespondentAssignment, err error)
}

type assignmentRepository struct {
	db *gorm.DB
}

func (repo *assignmentRepository) FindAll(page int, limit int) (assignments []model.ProductRespondentAssignment, total int64, err error) {
	offset := (page - 1) * limit

	err = repo.db.
		Model(&model.ProductRespondentAssignment{}).
		Count(&total).
		Limit(limit).
		Offset(offset).
		Order("created_at desc").
		Find(&assignments).Error
	if err != nil {
		return
	}
	return
}

func (repo *assignmentRepository) GetById(id uuid.UUID) (assignment *model.ProductRespondentAssignment, err error) {
	if err = model.DB. /*.Joins("OperatorUser")*/ First(&assignment, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return assignment, nil
}

func (repo *assignmentRepository) Save(assignment *model.ProductRespondentAssignment) (err error) {
	if (uuid.UUID{} == assignment.ID) {
		//NEW - No ID yet
		return repo.db.Create(&assignment).Error
	}
	return repo.db.Updates(&assignment).Error
}

func (repo *assignmentRepository) Delete(assignment *model.ProductRespondentAssignment) (err error) {
	return repo.db.Delete(&assignment).Error
}

func (repo *assignmentRepository) DeleteById(id uuid.UUID) (assignment *model.ProductRespondentAssignment, err error) {
	assignment, err = repo.GetById(id)
	if err != nil {
		return nil, err
	}
	err = repo.db.Delete(&assignment).Error
	return assignment, err
}
