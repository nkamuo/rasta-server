package repository

import (
	"sync"

	"github.com/google/uuid"
	"github.com/nkamuo/rasta-server/model"
	"gorm.io/gorm"
)

var userRepo UserRepository
var userRepoMutext *sync.Mutex = &sync.Mutex{}

func GetUserRepository() UserRepository {
	userRepoMutext.Lock()
	if userRepo == nil {
		userRepo = &userRepository{db: model.DB}
	}
	userRepoMutext.Unlock()
	return userRepo
}

type UserRepository interface {
	FindAll(page int, limit int) (users []model.User, total int64, err error)
	GetById(id uuid.UUID) (user *model.User, err error)
	GetByEmail(email string) (user *model.User, err error)
	GetByPhone(phone string) (user *model.User, err error)
	Save(user *model.User) (err error)
	Delete(user *model.User) (error error)
	DeleteById(id uuid.UUID) (user *model.User, err error)
}

type userRepository struct {
	db *gorm.DB
}

func (repo *userRepository) FindAll(page int, limit int) (users []model.User, total int64, err error) {
	offset := (page - 1) * limit

	err = repo.db.
		Model(&model.User{}).
		Count(&total).
		Limit(limit).
		Offset(offset).
		Order("created_at desc").
		Find(&users).Error
	if err != nil {
		return
	}
	return
}

func (repo *userRepository) GetById(id uuid.UUID) (user *model.User, err error) {
	if err = model.DB.Where("id = ?", id).First(&user).Error; err != nil {
		return nil, err
	}
	return user, nil
}

func (repo *userRepository) GetByEmail(email string) (user *model.User, err error) {
	if err = model.DB.Where("email = ?", email).First(&user).Error; err != nil {
		return nil, err
	}
	return user, nil
}

func (repo *userRepository) GetByPhone(phone string) (user *model.User, err error) {
	if err = model.DB.Where("phone = ?", phone).First(&user).Error; err != nil {
		return nil, err
	}
	return user, nil
}

func (repo *userRepository) Save(user *model.User) (err error) {
	if (uuid.UUID{} == user.ID) {
		//NEW - No ID yet
		repo.db.Create(&user)
		return nil
	}
	repo.db.Updates(&user)
	return nil
}

func (repo *userRepository) Delete(user *model.User) (err error) {
	repo.db.Delete(&user)
	return nil
}

func (repo *userRepository) DeleteById(id uuid.UUID) (user *model.User, err error) {
	user, err = repo.GetById(id)
	if err != nil {
		return nil, err
	}
	return user, nil
}
