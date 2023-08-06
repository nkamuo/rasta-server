package repository

import (
	"sync"

	"github.com/google/uuid"
	"github.com/nkamuo/rasta-server/model"
	"gorm.io/gorm"
)

var companyRepo CompanyRepository
var companyRepoMutext *sync.Mutex = &sync.Mutex{}

func GetCompanyRepository() CompanyRepository {
	companyRepoMutext.Lock()
	if companyRepo == nil {
		companyRepo = &companyRepository{db: model.DB}
	}
	companyRepoMutext.Unlock()
	return companyRepo
}

type CompanyRepository interface {
	FindAll(page int, limit int) (companys []model.Company, total int64, err error)
	GetById(id uuid.UUID) (company *model.Company, err error)
	GetByEmail(email string) (company *model.Company, err error)
	GetByPhone(phone string) (company *model.Company, err error)
	GetByUser(user model.User) (company *model.Company, err error)
	GetByUserId(userID uuid.UUID) (company *model.Company, err error)
	Save(company *model.Company) (err error)
	Delete(company *model.Company) (error error)
	DeleteById(id uuid.UUID) (company *model.Company, err error)
}

type companyRepository struct {
	db *gorm.DB
}

func (repo *companyRepository) FindAll(page int, limit int) (companys []model.Company, total int64, err error) {
	offset := (page - 1) * limit

	err = repo.db.
		Model(&model.Company{}).
		Count(&total).
		Limit(limit).
		Offset(offset).
		Order("created_at desc").
		Find(&companys).Error
	if err != nil {
		return
	}
	return
}

func (repo *companyRepository) GetById(id uuid.UUID) (company *model.Company, err error) {
	if err = model.DB. /*.Joins("OperatorUser")*/ First(&company, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return company, nil
}

func (repo *companyRepository) GetByEmail(email string) (company *model.Company, err error) {
	if err = model.DB.Where("email = ?", email).First(&company).Error; err != nil {
		return nil, err
	}
	return company, nil
}

func (repo *companyRepository) GetByPhone(phone string) (company *model.Company, err error) {
	if err = model.DB.Where("phone = ?", phone).First(&company).Error; err != nil {
		return nil, err
	}
	return company, nil
}

func (repo *companyRepository) GetByUserId(userId uuid.UUID) (company *model.Company, err error) {
	if err = model.DB.
		Joins("JOIN users on users.id = companys.operator_user_id").
		Where("users.id = ?", userId).First(&company).Error; err != nil {
		return nil, err
	}
	return company, nil
}

func (repo *companyRepository) GetByUser(user model.User) (company *model.Company, err error) {
	return repo.GetByUserId(user.ID)
}

func (repo *companyRepository) Save(company *model.Company) (err error) {
	if (uuid.UUID{} == company.ID) {
		//NEW - No ID yet
		return repo.db.Create(&company).Error
	}
	return repo.db.Updates(&company).Error
}

func (repo *companyRepository) Delete(company *model.Company) (err error) {
	return repo.db.Delete(&company).Error
}

func (repo *companyRepository) DeleteById(id uuid.UUID) (company *model.Company, err error) {
	company, err = repo.GetById(id)
	if err != nil {
		return nil, err
	}
	err = repo.db.Delete(&company).Error
	return company, err
}
