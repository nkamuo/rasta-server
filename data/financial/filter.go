package financial

import (
	// "time"

	"github.com/nkamuo/rasta-server/dto"
	"gorm.io/gorm"
)

func FilterRequest(value interface{}, filter *dto.FinancialPageRequest, db *gorm.DB) func(db *gorm.DB) *gorm.DB {

	return func(db *gorm.DB) *gorm.DB {

		if filter.From != nil {
			db = db.Where("created_at >= ?", *filter.From)
		}

		if filter.To != nil {
			db = db.Where("created_at <= ?", *filter.To)
		}
		// if filter.Status != nil {
		// 	db = db.Where("status = ?", *filter.Status)
		// }

		return db
	}
}
