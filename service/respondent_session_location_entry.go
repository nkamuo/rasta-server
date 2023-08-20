package service

import (
	"sync"

	"github.com/google/uuid"
	"github.com/nkamuo/rasta-server/model"
	"github.com/nkamuo/rasta-server/repository"
)

var locationEntryService RespondentSessionLocationEntryService
var locationEntryRepoMutext *sync.Mutex = &sync.Mutex{}

func GetRespondentSessionLocationEntryService() RespondentSessionLocationEntryService {
	locationEntryRepoMutext.Lock()
	if locationEntryService == nil {
		locationEntryService = &locationEntryServiceImpl{repo: repository.GetRespondentSessionLocationEntryRepository()}
	}
	locationEntryRepoMutext.Unlock()
	return locationEntryService
}

type RespondentSessionLocationEntryService interface {
	GetById(id uuid.UUID) (locationEntry *model.RespondentSessionLocationEntry, err error)
	Save(locationEntry *model.RespondentSessionLocationEntry) (err error)
	Delete(locationEntry *model.RespondentSessionLocationEntry) (error error)
}

type locationEntryServiceImpl struct {
	repo repository.RespondentSessionLocationEntryRepository
}

func (service *locationEntryServiceImpl) GetById(id uuid.UUID) (locationEntry *model.RespondentSessionLocationEntry, err error) {
	return service.repo.GetById(id)
}

func (service *locationEntryServiceImpl) Save(locationEntry *model.RespondentSessionLocationEntry) (err error) {
	return service.repo.Save(locationEntry)
}

func (service *locationEntryServiceImpl) Delete(locationEntry *model.RespondentSessionLocationEntry) (err error) {
	err = service.repo.Delete(locationEntry)

	return err
}

func (service *locationEntryServiceImpl) DeleteById(id uuid.UUID) (locationEntry *model.RespondentSessionLocationEntry, err error) {
	locationEntry, err = service.repo.DeleteById(id)
	return locationEntry, err
}
