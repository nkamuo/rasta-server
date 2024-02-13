package repository

import (
	"sync"

	"github.com/google/uuid"
	"github.com/nkamuo/rasta-server/model"
	"gorm.io/gorm"
)

var respondentRepo RespondentRepository
var respondentRepoMutext *sync.Mutex = &sync.Mutex{}

func GetRespondentRepository() RespondentRepository {
	respondentRepoMutext.Lock()
	if respondentRepo == nil {
		respondentRepo = &respondentRepository{db: model.DB}
	}
	respondentRepoMutext.Unlock()
	return respondentRepo
}

type RespondentRepository interface {
	FindAll(input ...int) (respondents []model.Respondent, total int64, err error)
	FindAllByCompanyID(companyID uuid.UUID, page int, limit int) (respondents []model.Respondent, total int64, err error)
	GetById(id uuid.UUID, preload ...string) (respondent *model.Respondent, err error)
	GetByEmail(email string, preload ...string) (respondent *model.Respondent, err error)
	GetByPhone(phone string, preload ...string) (respondent *model.Respondent, err error)
	GetByUser(user model.User, preload ...string) (respondent *model.Respondent, err error)
	GetByUserId(userID uuid.UUID, preload ...string) (respondent *model.Respondent, err error)
	Save(respondent *model.Respondent) (err error)
	Delete(respondent *model.Respondent) (error error)
	DeleteById(id uuid.UUID) (respondent *model.Respondent, err error)
}

type respondentRepository struct {
	db *gorm.DB
}

func (repo *respondentRepository) FindAll(input ...int) (respondents []model.Respondent, total int64, err error) {
	page := 1
	limit := 10
	if len(input) > 0 {
		page = input[0]
	}
	if len(input) > 1 {
		limit = input[1]
	}
	offset := (page - 1) * limit

	err = repo.db.
		Model(&model.Respondent{}).
		Count(&total).
		Limit(limit).
		Offset(offset).
		Order("created_at desc").
		Find(&respondents).Error
	if err != nil {
		return
	}
	return
}

func (repo *respondentRepository) FindAllByCompanyID(companyID uuid.UUID, page int, limit int) (respondents []model.Respondent, total int64, err error) {
	offset := (page - 1) * limit

	err = repo.db.
		Joins("JOIN companies ON companies.id = respondent.company_id").
		Where("companies.id = ?", companyID).
		Model(&model.Respondent{}).
		Count(&total).
		Limit(limit).
		Offset(offset).
		Order("created_at desc").
		Find(&respondents).Error
	if err != nil {
		return
	}
	return
}

func (repo *respondentRepository) GetById(id uuid.UUID, preload ...string) (respondent *model.Respondent, err error) {

	query := model.DB
	for _, pLoad := range preload {
		query = query.Preload(pLoad)
	}

	if err = query.First(&respondent, "respondents.id = ?", id).Error; err != nil {
		return nil, err
	}
	return respondent, nil
}

func (repo *respondentRepository) GetByEmail(email string, preload ...string) (respondent *model.Respondent, err error) {
	query := model.DB

	for _, pLoad := range preload {
		query = query.Preload(pLoad)
	}

	if err = query.Where("email = ?", email).First(&respondent).Error; err != nil {
		return nil, err
	}
	return respondent, nil
}

func (repo *respondentRepository) GetByPhone(phone string, preload ...string) (respondent *model.Respondent, err error) {

	query := model.DB
	for _, pLoad := range preload {
		query = query.Preload(pLoad)
	}

	if err = query.Where("phone = ?", phone).First(&respondent).Error; err != nil {
		return nil, err
	}
	return respondent, nil
}

func (repo *respondentRepository) GetByUserId(userId uuid.UUID, preload ...string) (respondent *model.Respondent, err error) {

	query := model.DB.
		Joins("JOIN users on users.id = respondents.user_id").
		Where("users.id = ?", userId)

	for _, pLoad := range preload {
		query = query.Preload(pLoad)
	}

	if err = query.First(&respondent).Error; err != nil {
		return nil, err
	}
	return respondent, nil
}

func (repo *respondentRepository) GetByUser(user model.User, preload ...string) (respondent *model.Respondent, err error) {
	return repo.GetByUserId(user.ID, preload...)
}

func (repo *respondentRepository) Save(respondent *model.Respondent) (err error) {
	if (uuid.UUID{} == respondent.ID) {
		//NEW - No ID yet
		return repo.db.Create(&respondent).Error
	}
	return repo.db.Updates(&respondent).Error
}

func (repo *respondentRepository) Delete(respondent *model.Respondent) (err error) {
	return repo.db.Delete(&respondent).Error
}

func (repo *respondentRepository) DeleteById(id uuid.UUID) (respondent *model.Respondent, err error) {
	respondent, err = repo.GetById(id)
	if err != nil {
		return nil, err
	}
	err = repo.db.Delete(&respondent).Error
	return respondent, err
}
