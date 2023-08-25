package service

import (
	"sync"

	"github.com/google/uuid"
	"github.com/nkamuo/rasta-server/model"
	"github.com/nkamuo/rasta-server/repository"
)

var motoristRequestSituationService MotoristRequestSituationService
var motoristRequestsituationRepoMutext *sync.Mutex = &sync.Mutex{}

func GetMotoristRequestSituationService() MotoristRequestSituationService {
	motoristRequestsituationRepoMutext.Lock()
	if motoristRequestSituationService == nil {
		motoristRequestSituationService = &motoristRequestSituationServiceImpl{repo: repository.GetMotoristRequestSituationRepository()}
	}
	motoristRequestsituationRepoMutext.Unlock()
	return motoristRequestSituationService
}

type MotoristRequestSituationService interface {
	GetById(id uuid.UUID) (motoristRequestsituation *model.MotoristRequestSituation, err error)
	// GetByEmail(email string) (motoristRequestsituation *model.MotoristRequestSituation, err error)
	// GetByPhone(phone string) (motoristRequestsituation *model.MotoristRequestSituation, err error)
	Save(motoristRequestsituation *model.MotoristRequestSituation) (err error)
	Delete(motoristRequestsituation *model.MotoristRequestSituation) (error error)
}

type motoristRequestSituationServiceImpl struct {
	repo repository.MotoristRequestSituationRepository
}

func (service *motoristRequestSituationServiceImpl) GetById(id uuid.UUID) (motoristRequestsituation *model.MotoristRequestSituation, err error) {
	return service.repo.GetById(id)
}

func (service *motoristRequestSituationServiceImpl) Save(motoristRequestsituation *model.MotoristRequestSituation) (err error) {
	return service.repo.Save(motoristRequestsituation)
}

func (service *motoristRequestSituationServiceImpl) Delete(motoristRequestsituation *model.MotoristRequestSituation) (err error) {
	err = service.repo.Delete(motoristRequestsituation)

	return err
}

func (service *motoristRequestSituationServiceImpl) DeleteById(id uuid.UUID) (motoristRequestsituation *model.MotoristRequestSituation, err error) {
	motoristRequestsituation, err = service.repo.DeleteById(id)
	return motoristRequestsituation, err
}
