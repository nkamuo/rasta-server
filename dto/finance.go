package dto

import (
	"time"

	"github.com/nkamuo/rasta-server/data/pagination"
)

type FinancialPageRequest struct {
	pagination.Page
	//
	From   *time.Time `json:"sort,omitempty;" form:"sort"`
	To     *time.Time `json:"total_rows"`
	Status *string    `json:"total_pages"`
}
