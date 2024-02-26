package repository

import (
	"sync"

	"github.com/google/uuid"
	"github.com/nkamuo/rasta-server/model"
	"gorm.io/gorm"
)

var userVerificationRequestRepo UserVerificationRequestRepository
var userVerificationRequestRepoMutext *sync.Mutex = &sync.Mutex{}

func GetUserVerificationRequestRepository() UserVerificationRequestRepository {
	userVerificationRequestRepoMutext.Lock()
	if userVerificationRequestRepo == nil {
		userVerificationRequestRepo = &userVerificationRequestRepository{db: model.DB}
	}
	userVerificationRequestRepoMutext.Unlock()
	return userVerificationRequestRepo
}

type UserVerificationRequestRepository interface {
	FindAll(page int, limit int) (userVerificationRequests []model.UserVerificationRequest, total int64, err error)
	GetById(id uuid.UUID, preload ...string) (userVerificationRequest *model.UserVerificationRequest, err error)
	GetByCodeAndEmail(code string, email string, preload ...string) (userVerificationRequest *model.UserVerificationRequest, err error)
	Save(userVerificationRequest *model.UserVerificationRequest) (err error)
	Delete(userVerificationRequest *model.UserVerificationRequest) (error error)
	DeleteById(id uuid.UUID) (userVerificationRequest *model.UserVerificationRequest, err error)
}

type userVerificationRequestRepository struct {
	db *gorm.DB
}

func (repo *userVerificationRequestRepository) FindAll(page int, limit int) (userVerificationRequests []model.UserVerificationRequest, total int64, err error) {
	offset := (page - 1) * limit

	err = repo.db.
		Model(&model.UserVerificationRequest{}).
		Count(&total).
		Limit(limit).
		Offset(offset).
		Order("created_at desc").
		Find(&userVerificationRequests).Error
	if err != nil {
		return
	}
	return
}

func (repo *userVerificationRequestRepository) GetById(id uuid.UUID, preload ...string) (userVerificationRequest *model.UserVerificationRequest, err error) {

	query := model.DB
	for _, pLoad := range preload {
		query = query.Preload(pLoad)
	}

	if err = query.Where("id = ?", id).First(&userVerificationRequest).Error; err != nil {
		return nil, err
	}
	return userVerificationRequest, nil
}

func (repo *userVerificationRequestRepository) Save(userVerificationRequest *model.UserVerificationRequest) (err error) {
	if (uuid.UUID{} == userVerificationRequest.ID) {
		//NEW - No ID yet
		return repo.db.Create(&userVerificationRequest).Error
	}
	return repo.db.Updates(&userVerificationRequest).Error
}

func (repo *userVerificationRequestRepository) Delete(userVerificationRequest *model.UserVerificationRequest) (err error) {
	return repo.db.Delete(&userVerificationRequest).Error
}

func (repo *userVerificationRequestRepository) DeleteById(id uuid.UUID) (userVerificationRequest *model.UserVerificationRequest, err error) {
	userVerificationRequest, err = repo.GetById(id)
	if err != nil {
		return nil, err
	}
	err = repo.db.Delete(userVerificationRequest).Error
	return userVerificationRequest, err
}

func (repo *userVerificationRequestRepository) GetByCodeAndEmail(code string, email string, preload ...string) (userVerificationRequest *model.UserVerificationRequest, err error) {

	query := model.DB
	for _, pLoad := range preload {
		query = query.Preload(pLoad)
	}

	if err = query.Where("code = ? AND email = ?", code, email).First(&userVerificationRequest).Error; err != nil {
		return nil, err
	}
	return userVerificationRequest, nil
}
