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
	FindAll(page int, limit int) (respondents []model.Respondent, total int64, err error)
	FindAllByCompanyID(companyID uuid.UUID, page int, limit int) (respondents []model.Respondent, total int64, err error)
	GetById(id uuid.UUID) (respondent *model.Respondent, err error)
	GetByEmail(email string) (respondent *model.Respondent, err error)
	GetByPhone(phone string) (respondent *model.Respondent, err error)
	GetByUser(user model.User) (respondent *model.Respondent, err error)
	GetByUserId(userID uuid.UUID) (respondent *model.Respondent, err error)
	Save(respondent *model.Respondent) (err error)
	Delete(respondent *model.Respondent) (error error)
	DeleteById(id uuid.UUID) (respondent *model.Respondent, err error)
}

type respondentRepository struct {
	db *gorm.DB
}

func (repo *respondentRepository) FindAll(page int, limit int) (respondents []model.Respondent, total int64, err error) {
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

func (repo *respondentRepository) GetById(id uuid.UUID) (respondent *model.Respondent, err error) {
	if err = model.DB. /*.Joins("User")*/ First(&respondent, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return respondent, nil
}

func (repo *respondentRepository) GetByEmail(email string) (respondent *model.Respondent, err error) {
	if err = model.DB.Where("email = ?", email).First(&respondent).Error; err != nil {
		return nil, err
	}
	return respondent, nil
}

func (repo *respondentRepository) GetByPhone(phone string) (respondent *model.Respondent, err error) {
	if err = model.DB.Where("phone = ?", phone).First(&respondent).Error; err != nil {
		return nil, err
	}
	return respondent, nil
}

func (repo *respondentRepository) GetByUserId(userId uuid.UUID) (respondent *model.Respondent, err error) {
	if err = model.DB.
		Joins("JOIN users on users.id = respondents.user_id").
		Where("users.id = ?", userId).First(&respondent).Error; err != nil {
		return nil, err
	}
	return respondent, nil
}

func (repo *respondentRepository) GetByUser(user model.User) (respondent *model.Respondent, err error) {
	return repo.GetByUserId(user.ID)
}

func (repo *respondentRepository) Save(respondent *model.Respondent) (err error) {
	if (uuid.UUID{} == respondent.ID) {
		//NEW - No ID yet
		repo.db.Create(&respondent)
		return repo.db.Error
	}
	repo.db.Updates(&respondent)
	return repo.db.Error
}

func (repo *respondentRepository) Delete(respondent *model.Respondent) (err error) {
	repo.db.Delete(&respondent)
	return repo.db.Error
}

func (repo *respondentRepository) DeleteById(id uuid.UUID) (respondent *model.Respondent, err error) {
	respondent, err = repo.GetById(id)
	if err != nil {
		return nil, err
	}
	repo.db.Delete(&respondent)
	return respondent, repo.db.Error
}
