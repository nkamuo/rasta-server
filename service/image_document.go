package service

import (
	"sync"

	"github.com/google/uuid"
	"github.com/nkamuo/rasta-server/model"
	"github.com/nkamuo/rasta-server/repository"
)

var documentService ImageDocumentService
var documentRepoMutext *sync.Mutex = &sync.Mutex{}

func GetImageDocumentService() ImageDocumentService {
	documentRepoMutext.Lock()
	if documentService == nil {
		documentService = &documentServiceImpl{repo: repository.GetImageDocumentRepository()}
	}
	documentRepoMutext.Unlock()
	return documentService
}

type ImageDocumentService interface {
	GetById(id uuid.UUID) (document *model.ImageDocument, err error)
	Save(document *model.ImageDocument) (err error)
	Delete(document *model.ImageDocument) (error error)
}

type documentServiceImpl struct {
	repo repository.ImageDocumentRepository
}

func (service *documentServiceImpl) GetById(id uuid.UUID) (document *model.ImageDocument, err error) {
	return service.repo.GetById(id)
}

func (service *documentServiceImpl) Save(document *model.ImageDocument) (err error) {
	return service.repo.Save(document)
}

func (service *documentServiceImpl) Delete(document *model.ImageDocument) (err error) {
	err = service.repo.Delete(document)

	return err
}

func (service *documentServiceImpl) DeleteById(id uuid.UUID) (document *model.ImageDocument, err error) {
	document, err = service.repo.DeleteById(id)
	return document, err
}
