package service

import (
	"sync"

	"github.com/google/uuid"
	"github.com/nkamuo/rasta-server/model"
	"github.com/nkamuo/rasta-server/repository"
	"golang.org/x/crypto/bcrypt"
)

var userService UserService
var userRepoMutext *sync.Mutex = &sync.Mutex{}

func GetUserService() UserService {
	userRepoMutext.Lock()
	if userService == nil {
		userService = &userServiceImpl{repo: repository.GetUserRepository()}
	}
	userRepoMutext.Unlock()
	return userService
}

type UserService interface {
	GetById(id uuid.UUID) (user *model.User, err error)
	// GetByEmail(email string) (user *model.User, err error)
	// GetByPhone(phone string) (user *model.User, err error)
	Save(user *model.User) (err error)
	Delete(user *model.User) (error error)
	HashUserPassword(user *model.User, Password string) (err error)
	ValidateEmail(user *model.User) (err error)
}

type userServiceImpl struct {
	repo repository.UserRepository
}

func (service *userServiceImpl) GetById(id uuid.UUID) (user *model.User, err error) {
	return service.repo.GetById(id)
}

func (service *userServiceImpl) Save(user *model.User) (err error) {
	return service.repo.Save(user)
}

func (service *userServiceImpl) ValidateEmail(user *model.User) (err error) {
	return nil
}

func (service *userServiceImpl) Delete(user *model.User) (err error) {
	err = service.repo.Delete(user)

	return err
}

func (service *userServiceImpl) DeleteById(id uuid.UUID) (user *model.User, err error) {
	user, err = service.repo.DeleteById(id)
	return user, err
}

func (service *userServiceImpl) HashUserPassword(user *model.User, Password string) (err error) {
	if err != nil {
		return err
	}
	var password *model.UserPassword

	password, err = service.repo.GetPassword(user)
	if err != nil {
		if err.Error() == "record not found" {
			password = &model.UserPassword{}
		} else {
			return err
		}
	}
	//turn password into hash
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(Password), bcrypt.DefaultCost)
	password.HashedPassword = string(hashedPassword)
	user.Password = password

	return nil
}
