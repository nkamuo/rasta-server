package service

import (
	"sync"

	"github.com/google/uuid"
	"github.com/nkamuo/rasta-server/model"
	"github.com/nkamuo/rasta-server/repository"
)

var respondentServiceReviewService RespondentServiceReviewService
var respondentServiceReviewRepoMutext *sync.Mutex = &sync.Mutex{}

func GetRespondentServiceReviewService() RespondentServiceReviewService {
	respondentServiceReviewRepoMutext.Lock()
	if respondentServiceReviewService == nil {
		respondentServiceReviewService = &respondentServiceReviewServiceImpl{repo: repository.GetRespondentServiceReviewRepository()}
	}
	respondentServiceReviewRepoMutext.Unlock()
	return respondentServiceReviewService
}

type RespondentServiceReviewService interface {
	GetById(id uuid.UUID) (respondentServiceReview *model.RespondentServiceReview, err error)
	Save(respondentServiceReview *model.RespondentServiceReview) (err error)
	Delete(respondentServiceReview *model.RespondentServiceReview) (error error)
}

type respondentServiceReviewServiceImpl struct {
	repo repository.RespondentServiceReviewRepository
}

func (service *respondentServiceReviewServiceImpl) GetById(id uuid.UUID) (respondentServiceReview *model.RespondentServiceReview, err error) {
	return service.repo.GetById(id)
}

func (service *respondentServiceReviewServiceImpl) Save(respondentServiceReview *model.RespondentServiceReview) (err error) {
	return service.repo.Save(respondentServiceReview)
}

func (service *respondentServiceReviewServiceImpl) Delete(respondentServiceReview *model.RespondentServiceReview) (err error) {
	err = service.repo.Delete(respondentServiceReview)

	return err
}

func (service *respondentServiceReviewServiceImpl) DeleteById(id uuid.UUID) (respondentServiceReview *model.RespondentServiceReview, err error) {
	respondentServiceReview, err = service.repo.DeleteById(id)
	return respondentServiceReview, err
}
