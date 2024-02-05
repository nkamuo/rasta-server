package service

import (
	"sync"

	"github.com/google/uuid"
	"github.com/nkamuo/rasta-server/model"
	"github.com/nkamuo/rasta-server/repository"
	"github.com/stripe/stripe-go/v74"
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
	GetByRespondent(respondent *model.Respondent, preload ...string) (balance *model.RespondentAccessProductSubscription, err error)
	GetActiveStripeSubscriptionByCustomerID(customerID string) (subscription *stripe.Subscription, err error)
	Close(subscription *model.RespondentAccessProductSubscription) (err error)
	Save(subscription *model.RespondentAccessProductSubscription) (err error)
	Delete(subscription *model.RespondentAccessProductSubscription) (error error)
}

type subscriptionServiceImpl struct {
	repo repository.RespondentAccessProductSubscriptionRepository
}

func (service *subscriptionServiceImpl) GetById(id uuid.UUID, preload ...string) (subscription *model.RespondentAccessProductSubscription, err error) {
	return service.repo.GetById(id, preload...)
}

func (service subscriptionServiceImpl) GetByRespondent(respondent *model.Respondent, preload ...string) (balance *model.RespondentAccessProductSubscription, err error) {
	return service.repo.GetActiveByRespondent(*respondent, preload...)
	// return nil, errors.New("Could not resolve Access Balance for the given responder");
}

func (service *subscriptionServiceImpl) Close(subscription *model.RespondentAccessProductSubscription) (err error) {
	panic("You cannot close a subscription like this yet")
	// now := time.Now()
	// subscription.EndedAt = &now
	// *subscription.Active = false
	// return service.repo.Save(subscription)
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

func (service *subscriptionServiceImpl) GetActiveStripeSubscriptionByCustomerID(customerID string) (subscription *stripe.Subscription, err error) {
	return service.repo.GetActiveByStripeCustomerID(customerID)
}
