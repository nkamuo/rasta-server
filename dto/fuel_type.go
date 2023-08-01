package dto

type FuelTypeCreationInput struct {
	Code        string `json:"code" binding:""`
	Rate        uint64 `json:"rate" binding:"required"`
	Title       string `json:"title" binding:""`
	ShortName   string `json:"shortName" binding:""`
	Description string `json:"description" binding:""`
	Published   bool   `json:"published" binding:"required"`
}

type FuelTypeUpdateInput struct {
	Code        *string `json:"code" binding:""`
	Rate        *uint64 `json:"rate" binding:"required"`
	Title       *string `json:"title" binding:""`
	ShortName   *string `json:"shortName" binding:""`
	Description *string `json:"description" binding:""`
	Published   *bool   `json:"published" binding:"required"`
}
