package repository

import (
	"errors"
	"fmt"
	"sync"

	"github.com/google/uuid"
	"github.com/nkamuo/rasta-server/initializers"
	"github.com/nkamuo/rasta-server/model"
	"github.com/stripe/stripe-go/v74"
	"github.com/stripe/stripe-go/v74/subscription"
	"gorm.io/gorm"
)

var subscriptionRepo RespondentAccessProductSubscriptionRepository
var subscriptionRepoMutext *sync.Mutex = &sync.Mutex{}

func GetRespondentAccessProductSubscriptionRepository() RespondentAccessProductSubscriptionRepository {
	subscriptionRepoMutext.Lock()
	if subscriptionRepo == nil {
		subscriptionRepo = &subscriptionRepository{db: model.DB}
	}
	subscriptionRepoMutext.Unlock()
	return subscriptionRepo
}

type RespondentAccessProductSubscriptionRepository interface {
	FindAll(page int, limit int) (subscriptions []model.RespondentAccessProductSubscription, total int64, err error)
	GetById(id uuid.UUID, prefetch ...string) (subscription *model.RespondentAccessProductSubscription, err error)
	GetActiveByRespondent(respondent model.Respondent, prefetch ...string) (subscription *model.RespondentAccessProductSubscription, err error)
	GetActiveByStripeCustomerID(customerID string) (sub *stripe.Subscription, err error)
	Save(subscription *model.RespondentAccessProductSubscription) (err error)
	Delete(subscription *model.RespondentAccessProductSubscription) (error error)
	DeleteById(id uuid.UUID) (subscription *model.RespondentAccessProductSubscription, err error)
}

type subscriptionRepository struct {
	db *gorm.DB
}

func (repo *subscriptionRepository) FindAll(page int, limit int) (subscriptions []model.RespondentAccessProductSubscription, total int64, err error) {
	offset := (page - 1) * limit

	err = repo.db.
		Model(&model.RespondentAccessProductSubscription{}).
		Count(&total).
		Limit(limit).
		Offset(offset).
		Order("created_at desc").
		Find(&subscriptions).Error
	if err != nil {
		return
	}
	return
}

func (repo *subscriptionRepository) GetById(id uuid.UUID, prefetch ...string) (subscription *model.RespondentAccessProductSubscription, err error) {
	query := model.DB

	for _, pFetch := range prefetch {
		query = query.Preload(pFetch)
	}

	if err = query.First(&subscription, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return subscription, nil
}

func (repo *subscriptionRepository) GetActiveByRespondent(respondent model.Respondent, prefetch ...string) (subscription *model.RespondentAccessProductSubscription, err error) {
	query := model.DB //.Preload("Assignments.Assignment.Product").
	query = query.Where("respondent_id = ? AND active = ? AND ended_at IS NULL", respondent.ID, true)

	for _, pFetch := range prefetch {
		query = query.Preload(pFetch)
	}

	if err = query.First(&subscription).Error; err != nil {
		return nil, err
	}
	return subscription, nil
}

func (repo *subscriptionRepository) GetActiveByStripeCustomerID(customerID string) (sub *stripe.Subscription, err error) {

	// params := &stripe.SubscriptionSearchParams{
	// 	SearchParams: stripe.SearchParams{
	// 	  Query: "status:'active' AND metadata['order_id']:'6735'",

	// 	},
	//   };

	config, err := initializers.LoadConfig()
	if err != nil {
		return nil, err
	}

	params := &stripe.SubscriptionListParams{
		Customer: stripe.String(customerID),
		Status:   stripe.String("active"),
	}
	iter := subscription.List(params)

	for iter.Next() {
		sub := iter.Subscription()

		for _, item := range sub.Items.Data {
			if item.Price.Product.ID != config.STRIPE_RESPONDENT_SUBSCRIPTION_PRODUCT_ID {
				continue
			}
			return sub, nil
			// fmt.Println(item.Price.ID)
		}
		// prices = append(prices, iter.Price())
	}

	if iter.Err() != nil {
		message := fmt.Sprintf(" %s", iter.Err().Error())
		return nil, errors.New(message)
	}

	return nil, errors.New("Not implemented")
}

func (repo *subscriptionRepository) Save(subscription *model.RespondentAccessProductSubscription) (err error) {
	if (uuid.UUID{} == subscription.ID) {
		//NEW - No ID yet
		return repo.db.Create(&subscription).Error
	}
	return repo.db.Updates(&subscription).Error
}

func (repo *subscriptionRepository) Delete(subscription *model.RespondentAccessProductSubscription) (err error) {
	return repo.db.Delete(&subscription).Error
}

func (repo *subscriptionRepository) DeleteById(id uuid.UUID) (subscription *model.RespondentAccessProductSubscription, err error) {
	subscription, err = repo.GetById(id)
	if err != nil {
		return nil, err
	}
	err = repo.db.Delete(&subscription).Error
	return subscription, err
}
