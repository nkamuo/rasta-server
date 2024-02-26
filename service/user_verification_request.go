package service

import (
	"sync"

	"github.com/google/uuid"
	"github.com/nkamuo/rasta-server/model"
	"github.com/nkamuo/rasta-server/repository"
)

var userVerificationRequestService UserVerificationRequestService
var userVerificationRequestRepoMutext *sync.Mutex = &sync.Mutex{}

func GetUserVerificationRequestService() UserVerificationRequestService {
	userVerificationRequestRepoMutext.Lock()
	if userVerificationRequestService == nil {
		userVerificationRequestService = &userVerificationRequestServiceImpl{repo: repository.GetUserVerificationRequestRepository()}
	}
	userVerificationRequestRepoMutext.Unlock()
	return userVerificationRequestService
}

type UserVerificationRequestService interface {
	GetById(id uuid.UUID, preload ...string) (userVerificationRequest *model.UserVerificationRequest, err error)
	// GetByEmail(email string) (userVerificationRequest *model.UserVerificationRequest, err error)
	// GetByPhone(phone string) (userVerificationRequest *model.UserVerificationRequest, err error)
	GetByCodeAndEmail(code string, email string, preload ...string) (userVerificationRequest *model.UserVerificationRequest, err error)
	Save(userVerificationRequest *model.UserVerificationRequest) (err error)
	Delete(userVerificationRequest *model.UserVerificationRequest) (error error)
}

type userVerificationRequestServiceImpl struct {
	repo repository.UserVerificationRequestRepository
}

func (service *userVerificationRequestServiceImpl) GetById(id uuid.UUID, preload ...string) (userVerificationRequest *model.UserVerificationRequest, err error) {
	return service.repo.GetById(id, preload...)
}

func (service *userVerificationRequestServiceImpl) Save(userVerificationRequest *model.UserVerificationRequest) (err error) {
	return service.repo.Save(userVerificationRequest)
}

func (service *userVerificationRequestServiceImpl) Delete(userVerificationRequest *model.UserVerificationRequest) (err error) {
	err = service.repo.Delete(userVerificationRequest)
	return err
}

func (service *userVerificationRequestServiceImpl) DeleteById(id uuid.UUID) (userVerificationRequest *model.UserVerificationRequest, err error) {
	userVerificationRequest, err = service.repo.DeleteById(id)
	return userVerificationRequest, err
}

func (service *userVerificationRequestServiceImpl) GetByCodeAndEmail(code string, email string, preload ...string) (userVerificationRequest *model.UserVerificationRequest, err error) {
	return service.repo.GetByCodeAndEmail(code, email, preload...)
}
