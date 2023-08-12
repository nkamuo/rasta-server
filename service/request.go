package service

import (
	"sync"

	"github.com/google/uuid"
	"github.com/nkamuo/rasta-server/model"
	"github.com/nkamuo/rasta-server/repository"
)

var requestService RequestService
var requestRepoMutext *sync.Mutex = &sync.Mutex{}

func GetRequestService() RequestService {
	requestRepoMutext.Lock()
	if requestService == nil {
		requestService = &requestServiceImpl{repo: repository.GetRequestRepository()}
	}
	requestRepoMutext.Unlock()
	return requestService
}

type RequestService interface {
	GetById(id uuid.UUID) (request *model.Request, err error)
	// GetByEmail(email string) (request *model.Request, err error)
	// GetByPhone(phone string) (request *model.Request, err error)
	Save(request *model.Request) (err error)
	Delete(request *model.Request) (error error)
}

type requestServiceImpl struct {
	repo repository.RequestRepository
}

func (service *requestServiceImpl) GetById(id uuid.UUID) (request *model.Request, err error) {
	return service.repo.GetById(id)
}

func (service *requestServiceImpl) Save(request *model.Request) (err error) {
	return service.repo.Save(request)
}

func (service *requestServiceImpl) Delete(request *model.Request) (err error) {
	err = service.repo.Delete(request)

	return err
}

func (service *requestServiceImpl) DeleteById(id uuid.UUID) (request *model.Request, err error) {
	request, err = service.repo.DeleteById(id)
	return request, err
}
