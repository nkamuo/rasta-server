package repository

import (
	"sync"

	"github.com/google/uuid"
	"github.com/nkamuo/rasta-server/model"
	"gorm.io/gorm"
)

var locationEntryRepo RespondentSessionLocationEntryRepository
var locationEntryRepoMutext *sync.Mutex = &sync.Mutex{}

func GetRespondentSessionLocationEntryRepository() RespondentSessionLocationEntryRepository {
	locationEntryRepoMutext.Lock()
	if locationEntryRepo == nil {
		locationEntryRepo = &locationEntryRepository{db: model.DB}
	}
	locationEntryRepoMutext.Unlock()
	return locationEntryRepo
}

type RespondentSessionLocationEntryRepository interface {
	FindAll(page int, limit int) (locationEntrys []model.RespondentSessionLocationEntry, total int64, err error)
	GetById(id uuid.UUID) (locationEntry *model.RespondentSessionLocationEntry, err error)
	Save(locationEntry *model.RespondentSessionLocationEntry) (err error)
	Delete(locationEntry *model.RespondentSessionLocationEntry) (error error)
	DeleteById(id uuid.UUID) (locationEntry *model.RespondentSessionLocationEntry, err error)
}

type locationEntryRepository struct {
	db *gorm.DB
}

func (repo *locationEntryRepository) FindAll(page int, limit int) (locationEntrys []model.RespondentSessionLocationEntry, total int64, err error) {
	offset := (page - 1) * limit

	err = repo.db.
		Model(&model.RespondentSessionLocationEntry{}).
		Count(&total).
		Limit(limit).
		Offset(offset).
		Order("created_at desc").
		Find(&locationEntrys).Error
	if err != nil {
		return
	}
	return
}

func (repo *locationEntryRepository) GetById(id uuid.UUID) (locationEntry *model.RespondentSessionLocationEntry, err error) {
	if err = model.DB. /*.Joins("OperatorUser")*/ First(&locationEntry, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return locationEntry, nil
}

func (repo *locationEntryRepository) Save(locationEntry *model.RespondentSessionLocationEntry) (err error) {
	if (uuid.UUID{} == locationEntry.ID) {
		//NEW - No ID yet
		return repo.db.Create(&locationEntry).Error
	}
	return repo.db.Updates(&locationEntry).Error
}

func (repo *locationEntryRepository) Delete(locationEntry *model.RespondentSessionLocationEntry) (err error) {
	return repo.db.Delete(&locationEntry).Error
}

func (repo *locationEntryRepository) DeleteById(id uuid.UUID) (locationEntry *model.RespondentSessionLocationEntry, err error) {
	locationEntry, err = repo.GetById(id)
	if err != nil {
		return nil, err
	}
	err = repo.db.Delete(&locationEntry).Error
	return locationEntry, err
}
