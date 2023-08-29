package service

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"errors"
	"fmt"

	// "go/token"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/nkamuo/rasta-server/initializers"
	"github.com/nkamuo/rasta-server/model"
	"github.com/nkamuo/rasta-server/repository"
)

var userPasswordResetRequestService UserPasswordResetRequestService
var userPasswordResetRequestRepoMutext *sync.Mutex = &sync.Mutex{}

func GetUserPasswordResetRequestService() UserPasswordResetRequestService {
	userPasswordResetRequestRepoMutext.Lock()
	if userPasswordResetRequestService == nil {
		userPasswordResetRequestService = &userPasswordResetRequestServiceImpl{repo: repository.GetUserPasswordResetRequestRepository()}
	}
	userPasswordResetRequestRepoMutext.Unlock()
	return userPasswordResetRequestService
}

type UserPasswordResetRequestService interface {
	GetById(id uuid.UUID) (userPasswordResetRequest *model.UserPasswordResetRequest, err error)
	GenerateForUser(user *model.User) (request *model.UserPasswordResetRequest, err error)
	Fullfil(request model.UserPasswordResetRequest, newPassword string) (err error)
	Save(userPasswordResetRequest *model.UserPasswordResetRequest) (err error)
	Delete(userPasswordResetRequest *model.UserPasswordResetRequest) (error error)
}

type userPasswordResetRequestServiceImpl struct {
	repo repository.UserPasswordResetRequestRepository
}

func (service *userPasswordResetRequestServiceImpl) GetById(id uuid.UUID) (userPasswordResetRequest *model.UserPasswordResetRequest, err error) {
	return service.repo.GetById(id)
}

func (service *userPasswordResetRequestServiceImpl) Save(userPasswordResetRequest *model.UserPasswordResetRequest) (err error) {
	return service.repo.Save(userPasswordResetRequest)
}

func (service *userPasswordResetRequestServiceImpl) Delete(userPasswordResetRequest *model.UserPasswordResetRequest) (err error) {
	err = service.repo.Delete(userPasswordResetRequest)

	return err
}

func (service *userPasswordResetRequestServiceImpl) GenerateForUser(user *model.User) (request *model.UserPasswordResetRequest, err error) {
	_token, err := generateRandomToken() //service.repo.Delete(userPasswordResetRequest)
	if err != nil {
		return nil, err
	}

	hashedToken := string(_token)
	expiresAt := time.Now().Add(time.Hour)

	request = &model.UserPasswordResetRequest{
		UserID:    &user.ID,
		Token:     hashedToken,
		ExpiresAt: expiresAt,
	}

	if err := service.Save(request); err != nil {
		return nil, err
	}

	return request, err
}

func (service *userPasswordResetRequestServiceImpl) Fullfil(request model.UserPasswordResetRequest, newPassword string) (err error) {

	userService := GetUserService()

	now := time.Now()

	if now.After(request.ExpiresAt) {
		//DELETE EXPIRED requests
		if err := service.Delete(&request); err != nil {
			return err
		}
		return errors.New("Token expired: Please try again")
	}

	user, err := userService.GetById(*request.UserID)
	if err != nil {
		message := fmt.Sprintf("Error loading user info: %s", err.Error())
		return errors.New(message)
	}

	if err := userService.HashUserPassword(user, newPassword); err != nil {
		message := fmt.Sprintf("Error updating user password: %s", err.Error())
		return errors.New(message)
	}

	return nil
}

func (service *userPasswordResetRequestServiceImpl) DeleteById(id uuid.UUID) (userPasswordResetRequest *model.UserPasswordResetRequest, err error) {
	userPasswordResetRequest, err = service.repo.DeleteById(id)
	return userPasswordResetRequest, err
}

func generateRandomToken() ([]byte, error) {
	token := make([]byte, 32) // Change the length according to your needs
	_, err := rand.Read(token)
	if err != nil {
		return nil, err
	}
	return token, nil
}

func generateHMAC(data []byte, token []byte) []byte {
	appSecret := initializers.CONFIG.APP_SECRET
	h := hmac.New(sha256.New, []byte(appSecret))
	h.Write(data)
	h.Write(token)
	return h.Sum(nil)
}

func validateHMAC(data []byte, hmacToken []byte) bool {
	calculatedHMAC := generateHMAC(data, data) // Recalculate HMAC using the received data and token
	return hmac.Equal(calculatedHMAC, hmacToken)
}
