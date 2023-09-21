package service

import (
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/nkamuo/rasta-server/model"
	"github.com/nkamuo/rasta-server/repository"
	"github.com/stripe/stripe-go/v74"
	"github.com/stripe/stripe-go/v74/payout"
	"gorm.io/gorm"
)

var respondentEarningService RespondentEarningService
var respondentEarningRepoMutext *sync.Mutex = &sync.Mutex{}

func GetRespondentEarningService() RespondentEarningService {
	respondentEarningRepoMutext.Lock()
	if respondentEarningService == nil {
		respondentEarningService = &respondentEarningServiceImpl{repo: repository.GetRespondentEarningRepository()}
	}
	respondentEarningRepoMutext.Unlock()
	return respondentEarningService
}

type RespondentEarningService interface {
	GetById(id uuid.UUID) (respondentEarning *model.RespondentEarning, err error)
	// GetByEmail(email string) (respondentEarning *model.RespondentEarning, err error)
	// GetByPhone(phone string) (respondentEarning *model.RespondentEarning, err error)'
	Commit(respondentEarning *model.RespondentEarning) (err error)
	Save(respondentEarning *model.RespondentEarning) (err error)
	Delete(respondentEarning *model.RespondentEarning) (error error)
}

type respondentEarningServiceImpl struct {
	repo repository.RespondentEarningRepository
}

func (service *respondentEarningServiceImpl) GetById(id uuid.UUID) (respondentEarning *model.RespondentEarning, err error) {
	return service.repo.GetById(id)
}

func (service *respondentEarningServiceImpl) Save(respondentEarning *model.RespondentEarning) (err error) {
	return service.repo.Save(respondentEarning)
}

func (service *respondentEarningServiceImpl) Commit(earning *model.RespondentEarning) (err error) {
	walletService := GetRespondentWalletService()
	wallet, err := service.GetWallet(earning)
	mutex := walletService.GetMutex(wallet)

	mutex.Lock()
	defer mutex.Unlock()
	now := time.Now()

	if !service.EnsureIsCommitable(earning, wallet) {
		message := fmt.Sprintf("Cannot commit this earning ")
		return errors.New(message)
	}

	err = model.DB.Transaction(func(tx *gorm.DB) error {

		// payout, err := service.Payout(earning)
		// if err != nil {
		// 	return err
		// }

		earning.Status = model.ORDER_EARNING_STATUS_COMPLETED
		earning.CommittedAt = &now

		if err = wallet.CommiteEarning(earning.OrderEarning); err != nil {
			return err
		}
		// earning.StripePayoutID = &payout.ID

		if tx.Save(earning); err != nil {
			return err
		}
		if tx.Save(wallet); err != nil {
			return err
		}

		return nil
	})
	return err
}

func (service *respondentEarningServiceImpl) Payout(earning *model.RespondentEarning) (stripPayout *stripe.Payout, err error) {
	payoutParams := &stripe.PayoutParams{
		Amount:      stripe.Int64(int64(earning.Amount)), // Amount in cents (e.g., $10.00)
		Currency:    stripe.String("usd"),
		Method:      stripe.String("instant"),
		Destination: stripe.String("your_customer_account_id"), // Customer's Stripe account ID
	}

	stripPayout, err = payout.New(payoutParams)
	if err != nil {
		message := fmt.Sprintf("Could not init order payment: %s", err.Error())
		return nil, errors.New(message)

	}
	return stripPayout, err
}

func (service *respondentEarningServiceImpl) EnsureIsCommitable(earning *model.RespondentEarning, wallet *model.RespondentWallet) (canCommit bool) {
	//CHECK IF WALLET IS LOCKED FOR DEPOSIT
	if earning.Status == model.ORDER_EARNING_STATUS_PENDING {
		return true
	}
	return false
}

func (service *respondentEarningServiceImpl) GetWallet(earning *model.RespondentEarning) (wallet *model.RespondentWallet, err error) {
	respondentRepo := repository.GetRespondentRepository()
	walletRepo := repository.GetRespondentWalletRepository()

	respondent, err := respondentRepo.GetById(*earning.RespondentID)
	if err != nil {
		return nil, err
	}

	if wallet, err = walletRepo.GetByRespondent(*respondent); err != nil {
		return nil, err
	}
	return wallet, nil
}

func (service *respondentEarningServiceImpl) Delete(respondentEarning *model.RespondentEarning) (err error) {
	err = service.repo.Delete(respondentEarning)

	return err
}

func (service *respondentEarningServiceImpl) DeleteById(id uuid.UUID) (respondentEarning *model.RespondentEarning, err error) {
	respondentEarning, err = service.repo.DeleteById(id)
	return respondentEarning, err
}
