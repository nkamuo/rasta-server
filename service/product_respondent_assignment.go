package service

import (
	"sync"

	"github.com/google/uuid"
	"github.com/nkamuo/rasta-server/model"
	"github.com/nkamuo/rasta-server/repository"
)

var assignmentService ProductRespondentAssignmentService
var assignmentRepoMutext *sync.Mutex = &sync.Mutex{}

func GetProductRespondentAssignmentService() ProductRespondentAssignmentService {
	assignmentRepoMutext.Lock()
	if assignmentService == nil {
		assignmentService = &assignmentServiceImpl{repo: repository.GetProductRespondentAssignmentRepository()}
	}
	assignmentRepoMutext.Unlock()
	return assignmentService
}

type ProductRespondentAssignmentService interface {
	GetById(id uuid.UUID) (assignment *model.ProductRespondentAssignment, err error)
	Save(assignment *model.ProductRespondentAssignment) (err error)
	Delete(assignment *model.ProductRespondentAssignment) (error error)
}

type assignmentServiceImpl struct {
	repo repository.ProductRespondentAssignmentRepository
}

func (service *assignmentServiceImpl) GetById(id uuid.UUID) (assignment *model.ProductRespondentAssignment, err error) {
	return service.repo.GetById(id)
}

func (service *assignmentServiceImpl) Save(assignment *model.ProductRespondentAssignment) (err error) {
	return service.repo.Save(assignment)
}

func (service *assignmentServiceImpl) Delete(assignment *model.ProductRespondentAssignment) (err error) {
	err = service.repo.Delete(assignment)

	return err
}

func (service *assignmentServiceImpl) DeleteById(id uuid.UUID) (assignment *model.ProductRespondentAssignment, err error) {
	assignment, err = service.repo.DeleteById(id)
	return assignment, err
}
