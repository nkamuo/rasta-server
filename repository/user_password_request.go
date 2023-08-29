package repository

import (
	"sync"

	"github.com/google/uuid"
	"github.com/nkamuo/rasta-server/model"
	"gorm.io/gorm"
)

var userPasswordResetRequestRepo UserPasswordResetRequestRepository
var userPasswordResetRequestRepoMutext *sync.Mutex = &sync.Mutex{}

func GetUserPasswordResetRequestRepository() UserPasswordResetRequestRepository {
	userPasswordResetRequestRepoMutext.Lock()
	if userPasswordResetRequestRepo == nil {
		userPasswordResetRequestRepo = &userPasswordResetRequestRepository{db: model.DB}
	}
	userPasswordResetRequestRepoMutext.Unlock()
	return userPasswordResetRequestRepo
}

type UserPasswordResetRequestRepository interface {
	FindAll(page int, limit int) (userPasswordResetRequests []model.UserPasswordResetRequest, total int64, err error)
	GetById(id uuid.UUID) (userPasswordResetRequest *model.UserPasswordResetRequest, err error)
	Save(userPasswordResetRequest *model.UserPasswordResetRequest) (err error)
	Delete(userPasswordResetRequest *model.UserPasswordResetRequest) (error error)
	DeleteById(id uuid.UUID) (userPasswordResetRequest *model.UserPasswordResetRequest, err error)
}

type userPasswordResetRequestRepository struct {
	db *gorm.DB
}

func (repo *userPasswordResetRequestRepository) FindAll(page int, limit int) (userPasswordResetRequests []model.UserPasswordResetRequest, total int64, err error) {
	offset := (page - 1) * limit

	err = repo.db.
		Model(&model.UserPasswordResetRequest{}).
		Count(&total).
		Limit(limit).
		Offset(offset).
		Order("created_at desc").
		Find(&userPasswordResetRequests).Error
	if err != nil {
		return
	}
	return
}

func (repo *userPasswordResetRequestRepository) GetById(id uuid.UUID) (userPasswordResetRequest *model.UserPasswordResetRequest, err error) {
	if err = model.DB. /*.Preload("Place")*/ Where("id = ?", id).First(&userPasswordResetRequest).Error; err != nil {
		return nil, err
	}
	return userPasswordResetRequest, nil
}

func (repo *userPasswordResetRequestRepository) Save(userPasswordResetRequest *model.UserPasswordResetRequest) (err error) {
	if (uuid.UUID{} == userPasswordResetRequest.ID) {
		//NEW - No ID yet
		return repo.db.Create(&userPasswordResetRequest).Error
	}
	return repo.db.Updates(&userPasswordResetRequest).Error
}

func (repo *userPasswordResetRequestRepository) GetByToken(token string) (request *model.UserPasswordResetRequest, err error) {
	if err = model.DB.Where("token = ?", token).First(request).Error; err != nil {
		return nil, err
	}
	return request, nil
}

// func (repo *userPasswordResetRequestRepository) GetByUser(user *model.User) (request *model.UserPasswordResetRequest,err error){
// 	if err = model.DB.Where("user_id = ?", user.ID).First(request).Error; err != nil {
// 		return nil, err
// 	}
// 	return request, nil
// }

func (repo *userPasswordResetRequestRepository) Delete(userPasswordResetRequest *model.UserPasswordResetRequest) (err error) {
	return repo.db.Delete(&userPasswordResetRequest).Error
}

func (repo *userPasswordResetRequestRepository) DeleteById(id uuid.UUID) (userPasswordResetRequest *model.UserPasswordResetRequest, err error) {
	userPasswordResetRequest, err = repo.GetById(id)
	if err != nil {
		return nil, err
	}
	err = repo.db.Delete(userPasswordResetRequest).Error
	return userPasswordResetRequest, err
}
