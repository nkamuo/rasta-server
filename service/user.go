package service

import (
	"encoding/base64"
	"errors"
	"fmt"
	"sync"

	"github.com/google/uuid"
	"github.com/nkamuo/rasta-server/initializers"
	"github.com/nkamuo/rasta-server/model"
	"github.com/nkamuo/rasta-server/repository"
	"github.com/stripe/stripe-go/v74"
	"github.com/stripe/stripe-go/v74/customer"
	"golang.org/x/crypto/bcrypt"
)

var userService UserService
var userRepoMutext *sync.Mutex = &sync.Mutex{}

func GetUserService() UserService {
	userRepoMutext.Lock()
	if userService == nil {
		userService = &userServiceImpl{repo: repository.GetUserRepository()}
	}
	userRepoMutext.Unlock()
	return userService
}

type UserService interface {
	GetById(id uuid.UUID, preload ...string) (user *model.User, err error)
	// GetByEmail(email string) (user *model.User, err error)
	// GetByPhone(phone string) (user *model.User, err error)
	UpdateStripeCustomer(user *model.User, save bool) (customer *stripe.Customer, err error)
	Save(user *model.User) (err error)
	Delete(user *model.User) (error error)
	HashUserPassword(user *model.User, Password string) (err error)
	ValidateEmail(user *model.User) (err error)
}

type userServiceImpl struct {
	repo repository.UserRepository
}

func (service *userServiceImpl) GetById(id uuid.UUID, preload ...string) (user *model.User, err error) {
	return service.repo.GetById(id, preload...)
}

func (service *userServiceImpl) Save(user *model.User) (err error) {
	return service.repo.Save(user)
}

func (service *userServiceImpl) ValidateEmail(user *model.User) (err error) {
	return nil
}

func (service *userServiceImpl) ResetPassword(user *model.User) (err error) {

	requestService := GetUserPasswordResetRequestService()
	request, err := requestService.GenerateForUser(user)
	if err != nil {
		message := fmt.Sprintf("%s", err.Error())
		return errors.New(message)
	}

	token := request.Token
	tokenString := []byte(token)
	appUrl := initializers.CONFIG.APP_URL

	resetLink := fmt.Sprintf("%s/reset?token=%s", appUrl, base64.URLEncoding.EncodeToString(tokenString))

	fmt.Print(resetLink)

	return nil
}

func (service *userServiceImpl) Delete(user *model.User) (err error) {
	err = service.repo.Delete(user)

	return err
}

func (service *userServiceImpl) DeleteById(id uuid.UUID) (user *model.User, err error) {
	user, err = service.repo.DeleteById(id)
	return user, err
}

func (service *userServiceImpl) HashUserPassword(user *model.User, Password string) (err error) {
	if err != nil {
		return err
	}
	var password *model.UserPassword

	password, err = service.repo.GetPassword(user)
	if err != nil {
		if err.Error() == "record not found" {
			password = &model.UserPassword{}
		} else {
			return err
		}
	}
	//turn password into hash
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(Password), bcrypt.DefaultCost)
	password.HashedPassword = string(hashedPassword)
	user.Password = password

	return nil
}

func (service *userServiceImpl) UpdateStripeCustomer(user *model.User, save bool) (_customer *stripe.Customer, err error) {
	var userService = GetUserService()
	// customer.New()

	// var _customer stripe.Customer;
	if user.StripeCustomerID == nil {
		_customer, err = service.CreateStripeCustomer(user)
		if err != nil {
			return nil, err
		}
		user.StripeCustomerID = &_customer.ID

		if save {
			/**
			 * Update the database to work with this system
			 */
			if err = userService.Save(user); err != nil {
				return _customer, err
			}
		}
	} else {
		params := &stripe.CustomerParams{}
		_customer, err = customer.Get(*user.StripeCustomerID, params)
		if err != nil {
			return nil, err
		}

		shouldUpdate := false
		if _customer.Name != user.FullName() || _customer.Email != user.Email || _customer.Phone != user.Phone {
			shouldUpdate = true
		}

		if shouldUpdate {
			params = &stripe.CustomerParams{
				Name:  stripe.String(user.FullName()),
				Email: stripe.String(user.Email),
				Phone: stripe.String(user.Phone),
			}

			if _customer, err = customer.Update(_customer.ID, params); err != nil {
				return nil, err
			}
		}
	}

	return _customer, err
}

func (service *userServiceImpl) CreateStripeCustomer(user *model.User) (_customer *stripe.Customer, err error) {
	params := &stripe.CustomerParams{
		Name:  stripe.String(user.FullName()),
		Email: stripe.String(user.Email),
		Phone: stripe.String(user.Phone),
	}

	_customer, err = customer.New(params)
	if err != nil {
		return nil, err
	}
	return
}
