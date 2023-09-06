package dto

type WithdrawalRequest struct {
	Amount      uint64  `json:"amount" binding:"required,gte=5000"`
	Description *string `json:"description" binding:""`
}
