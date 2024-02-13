package service

import (
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/nkamuo/rasta-server/model"
	"github.com/nkamuo/rasta-server/repository"
	// "github.com/stripe/stripe-go/v74/subscription"
)

var subscriptionService RespondentAccessProductSubscriptionService
var subscriptionRepoMutext *sync.Mutex = &sync.Mutex{}

func GetRespondentAccessProductSubscriptionService() RespondentAccessProductSubscriptionService {
	subscriptionRepoMutext.Lock()
	if subscriptionService == nil {
		subscriptionService = &subscriptionServiceImpl{repo: repository.GetRespondentAccessProductSubscriptionRepository()}
	}
	subscriptionRepoMutext.Unlock()
	return subscriptionService
}

type RespondentAccessProductSubscriptionService interface {
	GetById(id uuid.UUID, preload ...string) (subscription *model.RespondentAccessProductSubscription, err error)
	GetByRespondent(respondent *model.Respondent, preload ...string) (subscription *model.RespondentAccessProductSubscription, err error)
	GetActiveByRespondentAndTime(respondent model.Respondent, time time.Time, prefetch ...string) (sub *model.RespondentAccessProductSubscription, err error)
	ExtendByDuration(subscription *model.RespondentAccessProductSubscription, duration time.Duration) (err error)
	//
	SetupForRespondent(respondent *model.Respondent) (subscription *model.RespondentAccessProductSubscription, err error)
	SetupForAllRespondents() (err error)
	//
	ExtendByDays(subscription *model.RespondentAccessProductSubscription, days int64) (err error)
	Save(subscription *model.RespondentAccessProductSubscription) (err error)
	Delete(subscription *model.RespondentAccessProductSubscription) (error error)
}

type subscriptionServiceImpl struct {
	repo repository.RespondentAccessProductSubscriptionRepository
}

func (service *subscriptionServiceImpl) GetById(id uuid.UUID, preload ...string) (subscription *model.RespondentAccessProductSubscription, err error) {
	return service.repo.GetById(id, preload...)
}

func (service subscriptionServiceImpl) GetByRespondent(respondent *model.Respondent, preload ...string) (subscription *model.RespondentAccessProductSubscription, err error) {
	return service.repo.GetByRespondent(*respondent, preload...)
	// return nil, errors.New("Could not resolve Access Balance for the given responder");
}

func (service *subscriptionServiceImpl) ExtendByDuration(subscription *model.RespondentAccessProductSubscription, duration time.Duration) (err error) {
	subscription.ExtendByDuration(duration)
	return service.repo.Save(subscription)
}

func (service *subscriptionServiceImpl) ExtendByDays(subscription *model.RespondentAccessProductSubscription, days int64) (err error) {
	duration := time.Duration(days) * 24 * time.Hour
	return service.ExtendByDuration(subscription, duration)
}

func (service *subscriptionServiceImpl) Save(subscription *model.RespondentAccessProductSubscription) (err error) {
	return service.repo.Save(subscription)
}

func (service *subscriptionServiceImpl) Delete(subscription *model.RespondentAccessProductSubscription) (err error) {
	err = service.repo.Delete(subscription)

	return err
}

func (service *subscriptionServiceImpl) DeleteById(id uuid.UUID) (subscription *model.RespondentAccessProductSubscription, err error) {
	subscription, err = service.repo.DeleteById(id)
	return subscription, err
}

//

func (service *subscriptionServiceImpl) GetActiveByRespondentAndTime(respondent model.Respondent, time time.Time, prefetch ...string) (sub *model.RespondentAccessProductSubscription, err error) {
	return service.repo.GetActiveByRespondentAndTime(respondent, time, prefetch...)
}

func (service subscriptionServiceImpl) SetupForRespondent(respondent *model.Respondent) (subscription *model.RespondentAccessProductSubscription, err error) {

	subscription, err = service.GetByRespondent(respondent)
	if err != nil {
		if err.Error() == "record not found" {

		} else {
			return nil, err
		}
	}

	if subscription == nil {
		subscription = &model.RespondentAccessProductSubscription{
			RespondentID: &respondent.ID,
			StartingAt:   time.Now(),
			EndingAt:     time.Now(),
			// Balance:      stripe.Int64(0),
		}

		if err = service.Save(subscription); err != nil {
			return nil, err
		}
	}
	return subscription, err
}

func (service subscriptionServiceImpl) SetupForAllRespondents() (err error) {
	respondantRepo := repository.GetRespondentRepository()
	respondents, _, err := respondantRepo.FindAll(1, 100000)
	if err != nil {
		return err
	}
	for _, respondent := range respondents {
		_, err = service.SetupForRespondent(&respondent)
		if err != nil {
			return err
		}
	}
	return err
}
