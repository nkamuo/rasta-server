package pagination

import (
	"math"

	"gorm.io/gorm"
)

type Page struct {
	Status     string      `json:"status,omitempty;" form:"status"`
	Search     string      `json:"search,omitempty;" form:"search"`
	Limit      *int        `json:"limit,omitempty;" form:"limit"`
	Page       int         `json:"page,omitempty;" form:"page"`
	Sort       string      `json:"sort,omitempty;" form:"sort"`
	TotalRows  int64       `json:"total_rows"`
	TotalPages int         `json:"total_pages"`
	Rows       interface{} `json:"rows"`
}

func (p *Page) GetOffset() int {
	return (p.GetPage() - 1) * p.GetLimit()
}

func (p *Page) GetLimit() int {
	// if nil == pLimit{
	// 	pLimit = 10;
	// }
	// if p.Limit == 0 {
	// 	p.Limit = 10
	// }
	Limit := 10
	pLimit := p.Limit
	if nil == pLimit {
		pLimit = &Limit
	} else if *pLimit > 100 {
		*pLimit = 100
	}
	return *pLimit
}

func (p *Page) GetPage() int {
	if p.Page == 0 {
		p.Page = 1
	}
	return p.Page
}

func (p *Page) GetSort() string {
	if p.Sort == "" {
		p.Sort = "created_at desc"
	}
	return p.Sort
}

func Paginate(value interface{}, pagination *Page, db *gorm.DB) func(db *gorm.DB) *gorm.DB {
	var totalRows int64
	var totalPages int
	db.Model(value).Count(&totalRows)

	pagination.TotalRows = totalRows

	limit := pagination.GetLimit()

	if int64(limit) > totalRows {
		totalPages = 1
	} else {
		if limit == 0 {
			totalPages = 0
		} else {
			itemsPage := float64(float32(totalRows) / float32(pagination.GetLimit()))
			totalPages = int(math.Ceil(itemsPage))
		}
	}
	pagination.TotalPages = totalPages

	return func(db *gorm.DB) *gorm.DB {
		return db.Offset(pagination.GetOffset()).Limit(pagination.GetLimit()).Order(pagination.GetSort())
	}
}
