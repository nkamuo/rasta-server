package repository

import (
	"sync"

	"github.com/google/uuid"
	"github.com/nkamuo/rasta-server/model"
	"gorm.io/gorm"
)

var sessionRepo RespondentSessionRepository
var sessionRepoMutext *sync.Mutex = &sync.Mutex{}

func GetRespondentSessionRepository() RespondentSessionRepository {
	sessionRepoMutext.Lock()
	if sessionRepo == nil {
		sessionRepo = &sessionRepository{db: model.DB}
	}
	sessionRepoMutext.Unlock()
	return sessionRepo
}

type RespondentSessionRepository interface {
	FindAll(page int, limit int) (sessions []model.RespondentSession, total int64, err error)
	GetById(id uuid.UUID) (session *model.RespondentSession, err error)
	GetActiveByRespondent(respondent model.Respondent, prefetch ...string) (session *model.RespondentSession, err error)
	Save(session *model.RespondentSession) (err error)
	Delete(session *model.RespondentSession) (error error)
	DeleteById(id uuid.UUID) (session *model.RespondentSession, err error)
}

type sessionRepository struct {
	db *gorm.DB
}

func (repo *sessionRepository) FindAll(page int, limit int) (sessions []model.RespondentSession, total int64, err error) {
	offset := (page - 1) * limit

	err = repo.db.
		Model(&model.RespondentSession{}).
		Count(&total).
		Limit(limit).
		Offset(offset).
		Order("created_at desc").
		Find(&sessions).Error
	if err != nil {
		return
	}
	return
}

func (repo *sessionRepository) GetById(id uuid.UUID) (session *model.RespondentSession, err error) {
	if err = model.DB. /*.Joins("OperatorUser")*/ First(&session, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return session, nil
}

func (repo *sessionRepository) GetActiveByRespondent(respondent model.Respondent, prefetch ...string) (session *model.RespondentSession, err error) {
	query := model.DB //.Preload("Assignments.Assignment.Product").
	query = query.Where("respondent_id = ? AND active = ? AND ended_at IS NULL", respondent.ID, true)

	for _, pFetch := range prefetch {
		query = query.Preload(pFetch)
	}

	if err = query.First(&session).Error; err != nil {
		return nil, err
	}
	return session, nil
}

func (repo *sessionRepository) Save(session *model.RespondentSession) (err error) {
	if (uuid.UUID{} == session.ID) {
		//NEW - No ID yet
		return repo.db.Create(&session).Error
	}
	return repo.db.Updates(&session).Error
}

func (repo *sessionRepository) Delete(session *model.RespondentSession) (err error) {
	return repo.db.Delete(&session).Error
}

func (repo *sessionRepository) DeleteById(id uuid.UUID) (session *model.RespondentSession, err error) {
	session, err = repo.GetById(id)
	if err != nil {
		return nil, err
	}
	err = repo.db.Delete(&session).Error
	return session, err
}
