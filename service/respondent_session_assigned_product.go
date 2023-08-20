package service

import (
	"sync"

	"github.com/google/uuid"
	"github.com/nkamuo/rasta-server/model"
	"github.com/nkamuo/rasta-server/repository"
)

var assignedProductService RespondentSessionAssignedProductService
var assignedProductRepoMutext *sync.Mutex = &sync.Mutex{}

func GetRespondentSessionAssignedProductService() RespondentSessionAssignedProductService {
	assignedProductRepoMutext.Lock()
	if assignedProductService == nil {
		assignedProductService = &assignedProductServiceImpl{repo: repository.GetRespondentSessionAssignedProductRepository()}
	}
	assignedProductRepoMutext.Unlock()
	return assignedProductService
}

type RespondentSessionAssignedProductService interface {
	GetById(id uuid.UUID) (assignedProduct *model.RespondentSessionAssignedProduct, err error)
	Save(assignedProduct *model.RespondentSessionAssignedProduct) (err error)
	Delete(assignedProduct *model.RespondentSessionAssignedProduct) (error error)
}

type assignedProductServiceImpl struct {
	repo repository.RespondentSessionAssignedProductRepository
}

func (service *assignedProductServiceImpl) GetById(id uuid.UUID) (assignedProduct *model.RespondentSessionAssignedProduct, err error) {
	return service.repo.GetById(id)
}

func (service *assignedProductServiceImpl) Save(assignedProduct *model.RespondentSessionAssignedProduct) (err error) {
	return service.repo.Save(assignedProduct)
}

func (service *assignedProductServiceImpl) Delete(assignedProduct *model.RespondentSessionAssignedProduct) (err error) {
	err = service.repo.Delete(assignedProduct)

	return err
}

func (service *assignedProductServiceImpl) DeleteById(id uuid.UUID) (assignedProduct *model.RespondentSessionAssignedProduct, err error) {
	assignedProduct, err = service.repo.DeleteById(id)
	return assignedProduct, err
}
