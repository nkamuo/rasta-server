package service

import (
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/nkamuo/rasta-server/model"
	"github.com/nkamuo/rasta-server/repository"
)

var sessionService RespondentSessionService
var sessionRepoMutext *sync.Mutex = &sync.Mutex{}

func GetRespondentSessionService() RespondentSessionService {
	sessionRepoMutext.Lock()
	if sessionService == nil {
		sessionService = &sessionServiceImpl{repo: repository.GetRespondentSessionRepository()}
	}
	sessionRepoMutext.Unlock()
	return sessionService
}

type RespondentSessionService interface {
	GetById(id uuid.UUID) (session *model.RespondentSession, err error)
	Close(session *model.RespondentSession) (err error)
	UpdateLocationEntry(session *model.RespondentSession, locationEntry model.RespondentSessionLocationEntry) (err error)
	Save(session *model.RespondentSession) (err error)
	Delete(session *model.RespondentSession) (error error)
}

type sessionServiceImpl struct {
	repo repository.RespondentSessionRepository
}

func (service *sessionServiceImpl) GetById(id uuid.UUID) (session *model.RespondentSession, err error) {
	return service.repo.GetById(id)
}

func (service *sessionServiceImpl) Close(session *model.RespondentSession) (err error) {
	now := time.Now()
	session.EndedAt = &now
	*session.Active = false
	return service.repo.Save(session)
}

func (service *sessionServiceImpl) Save(session *model.RespondentSession) (err error) {
	return service.repo.Save(session)
}

func (service *sessionServiceImpl) Delete(session *model.RespondentSession) (err error) {
	err = service.repo.Delete(session)

	return err
}

func (service *sessionServiceImpl) DeleteById(id uuid.UUID) (session *model.RespondentSession, err error) {
	session, err = service.repo.DeleteById(id)
	return session, err
}

func (service *sessionServiceImpl) UpdateLocationEntry(session *model.RespondentSession, locationEntry model.RespondentSessionLocationEntry) (err error) {
	orderService := GetOrderService()
	orderRepo := repository.GetOrderRepository()

	var fulfilments []model.OrderFulfilment

	query := model.DB.Where("session_id = ?", session.ID) //.Preload("Order.Items")

	if err := query.Find(&fulfilments).Error; err != nil {
		return err
	}

	for _, fulfilment := range fulfilments {
		if order, err := orderRepo.GetByFulfilment(fulfilment); err != nil {

		} else {
			/*err := */ go orderService.UpdateResponderLocationEntry(order, locationEntry)
		}
	}

	return nil
}
