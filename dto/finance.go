package dto

import (
	"time"

	"github.com/nkamuo/rasta-server/data/pagination"
)

type FinancialPageRequest struct {
	pagination.Page
	//
	From   *time.Time `json:"from,omitempty;" form:"from"`
	To     *time.Time `json:"to" form:"to"`
	Status *string    `json:"status" form:"status"`
}
