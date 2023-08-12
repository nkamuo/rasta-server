package repository

import (
	"sync"

	"github.com/google/uuid"
	"github.com/nkamuo/rasta-server/model"
	"gorm.io/gorm"
)

var respondentServiceReviewRepo RespondentServiceReviewRepository
var respondentServiceReviewRepoMutext *sync.Mutex = &sync.Mutex{}

func GetRespondentServiceReviewRepository() RespondentServiceReviewRepository {
	respondentServiceReviewRepoMutext.Lock()
	if respondentServiceReviewRepo == nil {
		respondentServiceReviewRepo = &respondentServiceReviewRepository{db: model.DB}
	}
	respondentServiceReviewRepoMutext.Unlock()
	return respondentServiceReviewRepo
}

type RespondentServiceReviewRepository interface {
	FindAll(page int, limit int) (respondentServiceReviews []model.RespondentServiceReview, total int64, err error)
	GetById(id uuid.UUID) (respondentServiceReview *model.RespondentServiceReview, err error)
	GetByCode(code string) (respondentServiceReview *model.RespondentServiceReview, err error)
	GetRateForFuelTypeInPlace(fuelType model.FuelType, place model.Place) (rate *model.RespondentServiceReview, err error)
	Save(respondentServiceReview *model.RespondentServiceReview) (err error)
	Delete(respondentServiceReview *model.RespondentServiceReview) (error error)
	DeleteById(id uuid.UUID) (respondentServiceReview *model.RespondentServiceReview, err error)
}

type respondentServiceReviewRepository struct {
	db *gorm.DB
}

func (repo *respondentServiceReviewRepository) FindAll(page int, limit int) (respondentServiceReviews []model.RespondentServiceReview, total int64, err error) {
	offset := (page - 1) * limit

	err = repo.db.
		Model(&model.RespondentServiceReview{}).
		Count(&total).
		Limit(limit).
		Offset(offset).
		Order("created_at desc").
		Find(&respondentServiceReviews).Error
	if err != nil {
		return
	}
	return
}

func (repo *respondentServiceReviewRepository) GetById(id uuid.UUID) (respondentServiceReview *model.RespondentServiceReview, err error) {
	if err = model.DB.Where("id = ?", id).First(&respondentServiceReview).Error; err != nil {
		return nil, err
	}
	return respondentServiceReview, nil
}
func (repo *respondentServiceReviewRepository) GetByCode(code string) (respondentServiceReview *model.RespondentServiceReview, err error) {
	if err = model.DB.Where("code = ?", code).First(&respondentServiceReview).Error; err != nil {
		return nil, err
	}
	return respondentServiceReview, nil
}

func (repo *respondentServiceReviewRepository) GetRateForFuelTypeInPlace(fuelType model.FuelType, place model.Place) (rate *model.RespondentServiceReview, err error) {
	if err = model.DB.Where("fuel_type_id = ? AND place_id = ?", fuelType.ID, place.ID).First(&rate).Error; err != nil {
		return nil, err
	}
	return rate, nil
}

func (repo *respondentServiceReviewRepository) Save(respondentServiceReview *model.RespondentServiceReview) (err error) {
	if (uuid.UUID{} == respondentServiceReview.ID) {
		//NEW - No ID yet
		return repo.db.Create(&respondentServiceReview).Error
	}
	return repo.db.Updates(&respondentServiceReview).Error
}

func (repo *respondentServiceReviewRepository) Delete(respondentServiceReview *model.RespondentServiceReview) (err error) {
	return repo.db.Delete(&respondentServiceReview).Error
}

func (repo *respondentServiceReviewRepository) DeleteById(id uuid.UUID) (respondentServiceReview *model.RespondentServiceReview, err error) {
	respondentServiceReview, err = repo.GetById(id)
	if err != nil {
		return nil, err
	}
	err = repo.db.Delete(respondentServiceReview).Error
	return respondentServiceReview, err
}
