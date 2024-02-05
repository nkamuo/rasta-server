package repository

import (
	"sync"

	"github.com/google/uuid"
	"github.com/nkamuo/rasta-server/model"
	"gorm.io/gorm"
)

var balanceRepo RespondentAccessProductBalanceRepository
var balanceRepoMutext *sync.Mutex = &sync.Mutex{}

func GetRespondentAccessProductBalanceRepository() RespondentAccessProductBalanceRepository {
	balanceRepoMutext.Lock()
	if balanceRepo == nil {
		balanceRepo = &balanceRepository{db: model.DB}
	}
	balanceRepoMutext.Unlock()
	return balanceRepo
}

type RespondentAccessProductBalanceRepository interface {
	FindAll(page int, limit int) (balances []model.RespondentAccessProductBalance, total int64, err error)
	GetById(id uuid.UUID, prefetch ...string) (balance *model.RespondentAccessProductBalance, err error)
	GetActiveByRespondent(respondent model.Respondent, prefetch ...string) (balance *model.RespondentAccessProductBalance, err error)
	Save(balance *model.RespondentAccessProductBalance) (err error)
	Delete(balance *model.RespondentAccessProductBalance) (error error)
	DeleteById(id uuid.UUID) (balance *model.RespondentAccessProductBalance, err error)
}

type balanceRepository struct {
	db *gorm.DB
}

func (repo *balanceRepository) FindAll(page int, limit int) (balances []model.RespondentAccessProductBalance, total int64, err error) {
	offset := (page - 1) * limit

	err = repo.db.
		Model(&model.RespondentAccessProductBalance{}).
		Count(&total).
		Limit(limit).
		Offset(offset).
		Order("created_at desc").
		Find(&balances).Error
	if err != nil {
		return
	}
	return
}

func (repo *balanceRepository) GetById(id uuid.UUID, prefetch ...string) (balance *model.RespondentAccessProductBalance, err error) {
	query := model.DB

	for _, pFetch := range prefetch {
		query = query.Preload(pFetch)
	}

	if err = query.First(&balance, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return balance, nil
}

func (repo *balanceRepository) GetActiveByRespondent(respondent model.Respondent, prefetch ...string) (balance *model.RespondentAccessProductBalance, err error) {
	query := model.DB //.Preload("Assignments.Assignment.Product").
	query = query.Where("respondent_id = ? AND active = ? AND ended_at IS NULL", respondent.ID, true)

	for _, pFetch := range prefetch {
		query = query.Preload(pFetch)
	}

	if err = query.First(&balance).Error; err != nil {
		return nil, err
	}
	return balance, nil
}

func (repo *balanceRepository) Save(balance *model.RespondentAccessProductBalance) (err error) {
	if (uuid.UUID{} == balance.ID) {
		//NEW - No ID yet
		return repo.db.Create(&balance).Error
	}
	return repo.db.Updates(&balance).Error
}

func (repo *balanceRepository) Delete(balance *model.RespondentAccessProductBalance) (err error) {
	return repo.db.Delete(&balance).Error
}

func (repo *balanceRepository) DeleteById(id uuid.UUID) (balance *model.RespondentAccessProductBalance, err error) {
	balance, err = repo.GetById(id)
	if err != nil {
		return nil, err
	}
	err = repo.db.Delete(&balance).Error
	return balance, err
}
