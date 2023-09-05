package service

import (
	"sync"

	"github.com/google/uuid"
	"github.com/nkamuo/rasta-server/model"
	"github.com/nkamuo/rasta-server/repository"
)

var respondentEarningService RespondentEarningService
var respondentEarningRepoMutext *sync.Mutex = &sync.Mutex{}

func GetRespondentEarningService() RespondentEarningService {
	respondentEarningRepoMutext.Lock()
	if respondentEarningService == nil {
		respondentEarningService = &respondentEarningServiceImpl{repo: repository.GetRespondentEarningRepository()}
	}
	respondentEarningRepoMutext.Unlock()
	return respondentEarningService
}

type RespondentEarningService interface {
	GetById(id uuid.UUID) (respondentEarning *model.RespondentEarning, err error)
	// GetByEmail(email string) (respondentEarning *model.RespondentEarning, err error)
	// GetByPhone(phone string) (respondentEarning *model.RespondentEarning, err error)
	Save(respondentEarning *model.RespondentEarning) (err error)
	Delete(respondentEarning *model.RespondentEarning) (error error)
}

type respondentEarningServiceImpl struct {
	repo repository.RespondentEarningRepository
}

func (service *respondentEarningServiceImpl) GetById(id uuid.UUID) (respondentEarning *model.RespondentEarning, err error) {
	return service.repo.GetById(id)
}

func (service *respondentEarningServiceImpl) Save(respondentEarning *model.RespondentEarning) (err error) {
	return service.repo.Save(respondentEarning)
}

func (service *respondentEarningServiceImpl) Delete(respondentEarning *model.RespondentEarning) (err error) {
	err = service.repo.Delete(respondentEarning)

	return err
}

func (service *respondentEarningServiceImpl) DeleteById(id uuid.UUID) (respondentEarning *model.RespondentEarning, err error) {
	respondentEarning, err = service.repo.DeleteById(id)
	return respondentEarning, err
}
