package service

import (
	"sync"

	"github.com/google/uuid"
	"github.com/nkamuo/rasta-server/model"
	"github.com/nkamuo/rasta-server/repository"
)

var priceService RespondentAccessProductPriceService
var priceRepoMutext *sync.Mutex = &sync.Mutex{}

func GetRespondentAccessProductPriceService() RespondentAccessProductPriceService {
	priceRepoMutext.Lock()
	if priceService == nil {
		priceService = &priceServiceImpl{repo: repository.GetRespondentAccessProductPriceRepository()}
	}
	priceRepoMutext.Unlock()
	return priceService
}

type RespondentAccessProductPriceService interface {
	GetById(id uuid.UUID, preload ...string) (price *model.RespondentAccessProductPrice, err error)
	GetByRespondent(respondent *model.Respondent, preload ...string) (price *model.RespondentAccessProductPrice, err error)
	// Close(price *model.RespondentAccessProductPrice) (err error)
	Save(price *model.RespondentAccessProductPrice) (err error)
	Delete(price *model.RespondentAccessProductPrice) (error error)
}

// RespondentAccessProductPrice

type priceServiceImpl struct {
	repo repository.RespondentAccessProductPriceRepository
}

func (service *priceServiceImpl) GetById(id uuid.UUID, preload ...string) (price *model.RespondentAccessProductPrice, err error) {
	return service.repo.GetById(id, preload...)
}

// func (service *priceServiceImpl) Close(price *model.RespondentAccessProductPrice) (err error) {
// 	now := time.Now()
// 	price.EndedAt = &now
// 	*price.Active = false
// 	return service.repo.Save(price)
// }

func (service priceServiceImpl) GetByRespondent(respondent *model.Respondent, preload ...string) (price *model.RespondentAccessProductPrice, err error) {
	return service.repo.GetActiveByRespondent(*respondent, preload...)
	// return nil, errors.New("Could not resolve Access Price for the given responder");
}

func (service *priceServiceImpl) Save(price *model.RespondentAccessProductPrice) (err error) {
	return service.repo.Save(price)
}

func (service *priceServiceImpl) Delete(price *model.RespondentAccessProductPrice) (err error) {
	err = service.repo.Delete(price)

	return err
}

func (service *priceServiceImpl) DeleteById(id uuid.UUID) (price *model.RespondentAccessProductPrice, err error) {
	price, err = service.repo.DeleteById(id)
	return price, err
}
