package service

import (
	"errors"
	"fmt"
	"sync"

	"github.com/google/uuid"
	"github.com/nkamuo/rasta-server/model"
	"github.com/nkamuo/rasta-server/repository"
	"github.com/stripe/stripe-go/v74/price"
	"gorm.io/gorm"
)

var purchaseService RespondentAccessProductPurchaseService
var purchaseRepoMutext *sync.Mutex = &sync.Mutex{}

func GetRespondentAccessProductPurchaseService() RespondentAccessProductPurchaseService {
	purchaseRepoMutext.Lock()
	defer purchaseRepoMutext.Unlock()
	if purchaseService == nil {
		purchaseService = &purchaseServiceImpl{repo: repository.GetRespondentAccessProductPurchaseRepository()}
	}
	return purchaseService
}

type RespondentAccessProductPurchaseService interface {
	GetById(id uuid.UUID, preload ...string) (purchase *model.RespondentAccessProductPurchase, err error)
	GetByStripeCheckoutID(id string, preload ...string) (purchase *model.RespondentAccessProductPurchase, err error)
	// GetByRespondent(respondent *model.Respondent, preload ...string) (purchase *model.RespondentAccessProductPurchase, err error)
	// Close(purchase *model.RespondentAccessProductPurchase) (err error)
	Commit(purchase *model.RespondentAccessProductPurchase) (err error)
	Cancel(purchase *model.RespondentAccessProductPurchase) (err error)
	Save(purchase *model.RespondentAccessProductPurchase) (err error)
	Delete(purchase *model.RespondentAccessProductPurchase) (error error)
}

type purchaseServiceImpl struct {
	repo repository.RespondentAccessProductPurchaseRepository
}

func (service *purchaseServiceImpl) GetById(id uuid.UUID, preload ...string) (purchase *model.RespondentAccessProductPurchase, err error) {
	return service.repo.GetById(id, preload...)
}

func (service *purchaseServiceImpl) GetByStripeCheckoutID(id string, preload ...string) (purchase *model.RespondentAccessProductPurchase, err error) {
	return service.repo.GetByStripeCheckoutID(id, preload...)
}

// func (service *purchaseServiceImpl) Close(purchase *model.RespondentAccessProductPurchase) (err error) {
// 	now := time.Now()
// 	purchase.EndedAt = &now
// 	*purchase.Active = false
// 	return service.repo.Save(purchase)
// }

// func (service purchaseServiceImpl) GetByRespondent(respondent *model.Respondent, preload ...string) (purchase *model.RespondentAccessProductPurchase, err error) {
// 	return service.repo.GetActiveByRespondent(*respondent, preload...)
// 	// return nil, errors.New("Could not resolve Access Purchase for the given responder");
// }

func (service *purchaseServiceImpl) Save(purchase *model.RespondentAccessProductPurchase) (err error) {
	return service.repo.Save(purchase)
}

func (service *purchaseServiceImpl) Commit(purchase *model.RespondentAccessProductPurchase) (err error) {

	respondentService := GetRespondentService()
	balanceService := GetRespondentAccessProductBalanceService()

	respondent, err := respondentService.GetById(*purchase.RespondentID)
	if err != nil {
		return nil
	}

	balance, err := balanceService.SetupForRespondent(respondent)
	if err != nil {
		return err
	}

	if purchase.Succeeded != nil && *purchase.Succeeded {
		message := fmt.Sprintf("Cannot commit and already succeed purchase")
		return errors.New(message)
	}

	if purchase.Cancelled != nil && *purchase.Cancelled {
		message := fmt.Sprintf("Purchase is alrady cancelled")
		return errors.New(message)
	}

	if purchase.StripePriceID != nil {
		stripePrice, err := price.Get(*purchase.StripePriceID, nil)
		if err != nil {
			return err
		}
		if stripePrice.TransformQuantity != nil {

			var multiplier = stripePrice.TransformQuantity.DivideBy
			balance.Increment(multiplier)
		}
	}
	*purchase.Succeeded = true
	err = model.DB.Transaction(func(tx *gorm.DB) error {
		if err = tx.Save(purchase).Error; err != nil {
			return err
		}
		if err = tx.Save(balance).Error; err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return err
	}
	return nil
}

func (service *purchaseServiceImpl) Cancel(purchase *model.RespondentAccessProductPurchase) (err error) {

	if purchase.Cancelled != nil && *purchase.Cancelled {
		message := fmt.Sprintf("Purchase is alrady cancelled")
		return errors.New(message)
	}
	if purchase.Succeeded != nil && *purchase.Succeeded {
		message := fmt.Sprintf("Cannot commit and already succeed purchase")
		return errors.New(message)
	}
	*purchase.Cancelled = true

	if err = service.Save(purchase); err != nil {
		return err
	}
	return nil
}

func (service *purchaseServiceImpl) Delete(purchase *model.RespondentAccessProductPurchase) (err error) {
	err = service.repo.Delete(purchase)

	return err
}

func (service *purchaseServiceImpl) DeleteById(id uuid.UUID) (purchase *model.RespondentAccessProductPurchase, err error) {
	purchase, err = service.repo.DeleteById(id)
	return purchase, err
}
