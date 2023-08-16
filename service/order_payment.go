package service

import (
	"errors"
	"fmt"
	"log"
	"sync"

	"github.com/google/uuid"
	"github.com/nkamuo/rasta-server/model"
	"github.com/nkamuo/rasta-server/repository"
	"github.com/stripe/stripe-go/v74"
	"github.com/stripe/stripe-go/v74/paymentintent"
	// "github.com/stripe/stripe-go/v74/paymentmethod"
)

var paymentService OrderPaymentService
var paymentRepoMutext *sync.Mutex = &sync.Mutex{}

func GetOrderPaymentService() OrderPaymentService {
	paymentRepoMutext.Lock()
	if paymentService == nil {
		paymentService = &paymentServiceImpl{repo: repository.GetOrderPaymentRepository()}
	}
	paymentRepoMutext.Unlock()
	return paymentService
}

type OrderPaymentService interface {
	GetById(id uuid.UUID) (payment *model.OrderPayment, err error)
	// GetByEmail(email string) (payment *model.OrderPayment, err error)
	// GetByPhone(phone string) (payment *model.OrderPayment, err error)
	InitOrderPayment(order *model.Order) (payment *model.OrderPayment, err error)
	Save(payment *model.OrderPayment) (err error)
	Delete(payment *model.OrderPayment) (error error)
}

type paymentServiceImpl struct {
	repo repository.OrderPaymentRepository
}

func (service *paymentServiceImpl) GetById(id uuid.UUID) (payment *model.OrderPayment, err error) {
	return service.repo.GetById(id)
}

func (service *paymentServiceImpl) Save(payment *model.OrderPayment) (err error) {
	return service.repo.Save(payment)
}

func (service *paymentServiceImpl) Delete(payment *model.OrderPayment) (err error) {
	err = service.repo.Delete(payment)

	return err
}

func (service *paymentServiceImpl) DeleteById(id uuid.UUID) (payment *model.OrderPayment, err error) {
	payment, err = service.repo.DeleteById(id)
	return payment, err
}

func (service *paymentServiceImpl) InitOrderPayment(order *model.Order) (payment *model.OrderPayment, err error) {

	orderService := GetOrderService()
	orderTotal, err := order.CalculateTotal()

	if nil != err {
		message := fmt.Sprintf("Could not calculate order total: %s", err.Error())
		return nil, errors.New(message)
	}
	paymentTitle := fmt.Sprintf("Payment for order[%s]", order.ID)

	// Create a PaymentIntent with amount and currency
	params := &stripe.PaymentIntentParams{
		Amount:   stripe.Int64(int64(orderTotal)),
		Currency: stripe.String(string(stripe.CurrencyUSD)),
		AutomaticPaymentMethods: &stripe.PaymentIntentAutomaticPaymentMethodsParams{
			Enabled: stripe.Bool(true),
		},
	}

	pi, err := paymentintent.New(params)
	log.Printf("pi.New: %v", pi.ClientSecret)

	if err != nil {
		return nil, err
	}

	payment = &model.OrderPayment{
		OrderID:            &order.ID,
		Amount:             int64(orderTotal),
		Title:              paymentTitle,
		Status:             "processing",
		StripeID:           &pi.ID,
		StripeClientSecret: &pi.ClientSecret,
	}

	if err := service.Save(payment); err != nil {
		return nil, err
	}

	order.PaymentID = &payment.ID

	if err := orderService.Save(order); err != nil {
		return nil, err
	}

	return payment, nil
}
