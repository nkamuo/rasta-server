package service

import (
	"errors"
	"fmt"
	"sync"

	"github.com/google/uuid"
	"github.com/nkamuo/rasta-server/initializers"
	"github.com/nkamuo/rasta-server/model"
	"github.com/nkamuo/rasta-server/repository"
	"github.com/stripe/stripe-go/v74"
	"github.com/stripe/stripe-go/v74/paymentintent"
	"github.com/stripe/stripe-go/v74/payout"
	"gorm.io/gorm"
)

var respondentChargeService RespondentOrderChargeService
var respondentChargeRepoMutext *sync.Mutex = &sync.Mutex{}

func GetRespondentOrderChargeService() RespondentOrderChargeService {
	respondentChargeRepoMutext.Lock()
	if respondentChargeService == nil {
		respondentChargeService = &respondentChargeServiceImpl{repo: repository.GetRespondentOrderChargeRepository()}
	}
	respondentChargeRepoMutext.Unlock()
	return respondentChargeService
}

type RespondentOrderChargeService interface {
	GetById(id uuid.UUID) (charge *model.RespondentOrderCharge, err error)
	// GetByEmail(email string) (charge *model.RespondentOrderCharge, err error)
	// GetByPhone(phone string) (charge *model.RespondentOrderCharge, err error)'
	// Commit(charge *model.RespondentOrderCharge) (err error)
	UpdateAllCharges() (err error)
	Update(charge *model.RespondentOrderCharge) (err error)
	Commit(charge *model.RespondentOrderCharge) (err error)
	CreateForOrder(order *model.Order) (charge *model.RespondentOrderCharge, err error)
	Save(charge *model.RespondentOrderCharge) (err error)
	Delete(charge *model.RespondentOrderCharge) (error error)
}

type respondentChargeServiceImpl struct {
	repo repository.RespondentOrderChargeRepository
}

func (service *respondentChargeServiceImpl) GetById(id uuid.UUID) (charge *model.RespondentOrderCharge, err error) {
	return service.repo.GetById(id)
}

func (service *respondentChargeServiceImpl) Save(charge *model.RespondentOrderCharge) (err error) {
	return service.repo.Save(charge)
}

func (service *respondentChargeServiceImpl) CreateForOrder(order *model.Order) (charge *model.RespondentOrderCharge, err error) {

	chargeRepo := repository.GetRespondentOrderChargeRepository()
	// orderService := GetOrderService();
	orderRepo := repository.GetOrderRepository()
	productService := GetProductService()
	// fulfilmentService := GetOrderFulfilmentService()
	// fulfilmentRepo := repository.GetOrderFulfilmentRepository();

	rOrder, err := orderRepo.GetById(order.ID, "Fulfilment.Responder", "Items.Product")
	if err != nil {
		return nil, err
	}

	fulfilment := rOrder.Fulfilment

	if nil == fulfilment {
		message := fmt.Sprintf("Order fulfilment not active")
		return nil, errors.New(message)
	}

	if nil == fulfilment.Responder {
		message := fmt.Sprintf("Order fulfilment responder not found")
		return nil, errors.New(message)
	}

	respondent := fulfilment.Responder

	if length := len(*rOrder.Items); length != 1 {
		message := fmt.Sprintf("Charges can only be generated for orders with one request. %d found", length)
		return nil, errors.New(message)
	}

	request := (*rOrder.Items)[0]

	existing, err := chargeRepo.GetByRequest(request)
	if err != nil {
		if err.Error() != "record not found" {
			return nil, err
		}
	}
	if existing != nil {
		message := fmt.Sprintf("There is already a charge for this order/request")
		return nil, errors.New(message)
	}

	product, err := productService.GetById(*request.ProductID)
	if err != nil {
		message := fmt.Sprintf("Error Fetching product[id:%s]: %s", request.ProductID, err.Error())
		return nil, errors.New(message)
	}

	amount := respondent.BillingAmount
	if amount == nil {
		amount = &initializers.CONFIG.RESPONDENT_ORDER_CHARGE_AMOUNT
	}
	// if *amount == 0{
	// 	amount = 2500;
	// }

	label := fmt.Sprintf("Charge for %s on %s", product.Label, fulfilment.CreatedAt)

	err = model.DB.Transaction(func(tx *gorm.DB) error {
		// Create a PaymentIntent with amount and currency
		params := &stripe.PaymentIntentParams{
			Amount:   stripe.Int64(int64(*amount)),
			Currency: stripe.String(string(stripe.CurrencyUSD)),
			AutomaticPaymentMethods: &stripe.PaymentIntentAutomaticPaymentMethodsParams{
				Enabled: stripe.Bool(true),
			},
		}
		pi, err := paymentintent.New(params)
		if err != nil {
			return err
		}

		charge = &model.RespondentOrderCharge{
			RequestID:                 &request.ID,
			RespondentID:              &respondent.ID,
			Amount:                    *amount,
			Label:                     label,
			Status:                    model.ORDER_EARNING_STATUS_PENDING,
			StripePaymentID:           &pi.ID,
			StripePaymentClientSecret: &pi.ClientSecret,
		}

		if err := tx.Save(charge).Error; err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return charge, nil
}

func (service *respondentChargeServiceImpl) Update(charge *model.RespondentOrderCharge) (err error) {
	if nil == charge.StripePaymentID {
		return errors.New("No stripe Payment Intent ID is associated with this charge")
	}

	pIntent, err := paymentintent.Get(*charge.StripePaymentID, nil)

	if err != nil {
		message := fmt.Sprintf("Could not fetch payment intent[%s]: %s", *charge.StripePaymentID, err.Error())
		return errors.New(message)
	}
	charge.Status = string(pIntent.Status)

	if err := service.Save(charge); err != nil {
		return err
	}

	return nil
}

func (service *respondentChargeServiceImpl) UpdateAllCharges() (err error) {

	var charges []model.RespondentOrderCharge

	query := model.DB

	if err = query.Find(&charges).Error; err != nil {
		return err
	}
	for _, charge := range charges {
		if err := service.Update(&charge); err != nil {
			message := fmt.Sprintf("There was an error updating charge %s", charge.ID)
			fmt.Println(message)
		}
	}
	return nil
}

func (service *respondentChargeServiceImpl) Payout(charge *model.RespondentOrderCharge) (stripPayout *stripe.Payout, err error) {
	payoutParams := &stripe.PayoutParams{
		Amount:      stripe.Int64(int64(charge.Amount)), // Amount in cents (e.g., $10.00)
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

func (service *respondentChargeServiceImpl) Commit(charge *model.RespondentOrderCharge) (err error) {

	if !service.EnsureIsCommitable(charge) {
		message := fmt.Sprintf("Charge of state \"%s\" cannot be commited", charge.Status)
		return errors.New(message)
	}

	charge.Status = model.ORDER_EARNING_STATUS_COMPLETED

	if err := service.Save(charge); err != nil {
		return err
	}

	return nil
}

func (service *respondentChargeServiceImpl) EnsureIsCommitable(charge *model.RespondentOrderCharge) (canCommit bool) {
	//CHECK IF WALLET IS LOCKED FOR DEPOSIT
	if charge.Status == model.ORDER_EARNING_STATUS_PENDING {
		return true
	}
	return false
}

func (service *respondentChargeServiceImpl) GetWallet(charge *model.RespondentOrderCharge) (wallet *model.RespondentWallet, err error) {
	respondentRepo := repository.GetRespondentRepository()
	walletRepo := repository.GetRespondentWalletRepository()

	respondent, err := respondentRepo.GetById(*charge.RespondentID)
	if err != nil {
		return nil, err
	}

	if wallet, err = walletRepo.GetByRespondent(*respondent); err != nil {
		return nil, err
	}
	return wallet, nil
}

func (service *respondentChargeServiceImpl) Delete(charge *model.RespondentOrderCharge) (err error) {
	err = service.repo.Delete(charge)

	return err
}

func (service *respondentChargeServiceImpl) DeleteById(id uuid.UUID) (charge *model.RespondentOrderCharge, err error) {
	charge, err = service.repo.DeleteById(id)
	return charge, err
}

// func (service *respondentChargeServiceImpl) Commit(charge *model.RespondentOrderCharge) (err error) {
// 	walletService := GetRespondentWalletService()
// 	wallet, err := service.GetWallet(charge)
// 	mutex := walletService.GetMutex(wallet)

// 	mutex.Lock()
// 	defer mutex.Unlock()
// 	now := time.Now()

// 	if !service.EnsureIsCommitable(charge, wallet) {
// 		message := fmt.Sprintf("Cannot commit this charge ")
// 		return errors.New(message)
// 	}

// 	err = model.DB.Transaction(func(tx *gorm.DB) error {

// 		// payout, err := service.Payout(charge)
// 		// if err != nil {
// 		// 	return err
// 		// }

// 		charge.Status = model.ORDER_EARNING_STATUS_COMPLETED
// 		charge.CommittedAt = &now

// 		if err = wallet.CommiteCharge(charge); err != nil {
// 			return err
// 		}
// 		// charge.StripePayoutID = &payout.ID

// 		if tx.Save(charge); err != nil {
// 			return err
// 		}
// 		if tx.Save(wallet); err != nil {
// 			return err
// 		}

// 		return nil
// 	})
// 	return err
// }
